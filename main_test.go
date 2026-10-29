package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
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

func TestHace(t *testing.T) {
	casos := map[time.Duration]string{
		45 * time.Second:             "45 s",
		3 * time.Minute:              "3 min",
		2*time.Hour + 5*time.Minute:  "2 h 5 min",
		3*24*time.Hour + 4*time.Hour: "3 días",
	}
	for d, quiero := range casos {
		if got := hace(d); got != quiero {
			t.Errorf("hace(%v) = %q, quiero %q", d, got, quiero)
		}
	}
}
