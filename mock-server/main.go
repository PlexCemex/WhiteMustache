package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

type PutVacancyRequest struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Body    struct {
		PutVacancy struct {
			Organization string `xml:"Organization"`
			Description  string `xml:"Description"`
			DateOfBegin  string `xml:"DateOfBegin"`
			DateOfEnd    string `xml:"DateOfEnd"`
			TypesOfWork  string `xml:"TypesOfWork"`
		} `xml:"PutVacancy"`
	} `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
}


func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, SOAPAction")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func putVacancyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fmt.Println("=== Получен XML запрос ===")
	fmt.Println(string(body))
	fmt.Println("===========================")

	var req PutVacancyRequest
	if err := xml.Unmarshal(body, &req); err != nil {
		fmt.Printf("Ошибка парсинга: %v\n", err)
	}

	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	response := `<?xml version="1.0" encoding="UTF-8"?><Response><Status>OK</Status></Response>`
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, response)

	fmt.Println("✓ Ответ отправлен")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/Job/ws/JobService.1cws", putVacancyHandler)

	handler := corsMiddleware(mux)

	fmt.Println("🚀 Сервер запущен на http://localhost:80")
	if err := http.ListenAndServe(":80", handler); err != nil {
		fmt.Println("Ошибка:", err)
	}
}
