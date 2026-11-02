package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/phuslu/log"
	"github.com/xuri/excelize/v2"
)

func Report2Handler(c *gin.Context) {
	// 1. Initialize a clean excelize file instance
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Error().Err(err).Msg("Error closing excelize file")
		}
	}()

	sheetName := "Report 2"
	f.SetSheetName("Sheet1", sheetName)

	// Mock Data (Replace with your database query or slice)
	headers := []string{"ID", "Product Name", "Price", "Created At"}
	rows := [][]interface{}{
		{1, "Premium Microservice Framework", 199.99, time.Now().Format("2006-01-02")},
		{2, "Go 1.26 Performance Tuning Guide", 49.50, time.Now().Format("2006-01-02")},
		{3, "Automated Cloud Sync Engine", 1250.00, "2026-03-02"},
		{4, "Next-Gen AI Toolkit", 13299.10, "2016-04-02"},
	}

	//  Write Headers
	for colNum, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(colNum+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// Write Data Rows
	for rowIdx, rowData := range rows {
		rowNum := rowIdx + 2 // Rows are 1-indexed; row 1 is headers
		for colIdx, value := range rowData {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowNum)
			f.SetCellValue(sheetName, cell, value)
		}
	}
	// generate in row 5 , col 3 a sum of previous prices
	f.SetCellFormula(sheetName, "C6", "SUM(C2:C5)")

	// Configure Content Headers for HTTP dynamic downloading
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=report2.xlsx")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	log.Info().Msg("Streaming Excel report: report2.xlsx")
	//  STREAM ON THE FLY: Write directly into Gin's HTTP response writer
	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stream Excel payload: " + err.Error()})
		return
	}
}
