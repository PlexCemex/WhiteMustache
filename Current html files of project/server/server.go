package main

import (
	"log"
	"net/http"
)

func main() {
	// Настройка маршрутов
	http.HandleFunc("/", serveMain)
	http.HandleFunc("/main.html", serveMain)
	http.HandleFunc("/vacancy.html", serveVacancy)
	http.HandleFunc("/employer.html", serveEmployer)

	// Запуск сервера на порту 8080
	log.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveMain(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "main.html")
}

func serveVacancy(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "vacancy.html")
}

func serveEmployer(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "employer.html")
}
