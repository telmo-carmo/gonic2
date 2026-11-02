package controllers

/*
Or use:  signintech/gopdf

*/
import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/phuslu/log"
	"seehuhn.de/go/pdf"
	"seehuhn.de/go/pdf/document"
	"seehuhn.de/go/pdf/font"
	"seehuhn.de/go/pdf/font/gofont"
	"seehuhn.de/go/pdf/graphics/color"
)

type Customer struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Country string  `json:"country"`
	Balance float64 `json:"balance"`
}

type ReportRequest struct {
	Title     string     `json:"title"`
	Customers []Customer `json:"customers"`
}

func Report1Handler(c *gin.Context) {
	var request ReportRequest

	if c.Request.Method != http.MethodPost {
		json1 := `
{
  "title": "Customer Report",
  "customers": [
    {
      "id": 1001,
      "name": "João Simões",
      "country": "Portugal",
      "balance": 220.50
    },
    {
      "id": 1002,
      "name": "Maria Nuñez",
      "country": "Spain",
      "balance": 8920.75
    }
  ]
}`
		request = ReportRequest{}
		if err := json.Unmarshal([]byte(json1), &request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		log.Info().Msg("Using default JSON for report request")
	} else if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	regularFont, err := gofont.Regular.NewComposite(nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load report font: " + err.Error()})
		return
	}
	boldFont, err := gofont.Bold.NewComposite(nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load report font: " + err.Error()})
		return
	}

	var buffer bytes.Buffer
	doc, err := document.WriteMultiPage(&buffer, document.A4, pdf.V2_0, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not initialize PDF writer: " + err.Error()})
		return
	}

	headers := []string{
		"ID",
		"Name",
		"Country",
		"Balance",
	}

	widths := []float64{
		25,
		70,
		45,
		35,
	}

	pageHeight := document.A4.URy
	margin := mmToPoints(15)
	columnWidths := make([]float64, len(widths))
	for i, width := range widths {
		columnWidths[i] = mmToPoints(width)
	}

	page := doc.AddPage()
	titleTop := pageHeight - margin
	titleBaseline := titleTop - mmToPoints(12)/2 - 18*0.75
	drawReportText(page, boldFont, 18, margin, titleBaseline, document.A4.URx-2*margin, 0.5, request.Title)

	tableTop := titleTop - mmToPoints(17)
	headerHeight := mmToPoints(9)
	headerAlignment := []float64{0.5, 0.5, 0.5, 0.5}
	drawReportRow(page, headers, columnWidths, margin, tableTop, headerHeight, boldFont, 11, headerAlignment, true)
	currentY := tableTop - headerHeight
	rowHeight := mmToPoints(8)

	for _, customer := range request.Customers {
		if currentY-rowHeight < margin {
			if err := page.Close(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate PDF page: " + err.Error()})
				return
			}
			page = doc.AddPage()
			currentY = pageHeight - margin
			drawReportRow(page, headers, columnWidths, margin, currentY, headerHeight, boldFont, 11, headerAlignment, true)
			currentY -= headerHeight
		}

		values := []string{
			fmt.Sprintf("%d", customer.ID),
			customer.Name,
			customer.Country,
			fmt.Sprintf("%.2f EUR", customer.Balance),
		}
		alignments := []float64{0.5, 0, 0, 1}
		drawReportRow(page, values, columnWidths, margin, currentY, rowHeight, regularFont, 10, alignments, false)
		currentY -= rowHeight
	}

	if err := page.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate PDF page: " + err.Error()})
		return
	}
	if err := doc.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not finalize PDF: " + err.Error()})
		return
	}
	log.Warn().Msg("Generating PDF report: customer-report1.pdf")
	// send PDF as HTTP response
	c.Header("Content-Disposition", `attachment; filename="customer-report1.pdf"`)
	c.Data(http.StatusOK, "application/pdf", buffer.Bytes())
}

func drawReportRow(page *document.Page, values []string, widths []float64, x, top, height float64, face font.Instance, fontSize float64, alignments []float64, header bool) {
	headerFill := color.DeviceRGB{30.0 / 255, 90.0 / 255, 150.0 / 255}
	bodyFill := color.DeviceRGB{1, 1, 1}
	border := color.DeviceRGB{0, 0, 0}
	textColor := color.DeviceRGB{0, 0, 0}
	if header {
		textColor = color.DeviceRGB{1, 1, 1}
	}

	page.SetLineWidth(0.5)
	for i, value := range values {
		fill := bodyFill
		if header {
			fill = headerFill
		}
		page.SetFillColor(fill)
		page.SetStrokeColor(border)
		page.Rectangle(x, top-height, widths[i], height)
		page.FillAndStroke()

		padding := mmToPoints(2)
		baseline := top - (height-fontSize)/2 - fontSize*0.75
		page.SetFillColor(textColor)
		drawReportText(page, face, fontSize, x+padding, baseline, widths[i]-2*padding, alignments[i], value)
		x += widths[i]
	}
}

func drawReportText(page *document.Page, face font.Instance, fontSize, x, y, width, alignment float64, text string) {
	page.TextSetFont(face, fontSize)
	page.TextBegin()
	page.TextFirstLine(x, y)
	page.TextShowAligned(text, width, alignment)
	page.TextEnd()
}

func mmToPoints(mm float64) float64 {
	return mm * 72 / 25.4
}
