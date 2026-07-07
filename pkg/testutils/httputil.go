package testutils

import (
	"bytes"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

// NewTestContext builds a gin test context with a request body.
func NewTestContext(method, path string, body []byte) (*gin.Engine, *gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	r := gin.New()
	ctx, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx.Request = req
	return r, ctx, w
}
