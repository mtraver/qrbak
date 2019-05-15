// Program qrbak securely and durably backs up private keys using QR codes.
package main

import (
	"crypto/sha256"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path"
	"strings"

	qrcode "github.com/skip2/go-qrcode"

	"github.com/mtraver/qrbak/gpg"
	"github.com/mtraver/qrbak/pdf"
)

var (
	keyID  string
	outDir string

	verbose     bool
	saveImages  bool
	saveTxt     bool
	numQRCodes  int
	codesPerRow int
	pageSize    = pdf.PageSizeValue(pdf.Letter)
)

// fatalf calls fmt.Printf with the given arguments (adding a newline) and then exist with status 1.
func fatalf(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
	os.Exit(1)
}

// vprintf is equivalent to fmt.Printf, but it only prints if verbose output is enabled.
func vprintf(format string, a ...interface{}) {
	if verbose {
		fmt.Printf(format, a...)
	}
}

func split(s string, n int) []string {
	if n <= 0 {
		return []string{}
	} else if n > len(s) {
		// TODO(mtraver) Should this be an error? Should we return []string{s} instead?
		n = len(s)
	}
	partLen := int(math.Ceil(float64(len(s)) / float64(n)))

	parts := make([]string, n)
	for i := 0; i < n; i++ {
		if i == n-1 {
			parts[i] = s[i*partLen:]
		} else {
			parts[i] = s[i*partLen : i*partLen+partLen]
		}
	}

	return parts
}

func init() {
	flag.Usage = func() {
		message := `Usage: qrbak [options] keyid outdir

qrbak exports a private key from gpg, encrypts it with AES256, encodes the
ciphertext in base 64, creates QR codes from the base 64 text, and finally
produces a PDF containing the QR codes. The QR codes are rendered in the PDF
as a grid from left to right and from top to bottom.

To reconstruct the private key, scan each QR code, concatenate the text, decode
the base 64 text, and finally decrypt using the same passphrase you gave when
generating the PDF.

Positional arguments (required):
  keyid
        ID of GPG key
  outdir
        directory in which to save output

Options:
`
		fmt.Fprintf(flag.CommandLine.Output(), message)
		flag.PrintDefaults()
	}

	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.BoolVar(&saveImages, "img", false, "generate PNGs, one per QR code, in addition to a PDF")
	flag.BoolVar(&saveTxt, "txt", false, "save a text file containing the encrypted, base 64-encoded key in addition to a PDF")
	flag.Var(&pageSize, "pagesize", "PDF page size")
	flag.IntVar(&numQRCodes, "codes", 27, "number of QR codes to generate")
	flag.IntVar(&codesPerRow, "row", 3, "number of QR codes per row in the PDF")
}

func parseFlags() error {
	flag.Parse()
	if flag.NArg() != 2 {
		return fmt.Errorf("")
	}

	keyID = flag.Args()[0]
	outDir = flag.Args()[1]

	if keyID == "" {
		return fmt.Errorf("keyid must be given")
	}

	if outDir == "" {
		return fmt.Errorf("outdir must be given")
	}

	if numQRCodes < 1 {
		return fmt.Errorf("-codes must be >= 1")
	}

	if codesPerRow < 1 {
		return fmt.Errorf("-row must be >= 1")
	}

	return nil
}

func main() {
	if err := parseFlags(); err != nil {
		if err.Error() != "" {
			fmt.Printf("%v\n", err)
		}
		flag.Usage()
		os.Exit(2)
	}

	if !gpg.Installed() {
		fatalf("gpg is not installed")
	}

	key, err := gpg.ExportPrivateKey(keyID)
	if err != nil {
		fatalf("Failed to export private key: %v", err)
	}
	vprintf("Private key is %d bytes\n", len(key))

	// The key's fingerprint will be used in filenames and included in the PDF.
	fingerprint, err := gpg.Fingerprint(keyID)
	if err != nil {
		fatalf("Failed to get key fingerprint: %v", err)
	}

	vprintf("\n")
	fmt.Println("Encrypting private key. Enter a passphrase. You will need it to")
	fmt.Println("recover your key from the QR codes. Keep it secret, keep it safe!")
	enc, err := gpg.EncryptAES256(key)
	if err != nil {
		fatalf("Failed to encrypt private key: %v", err)
	}
	vprintf("\nEncrypted key is %d bytes\n", len(enc))

	// Convert the data to ASCII text so it can be encoded in a QR code.
	encb64 := base64.StdEncoding.EncodeToString(enc)
	vprintf("Encrypted, base 64-encoded key is %d bytes\n", len(encb64))

	filenameBase := fmt.Sprintf("qrbak_%s", strings.ToLower(fingerprint))

	parts := split(encb64, numQRCodes)
	pngs := make([][]byte, numQRCodes)
	for i := range parts {
		vprintf("Generating QR code %d with %d bytes\n", i, len(parts[i]))
		pngs[i], err = qrcode.Encode(parts[i], qrcode.Highest, 512)
		if err != nil {
			fatalf("Failed to generate QR code: %v", err)
		}

		if saveImages {
			imgFilename := fmt.Sprintf("%s_%d.png", filenameBase, i)
			vprintf("Writing %s\n", imgFilename)
			if err := ioutil.WriteFile(path.Join(outDir, imgFilename), pngs[i], 0600); err != nil {
				fatalf("Failed to write image file: %v", err)
			}
		}
	}

	if saveTxt {
		txtFilename := filenameBase + ".txt"
		vprintf("Writing %s\n", txtFilename)
		if err := ioutil.WriteFile(path.Join(outDir, txtFilename), []byte(encb64), 0600); err != nil {
			fatalf("Failed to write text file: %v", err)
		}
	}

	doc := pdf.New(pngs, sha256.Sum256([]byte(encb64)), fingerprint, string(pageSize), codesPerRow)
	pdfFilename := filenameBase + ".pdf"
	vprintf("Writing %s\n", pdfFilename)
	if err := doc.OutputFileAndClose(path.Join(outDir, pdfFilename)); err != nil {
		fatalf("Failed to write PDF file: %v", err)
	}
}