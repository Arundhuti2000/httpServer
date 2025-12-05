package main

import "net/http"

func main(){
	mux := http.NewServeMux()
	mux.Handle("/",)
	server:= &http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}