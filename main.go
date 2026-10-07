// rendimiento-demo es una página pequeña para probar rendimiento.ai desde cero.
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

var (
	version = "dev" // al construir, la confirmación (GIT_SHA)
	started = time.Now()
	visitas atomic.Int64
)

var page = template.Must(template.New("page").Parse(`<!doctype html>
<html lang="es-MX"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>rendimiento-demo</title>
<style>
  :root { --fondo: #fdf6e9; --texto: #2b2118; --suave: #7a6a5a; --cempasuchil: #f29f05; --grana: #b3261e; --tarjeta: #fffdf8; }
  @media (prefers-color-scheme: dark) { :root { --fondo: #14110f; --texto: #f3ece2; --suave: #b3a597; --tarjeta: #1f1a16; } }
  * { box-sizing: border-box; }
  body { margin: 0; min-height: 100vh; display: grid; place-items: center; padding: 16px;
         font-family: system-ui, sans-serif; background: var(--fondo); color: var(--texto); }
  .picado { position: fixed; top: 0; left: 0; right: 0; height: 14px;
            background: repeating-linear-gradient(90deg, var(--grana) 0 40px, var(--cempasuchil) 40px 80px, #2f6f9f 80px 120px, #d6457f 120px 160px); }
  main { max-width: 520px; width: 100%; text-align: center; background: var(--tarjeta); border-radius: 18px;
         padding: 32px 24px; box-shadow: 0 10px 30px rgb(0 0 0 / .08); }
  h1 { font-size: 2.2rem; margin: 12px 0 6px; }
  p { color: var(--suave); margin: 6px 0; }
  code { background: rgb(242 159 5 / .15); padding: 2px 8px; border-radius: 6px; color: var(--texto); }
  .flor { font-size: 3rem; display: inline-block; animation: latido 1.6s ease-in-out infinite; }
  @keyframes latido { 0%, 100% { transform: scale(1); } 15% { transform: scale(1.12); } 30% { transform: scale(1); } }
  @media (prefers-reduced-motion: reduce) { .flor { animation: none; } }
</style></head>
<body><div class="picado" aria-hidden="true"></div>
<main>
  <span class="flor" aria-hidden="true">🌼</span>
  <h1>{{.Saludo}}</h1>
  <p>Desplegada por <b>rendimiento</b> · versión <code>{{.Version}}</code></p>
  <p>servidor <code>{{.Host}}</code> · encendida hace {{.Encendida}} · visita n.º {{.Visitas}}</p>
</main></body></html>`))

func saludo() string {
	if s := os.Getenv("SALUDO"); s != "" {
		return s
	}
	return "¡Hola desde rendimiento-demo!"
}

func handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		host, _ := os.Hostname()
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = page.Execute(w, map[string]any{
			"Saludo":    saludo(),
			"Version":   version,
			"Host":      host,
			"Encendida": time.Since(started).Round(time.Second),
			"Visitas":   visitas.Add(1),
		})
	})
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	mux.HandleFunc("GET /api/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"saludo": saludo(), "version": version})
	})
	return mux
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("rendimiento-demo %s escuchando en :%s", version, port)
	srv := &http.Server{Addr: ":" + port, Handler: handler(), ReadHeaderTimeout: 5 * time.Second}
	log.Fatal(srv.ListenAndServe())
}
