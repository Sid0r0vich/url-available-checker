package pdf

import (
	"bytes"

	"github.com/Sid0r0vich/url-available-checker/internal/dto"
	"github.com/jung-kurt/gofpdf"
)

func GeneratePDFFromLinks(links []dto.Link) (*bytes.Buffer, error) {
	file := gofpdf.New("P", "mm", "A4", "")
	file.AddPage()
	file.SetFont("Arial", "", 12)
	file.Cell(0, 10, "Links availability report")
	file.Ln(12)

	for _, link := range links {
		status := "not available"
		if link.Availability {
			status = "available"
		}
		line := link.URL + ": " + status
		file.Cell(0, 10, line)
		file.Ln(10)
	}

	var buf bytes.Buffer
	err := file.Output(&buf)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}
