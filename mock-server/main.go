package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// Models
type Vacancy struct {
	Organization   string    `json:"Organization"`
	Description    string    `json:"Description"`
	DateOfBegin    time.Time `json:"DateOfBegin"`
	DateOfEnd      time.Time `json:"DateOfEnd"`
	Salary         int       `json:"Salary"`
	Title          string    `json:"Title"`
	DateOfDocument time.Time `json:"DateOfDocument"`
	TypesOfWork    []string  `json:"TypesOfWork"`
	Number         string    `json:"Number"`
}

type Request struct {
	Organization string `json:"Organization"`
	Student      string `json:"Student"`
	Description  string `json:"Description"`
	StartPeriod  string `json:"StartPeriod"`
	EndPeriod    string `json:"EndPeriod"`
	Number       string `json:"Number"`
	Good         bool   `json:"Good"`
}

type Notify struct {
	Text            string    `json:"Text"`
	Date            time.Time `json:"Date"`
	NumberOfRequest string    `json:"NumberOfRequest"`
}

type Account struct {
	Organization string `json:"Organization"`
	Student      string `json:"Student"`
}

// Mock data
var vacancies = []Vacancy{
	{
		Organization:   "Волонтеры ДВФУ",
		Description:    "Необходимо доставлять гуманитарную помощь, покупать лекарства для пожилых немобильных людей",
		DateOfBegin:    parseDate("20260101"),
		DateOfEnd:      parseDate("20270101"),
		Salary:         50000,
		Title:          "Волонтер",
		DateOfDocument: time.Date(2026, 1, 5, 19, 40, 47, 0, time.UTC),
		TypesOfWork:    []string{"Помощь пожилым", "Медицина"},
		Number:         "000000004",
	},
	{
		Organization:   "Волонтеры ДВФУ",
		Description:    "Сбор мусора на набережной",
		DateOfBegin:    parseDate("20260101"),
		DateOfEnd:      parseDate("20270101"),
		Salary:         1000,
		Title:          "Волонтер",
		DateOfDocument: time.Date(2026, 1, 5, 19, 41, 59, 0, time.UTC),
		TypesOfWork:    []string{"Общественная польза"},
		Number:         "000000005",
	},
	{
		Organization:   "CODE WORK",
		Description:    "Необходимо обучать программированию студентов 1-2 курсов",
		DateOfBegin:    parseDate("20260601"),
		DateOfEnd:      parseDate("20260701"),
		Salary:         10000,
		Title:          "Учитель по программированию на С++",
		DateOfDocument: time.Date(2026, 1, 6, 10, 30, 0, 0, time.UTC),
		TypesOfWork:    []string{"Программирование", "Обучение", "C++"},
		Number:         "000000001",
	},
	{
		Organization:   "Tech Startup",
		Description:    "Разработка backend API на Go, опыт обязателен",
		DateOfBegin:    parseDate("20260215"),
		DateOfEnd:      parseDate("20260630"),
		Salary:         150000,
		Title:          "Go разработчик",
		DateOfDocument: time.Date(2026, 1, 7, 14, 20, 0, 0, time.UTC),
		TypesOfWork:    []string{"Программирование", "Go", "API", "Backend"},
		Number:         "000000006",
	},
}

var requests = []Request{
	{
		Organization: "Волонтеры ДВФУ",
		Student:      "Тестов Тест Тестович",
		Description:  "Очень хочу попробовать поработать волонтером, но нет опыта, имею свой транспорт.",
		StartPeriod:  "26.01.2026 0:00:00",
		EndPeriod:    "06.02.2026 0:00:00",
		Number:       "000000004",
		Good:         false,
	},
	{
		Organization: "CODE WORK",
		Student:      "Иванов Иван Иванович",
		Description:  "Опыт преподавания 3 года, люблю работать со студентами",
		StartPeriod:  "01.06.2026 0:00:00",
		EndPeriod:    "01.07.2026 0:00:00",
		Number:       "000000002",
		Good:         true,
	},
	{
		Organization: "Tech Startup",
		Student:      "Петров Петр Петрович",
		Description:  "Разработчик с опытом 5 лет, знаю Go, PostgreSQL, Docker",
		StartPeriod:  "15.02.2026 0:00:00",
		EndPeriod:    "30.06.2026 0:00:00",
		Number:       "000000007",
		Good:         false,
	},
}

var notifies = []Notify{
	{
		Text:            "Уважаемый Николаев Николай Николаевич! \n Одобрена ваша заявка по вакансии на должность Учитель по программированию на С++. \n Сообщение от руководителя: Подходите в кабинет C315 14.01.2026 с 13 до 14",
		Date:            time.Date(2026, 1, 11, 0, 0, 0, 0, time.UTC),
		NumberOfRequest: "000000002",
	},
}

var tags = []string{"Наука", "Медицина", "Литература", "Технологии", "Творчество", "Программирование", "Алгоритмы", "ICPC", "Помощь пожилым", "Общественная польза", "Backend", "Go", "API", "Обучение", "C++"}

func parseDate(dateStr string) time.Time {
	year := dateStr[0:4]
	month := dateStr[4:6]
	day := dateStr[6:8]
	t, _ := time.Parse("2006-01-02", fmt.Sprintf("%s-%s-%s", year, month, day))
	return t
}

// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Logger
func logRequest(r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(strings.NewReader(string(body)))

	fmt.Printf("\n%s %s\n", r.Method, r.URL.Path)
	if r.URL.RawQuery != "" {
		fmt.Printf("Query: %s\n", r.URL.RawQuery)
	}
	if len(body) > 0 {
		fmt.Printf("Body: %s\n", string(body))
	}
}

// 1. Create Vacancy - POST /JobService/hs/jobservice/vacancy
func createVacancy(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	var data map[string]interface{}
	json.NewDecoder(r.Body).Decode(&data)

	fmt.Println("✓ Вакансия создана")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// 2. Create Request - POST /JobService/hs/jobservice/request
func createRequest(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	var data map[string]interface{}
	json.NewDecoder(r.Body).Decode(&data)

	fmt.Println("✓ Заявка на работу создана")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// 3. Get Vacancy List - GET /JobService/hs/jobservice/vacancylist
func getVacancyList(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	salaryMin := r.URL.Query().Get("salaryMIN")
	typeOfWork := r.URL.Query().Get("typesofwork")
	organization := r.URL.Query().Get("organization")

	fmt.Printf("Фильтры: salaryMIN=%s, typesofwork=%s, organization=%s\n", salaryMin, typeOfWork, organization)

	filtered := vacancies
	fmt.Printf("✓ Возвращены вакансии: %d шт.\n", len(filtered))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(filtered)
}

// 4. Get Tags - GET /JobService/hs/jobservice/tags
func getTags(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	fmt.Printf("✓ Возвращены теги: %d шт.\n", len(tags))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tags)
}

// 5. Get Request List - GET /JobService/hs/jobservice/requestlist
func getRequestList(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	vacancy := r.URL.Query().Get("vacancy")
	fmt.Printf("Вакансия: %s\n", vacancy)

	result := []interface{}{
		map[string]int{"count": len(requests)},
	}
	result = append(result, requests)

	fmt.Printf("✓ Возвращены отклики: %d шт.\n", len(requests))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// 6. Check Account - GET /JobService/hs/jobservice/checkaccount
func checkAccount(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	user := r.URL.Query().Get("user")
	fmt.Printf("Пользователь: %s\n", user)

	var response Account
	if user == "ivanov.ii" {
		response = Account{Organization: "", Student: "123-694-775 67"}
	} else if user == "ivanov.iv" {
		response = Account{Organization: "f2742040-cdb4-11f0-ae42-38d57ae2c1c1", Student: ""}
	} else {
		response = Account{Organization: "", Student: ""}
	}

	fmt.Printf("✓ Аккаунт найден: %v\n", response)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// 7. Send FAQ Suggestion - POST /JobService/hs/jobservice/faq
func sendFAQ(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	var data map[string]string
	json.NewDecoder(r.Body).Decode(&data)

	fmt.Printf("Предложение: %s\n", data["suggestion"])
	fmt.Println("✓ Предложение сохранено")

	w.WriteHeader(http.StatusOK)
}

// 8. Apply Request - POST /JobService/hs/jobservice/applyrequest
func applyRequest(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	var data map[string]string
	json.NewDecoder(r.Body).Decode(&data)

	fmt.Printf("Номер отклика: %s, Сообщение: %s\n", data["number"], data["text"])
	fmt.Println("✓ Отклик одобрен")

	w.WriteHeader(http.StatusOK)
}

// 9. Get Notifications - GET /JobService/hs/jobservice/mynotify
func getNotifications(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	student := r.URL.Query().Get("student")
	fmt.Printf("СНИЛС студента: %s\n", student)

	fmt.Printf("✓ Возвращены уведомления: %d шт.\n", len(notifies))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(notifies)
}

// 10. Get Vacancy From Notify - GET /JobService/hs/jobservice/vacancyfromnotify
func getVacancyFromNotify(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	numberOfRequest := r.URL.Query().Get("numberofrequest")
	fmt.Printf("Номер отклика: %s\n", numberOfRequest)

	var result []Vacancy
	for _, v := range vacancies {
		if v.Number == "000000001" {
			result = append(result, v)
			break
		}
	}

	fmt.Printf("✓ Вакансия найдена\n")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// 11. Close Vacancy - POST /JobService/hs/jobservice/closevacancy
func closeVacancy(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	number := r.URL.Query().Get("number")
	fmt.Printf("Вакансия закрыта: %s\n", number)
	fmt.Println("✓ Вакансия удалена из списка")

	w.WriteHeader(http.StatusOK)
}

func main() {
	fmt.Println("🚀 WhiteMustache Mock Server запущен")
	fmt.Println("📍 http://localhost:80")
	fmt.Println("─────────────────────────────────────\n")

	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/JobService/hs/jobservice/vacancy", createVacancy)
	mux.HandleFunc("/JobService/hs/jobservice/request", createRequest)
	mux.HandleFunc("/JobService/hs/jobservice/vacancylist/", getVacancyList)
	mux.HandleFunc("/JobService/hs/jobservice/tags", getTags)
	mux.HandleFunc("/JobService/hs/jobservice/requestlist/", getRequestList)
	mux.HandleFunc("/JobService/hs/jobservice/checkaccount/", checkAccount)
	mux.HandleFunc("/JobService/hs/jobservice/faq", sendFAQ)
	mux.HandleFunc("/JobService/hs/jobservice/applyrequest", applyRequest)
	mux.HandleFunc("/JobService/hs/jobservice/mynotify", getNotifications)
	mux.HandleFunc("/JobService/hs/jobservice/vacancyfromnotify", getVacancyFromNotify)
	mux.HandleFunc("/JobService/hs/jobservice/closevacancy", closeVacancy)

	handler := corsMiddleware(mux)

	log.Fatal(http.ListenAndServe(":80", handler))
}
