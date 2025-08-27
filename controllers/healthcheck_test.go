package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthController_Ping_Success(t *testing.T) {
	controller := NewHealthController()
	router := setupTestRouter()

	router.GET("/ping", controller.Ping)

	req, _ := http.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	expectedBody := `"pong"`
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, w.Body.String())
	}
}

func TestHealthController_Ping_Method(t *testing.T) {
	controller := NewHealthController()
	router := setupTestRouter()

	router.GET("/ping", controller.Ping)

	// Test with POST method (should return 404)
	req, _ := http.NewRequest("POST", "/ping", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d for POST method, got %d", http.StatusNotFound, w.Code)
	}
}

func TestHealthController_NewHealthController(t *testing.T) {
	controller := NewHealthController()

	if controller == nil {
		t.Error("Expected non-nil controller")
	}
}
