// Package pdf aids creation of a PDF containing a grid of QR codes.
package pdf

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	fpdf "github.com/jung-kurt/gofpdf/v2"
)

const (
	A3     = fpdf.PageSizeA3
	A4     = fpdf.PageSizeA4
	A5     = fpdf.PageSizeA5
	Legal  = fpdf.PageSizeLegal
	Letter = fpdf.PageSizeLetter

	// In reality there are 0.353 mm per point. But this is used for line spacing, and using a slightly
	// larger value makes the line spacing more visually appealing.
	mmPerPoint = 0.37
)

// sizes returns a map containing all valid sizes as keys. The values are all always true.
func sizes() map[string]bool {
	return map[string]bool{
		A3:     true,
		A4:     true,
		A5:     true,
		Letter: true,
		Legal:  true,
	}
}

// fontSize returns the font size to use when setting type on a page of the given size.
func fontSize(pageSize string) float64 {
	switch pageSize {
	case A3:
		return 11
	case A5:
		return 7
	case Letter, Legal:
		return 10
	default:
		return 9
	}
}

// PageSizeValue implements flag.Value, so it may be used with the flag package to get and validate a page size.
type PageSizeValue string

func (p *PageSizeValue) String() string {
	return string(*p)
}

func (p *PageSizeValue) Set(s string) error {
	s = strings.Title(strings.ToLower(s))

	validSizes := sizes()
	if _, ok := validSizes[s]; !ok {
		// Get the keys from validSizes to we can pretty-print the valid page sizes in the error.
		keys := make([]string, 0, len(validSizes))
		for k := range validSizes {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		return fmt.Errorf("valid sizes are %s", strings.Join(keys, ", "))
	}
	*p = PageSizeValue(s)
	return nil
}

type FooterFunc func(pageNumber int) string

func New(pngs [][]byte, codesPerRow int, pageSize string, footerFunc FooterFunc) *fpdf.Fpdf {
	pdf := fpdf.New("P", "mm", pageSize, "")

	// Set up the footer.
	pdf.SetFooterFunc(func() {
		text := footerFunc(pdf.PageNo())
		fs := fontSize(pageSize)

		pdf.SetY(-20)
		pdf.SetFont("Arial", "", fs)
		pdf.MultiCell(0, fs*mmPerPoint, text, "", "C", false)
	})
	pdf.AliasNbPages("")

	width, height := pdf.GetPageSize()
	left, _, right, bottom := pdf.GetMargins()

	// Compute the width of each QR code given how many per row.
	qrWidth := (width - left - right) / float64(codesPerRow)

	pdf.AddPage()

	initialX := pdf.GetX()
	for i := range pngs {
		if i != 0 && i%codesPerRow == 0 {
			// Set the x and y position for a new row.
			pdf.SetXY(initialX, pdf.GetY()+qrWidth)
		}

		if pdf.GetY()+qrWidth >= height-bottom {
			// The current position is reset when a page is added so we don't need to call SetX() or SetY().
			pdf.AddPage()
		}

		name := fmt.Sprintf("qr_%d", i)
		pdf.RegisterImageOptionsReader(name, fpdf.ImageOptions{ImageType: "PNG"}, bytes.NewReader(pngs[i]))
		pdf.Image(name, pdf.GetX(), pdf.GetY(), qrWidth, 0, false, "", 0, "")

		// Advance the x position.
		pdf.SetX(pdf.GetX() + qrWidth)
	}

	return pdf
}
