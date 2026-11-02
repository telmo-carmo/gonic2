package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestReport1HandlerProducesPDF(t *testing.T) {
	oldMode := gin.Mode()
	gin.SetMode(gin.TestMode)
	defer gin.SetMode(oldMode)

	customers := make([]Customer, 100)
	for i := range customers {
		customers[i] = Customer{
			ID:      i + 1,
			Name:    "José Silva",
			Country: "Portugal",
			Balance: 15220.50,
		}
	}
	body, err := json.Marshal(ReportRequest{
		Title:     "Customer Report",
		Customers: customers,
	})
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/reports", bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")

	Report1Handler(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; response: %s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	if got := recorder.Header().Get("Content-Type"); got != "application/pdf" {
		t.Fatalf("Content-Type = %q, want application/pdf", got)
	}
	if got := recorder.Header().Get("Content-Disposition"); got != `attachment; filename="customer-report1.pdf"` {
		t.Fatalf("Content-Disposition = %q", got)
	}
	if !strings.HasPrefix(recorder.Body.String(), "%PDF-2.0") {
		t.Fatal("response does not start with the expected PDF header")
	}
}
