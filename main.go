package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Arundhuti2000/httpserver/internal/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)


type apiConfig struct {
	fileserverHits atomic.Int32
	DB *database.Queries
}



func main(){
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err !=nil{
		fmt.Sprintf(err.Error())
	}
	dbQueries := database.New(db)
	const filepathRoot = "."
	const port = "8080"
	mux := http.NewServeMux()
	cfg := &apiConfig{
		fileserverHits: atomic.Int32{},
		DB: dbQueries,
	}
	fileserverhandler:=cfg.middlewareMetricsInc(http.StripPrefix("/app/",http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fileserverhandler)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	mux.HandleFunc("POST /api/validate_chirp", cfg.handlerValidateChirps)
	

	// handler:=mux.Handler(&http.Request{}){
	// 	return http.Handler.ServeHTTP()
	// }
	server:= &http.Server{
		Addr: ":" +port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}