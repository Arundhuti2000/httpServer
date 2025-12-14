package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)


type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w,r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request){
	fmt.Println("metrics handler called")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	hits:=cfg.fileserverHits.Load()
	w.Write([]byte(fmt.Sprintf("Hits: %d", hits)))
}
		


func main(){
	const filepathRoot = "."
	const port = "8080"
	mux := http.NewServeMux()
	cfg := &apiConfig{
		fileserverHits: atomic.Int32{},
	}
	
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/",http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /healthz", handlerReadiness)
	mux.HandleFunc("GET /metrics", cfg.handlerMetrics)

	mux.HandleFunc("POST /reset", cfg.handlerReset)
	

	// handler:=mux.Handler(&http.Request{}){
	// 	return http.Handler.ServeHTTP()
	// }
	server:= &http.Server{
		Addr: ":" +port,
		Handler: mux,
	}
	server.ListenAndServe()
}