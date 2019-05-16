# Securely and durably back up your private key using QR codes printed on paper

[![GoDoc](https://godoc.org/github.com/mtraver/qrbak?status.svg)](https://godoc.org/github.com/mtraver/qrbak)
[![Go Report Card](https://goreportcard.com/badge/github.com/mtraver/qrbak)](https://goreportcard.com/report/github.com/mtraver/qrbak)

`go run main.go example@example.com ~/key_bak`

The command above makes a PDF containing a secure backup of the key identified by `example@example.com`,
saving it in the directory `~/key_bak`.

## Cool tell me more

You know that your private key is important. You know that you should back it up. You know that
when you do so you should encrypt it. You even know that the most reliable backup medium is ink
on paper, and because a key is a fairly small amount of data it's feasible to back it up that way.

But it would be annoying to manually transcribe hundreds or thousands of characters in the
event that you do need to recover your key from the backup. If only there were some kind of
machine-readable represen...oh wait there is, it's QR codes, the answer is QR codes.

qrbak takes your key's ID as input and produces a PDF of QR codes as output. Simple as that.

## Great how about all the details?

`go run main.go -h`

```
usage: qrbak [options] keyid outdir

qrbak does the following:
  1. Export a private key from gpg.
  2. Encrypt the private key with AES256 (you will be prompted for a passphrase).
  3. Encode the result of step 2 in base 64.
  4. Split the result of step 3 into chunks and make a QR code for each chunk.
  5. Create a PDF containing the QR codes, rendered in a grid from left to right
     and top to bottom.

  Steps 1-3 are equivalent to executing

    gpg --export-secret-keys $KEY_ID | gpg --cipher-algo AES256 --symmetric | base64

To reconstruct the private key and import it into gpg, follow these steps:
  1. Scan each QR code.
  2. Concatenate the content of the QR codes to get a single block of base 64 text.
  3. Decode the base 64 text to get the encrypted private key.
  4. Decrypt the output of step 3 using the same passphrase you gave when
     generating the PDF.
  5. Import the result into gpg using
       gpg --import

  If the result of step 2 above is in a file named b64.txt, this is equivalent
  to executing

    base64 --decode b64.txt | gpg --decrypt | gpg --import

Positional arguments (required):
  keyid
        ID of GPG key
  outdir
        directory in which to save output

Options:
  -codes int
      number of QR codes to generate (default 27)
  -img
      save PNGs, one per QR code, in addition to a PDF
  -pagesize value
      PDF page size (default Letter)
  -row int
      number of QR codes per row in the PDF (default 3)
  -txt
      save a text file containing the encrypted, base 64-encoded key in addition to a PDF
  -v  verbose output

```
