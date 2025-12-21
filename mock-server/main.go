package main

import (
	"fmt"
	"net/http"
	"strings"
)

const soapResponse = `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:job="https://JobService.org">
    <soap:Body>
        <Vacancy>
            <Organization>ДВФУ</Organization>
            <Number>Стажер Go-разработчик</Number>
            <Description>Разработка backend сервисов на Go. Работа с gRPC и microservices. Оклад 25000₽/мес</Description>
            <DateOfBegin>2025-12-22</DateOfBegin>
            <DateOfEnd>2026-06-22</DateOfEnd>
            <TypesOfWork>Стажировка (40 ч/неделю)</TypesOfWork>
        </Vacancy>
        <Vacancy>
            <Organization>ДВФУ IT</Organization>
            <Number>Junior DevOps Engineer</Number>
            <Description>Работа с Docker, Kubernetes, CI/CD. Настройка инфраструктуры. Оклад 30000₽/мес</Description>
            <DateOfBegin>2025-12-22</DateOfBegin>
            <DateOfEnd>2026-08-22</DateOfEnd>
            <TypesOfWork>Полная занятость</TypesOfWork>
        </Vacancy>
        <Vacancy>
            <Organization>ДВФУ Lab</Organization>
            <Number>Backend Developer PostgreSQL</Number>
            <Description>Оптимизация БД, миграции, API на Go. Работа с Docker. Оклад 28000₽/мес</Description>
            <DateOfBegin>2025-12-22</DateOfBegin>
            <DateOfEnd>2026-12-22</DateOfEnd>
            <TypesOfWork>Подработка (20 ч/неделю)</TypesOfWork>
        </Vacancy>
    </soap:Body>
</soap:Envelope>`

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, SOAPAction")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func handleJobService(w http.ResponseWriter, r *http.Request) {
	body := make([]byte, r.ContentLength)
	r.Body.Read(body)
	bodyStr := string(body)

	if strings.Contains(bodyStr, "GetVacancy") {
		w.Header().Set("Content-Type", "text/xml; charset=UTF-8")
		fmt.Fprint(w, soapResponse)
	} else {
		http.Error(w, "Unknown request", http.StatusBadRequest)
	}
}

func main() {
	http.HandleFunc("/Job/ws/JobService.1cws", corsMiddleware(handleJobService))
	fmt.Println("SOAP сервер на http://localhost:80/Job/ws/JobService.1cws")
	http.ListenAndServe(":80", nil)
}
