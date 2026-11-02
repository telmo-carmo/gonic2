package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"github.com/phuslu/log"

	"bytes"
	"fmt"
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
      "name": "John Smith",
      "country": "Portugal",
      "balance": 15220.50
    },
    {
      "id": 1002,
      "name": "Maria Silva",
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

	pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.SetMargins(15, 15, 15)
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 18)
	pdf.CellFormat(
		0,
		12,
		request.Title,
		"",
		1,
		"C",
		false,
		0,
		"",
	)

	pdf.Ln(5)

	// Table header
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

	pdf.SetFont("Arial", "B", 11)
	pdf.SetFillColor(30, 90, 150)
	pdf.SetTextColor(255, 255, 255)

	for i, header := range headers {
		pdf.CellFormat(
			widths[i],
			9,
			header,
			"1",
			0,
			"C",
			true,
			0,
			"",
		)
	}

	pdf.Ln(-1)

	// Table rows
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(0, 0, 0)

	for _, customer := range request.Customers {

		pdf.CellFormat(
			widths[0],
			8,
			fmt.Sprintf("%d", customer.ID),
			"1",
			0,
			"C",
			false,
			0,
			"",
		)

		pdf.CellFormat(
			widths[1],
			8,
			customer.Name,
			"1",
			0,
			"L",
			false,
			0,
			"",
		)

		pdf.CellFormat(
			widths[2],
			8,
			customer.Country,
			"1",
			0,
			"L",
			false,
			0,
			"",
		)

		pdf.CellFormat(
			widths[3],
			8,
			fmt.Sprintf("%.2f EUR", customer.Balance),
			"1",
			0,
			"R",
			false,
			0,
			"",
		)

		pdf.Ln(-1)
	}

	// Generate into memory
	var buffer bytes.Buffer

	if err := pdf.Output(&buffer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not generate PDF",
		})
		return
	}
	log.Warn().Msg("Generating PDF report: customer-report1.pdf")
	// HTTP response
	c.Header(
		"Content-Disposition",
		`attachment; filename="customer-report1.pdf"`,
	)

	c.Data(
		http.StatusOK,
		"application/pdf",
		buffer.Bytes(),
	)
}
