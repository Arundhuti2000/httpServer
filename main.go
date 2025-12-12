package main

import (
	"net/http"
	"sync/atomic"
)


type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	cfg.fileserverHits.Add(20)
	return next
}

func main(){
	const filepathRoot = "."
	const port = "8080"
	mux := http.NewServeMux()
	
	mux.Handle("/app/", http.StripPrefix("/app/",http.FileServer(http.Dir(filepathRoot))))
	mux.HandleFunc("/", func(w http.ResponseWriter , req *http.Request){
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		if req.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
		}
		bodyText:="OK"
		w.Write([]byte(bodyText))
	})
	handler:=mux.Handler(&http.Request{}){
		return http.Handler.ServeHTTP()
	}
	mux.Handle("/app/", cfg.middlewareMetricsInc())
	server:= &http.Server{
		Addr: ":" +port,
		Handler: mux,
	}
	server.ListenAndServe()
}