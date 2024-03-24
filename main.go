package main

import (
	"apiperso/handler"
	"net/http"
	"os"
)

func main() {

	http.HandleFunc("/main", handler.HandlerMain)
	http.HandleFunc("/liste", handler.HandlerListe)
	http.HandleFunc("/info/", handler.HandlerInfo)
	http.HandleFunc("/categorie", handler.HandlerCategorie)
	http.HandleFunc("/deck", handler.HandlerDeck)
	http.HandleFunc("/deck/add/", handler.HandlerDeckAdd)
	http.HandleFunc("/deck/remove/", handler.HandlerDeckRemove)
	http.HandleFunc("/recherche", handler.HandlerRecherche)
	http.HandleFunc("/", handler.NotFoundHandler)
	http.HandleFunc("/404", handler.Handler404)
	http.HandleFunc("/aboutus", handler.HandlerAboutUs)

	rootDoc, _ := os.Getwd()
	fileserver := http.FileServer(http.Dir(rootDoc + "/assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fileserver))
	http.ListenAndServe(":8080", nil)
}
