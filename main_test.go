package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func get(t *testing.T, path string) (int, string) {
	t.Helper()
	rec := httptest.NewRecorder()
	handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	return rec.Code, rec.Body.String()
}

func TestHealthz(t *testing.T) {
	if code, body := get(t, "/healthz"); code != 200 || strings.TrimSpace(body) != "ok" {
		t.Fatalf("healthz = %d %q", code, body)
	}
}

func TestPaginaUsaSaludo(t *testing.T) {
	t.Setenv("SALUDO", "Buenas tardes")
	if code, body := get(t, "/"); code != 200 || !strings.Contains(body, "Buenas tardes") {
		t.Fatalf("página = %d", code)
	}
}

func TestInfo(t *testing.T) {
	if code, body := get(t, "/api/info"); code != 200 || !strings.Contains(body, `"saludo"`) {
		t.Fatalf("info = %d %q", code, body)
	}
}

func TestRutaDesconocida(t *testing.T) {
	if code, _ := get(t, "/nada"); code != 404 {
		t.Fatalf("ruta desconocida = %d", code)
	}
}
