package main

import (
	p "net/http"
	"rockpaperscissor/http"
)

func main() {
	han := http.NewGameHandler()
	p.ListenAndServe(":3000", han)
	srv := http.NewServer(":3000", han)
	defer srv.Close()
	srv.Open()
}
