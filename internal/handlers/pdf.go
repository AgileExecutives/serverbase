package handlers

import (
	"github.com/AgileExecutives/serverbase/pkg/core"
	pdfServices "github.com/AgileExecutives/shared-modules/pdf/services"
	"github.com/gin-gonic/gin"
)

// PdfHandler handles PDF generation requests and delegates to shared service.
type PdfHandler struct{ pdfService *pdfServices.PDFGenerator }

// NewPDFHandlerWithCtx constructs a PdfHandler using ModuleContext. If the
// pdf service is registered in ctx.Services it will be used; otherwise a
// default PDFGenerator is created.
func NewPDFHandlerWithCtx(ctx core.ModuleContext) *PdfHandler {
	if ctx.Services != nil {
		if s, ok := ctx.Services.Get("pdf-generator"); ok {
			if gen, ok := s.(*pdfServices.PDFGenerator); ok {
				return &PdfHandler{pdfService: gen}
			}
		}
	}
	// Fallback: simple default generator (no repo)
	return &PdfHandler{pdfService: pdfServices.NewPDFGenerator()}
}

// GeneratePDFFromTemplate delegates to the pdf service. Validation is still
// performed at the handler level to provide proper HTTP responses.
func (h *PdfHandler) GeneratePDFFromTemplate(c *gin.Context) {
	type PDFGenerateRequest struct {
		Data         map[string]interface{} `json:"data"`
		TemplateName string                 `json:"templateName"`
		FileName     string                 `json:"fileName"`
	}
	var req PDFGenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	if req.Data == nil || req.TemplateName == "" || req.FileName == "" {
		c.JSON(400, gin.H{"error": "Missing required fields"})
		return
	}
	name, err := h.pdfService.GeneratePDF(req.Data, req.TemplateName, req.FileName)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate PDF", "details": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "filename": name})
}

// Legacy constructor maintained for compatibility
func NewPDFHandler(pdfService *pdfServices.PDFGenerator) *PdfHandler {
	return &PdfHandler{pdfService: pdfService}
}
