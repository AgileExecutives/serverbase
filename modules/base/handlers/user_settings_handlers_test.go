package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AgileExecutives/serverbase/internal/models"
	baseRepo "github.com/AgileExecutives/serverbase/modules/base/repo"
	baseServices "github.com/AgileExecutives/serverbase/modules/base/services"
	"github.com/AgileExecutives/serverbase/pkg/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func setupRouterWithUserID(h *UserSettingsHandlers, userID uint) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	// middleware to set userID in context
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	r.GET("/user-settings", h.GetUserSettings)
	r.PUT("/user-settings", h.UpdateUserSettings)
	r.POST("/user-settings/reset", h.ResetUserSettings)
	return r
}

func TestUserSettingsHandlers_GetUpdateReset(t *testing.T) {
	repo := baseRepo.NewInMemoryUserSettingsRepo()
	svc := baseServices.NewUserSettingsService(repo)
	logger := testutils.NewMockLogger()
	h := NewUserSettingsHandlers(svc, logger)

	router := setupRouterWithUserID(h, 42)

	// GET -> should create defaults
	req := httptest.NewRequest(http.MethodGet, "/user-settings", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var resp models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.True(t, resp.Success)
	dataMap, ok := resp.Data.(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, float64(42), dataMap["user_id"].(float64))
	require.Equal(t, "en", dataMap["language"].(string))

	// PUT -> update language and theme
	payload := `{"language":"fr","theme":"dark","timezone":"CET","settings":"{\"k\":1}"}`
	req = httptest.NewRequest(http.MethodPut, "/user-settings", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.True(t, resp.Success)
	dataMap, ok = resp.Data.(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, "fr", dataMap["language"].(string))
	require.Equal(t, "dark", dataMap["theme"].(string))

	// POST reset -> back to defaults
	req = httptest.NewRequest(http.MethodPost, "/user-settings/reset", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.True(t, resp.Success)
	dataMap, ok = resp.Data.(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, "en", dataMap["language"].(string))
	require.Equal(t, "light", dataMap["theme"].(string))
}
