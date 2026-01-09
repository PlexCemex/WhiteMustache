package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ====== Модели ======

type Vacancy struct {
	Organization   string   `json:"Organization"`
	Description    string   `json:"Description"`
	DateOfBegin    string   `json:"DateOfBegin"`
	DateOfEnd      string   `json:"DateOfEnd"`
	Salary         int      `json:"Salary"`
	Title          string   `json:"Title"`
	DateOfDocument string   `json:"DateOfDocument"`
	TypesOfWork    []string `json:"TypesOfWork"`
	Number         string   `json:"Number"`
}

type VacancyCreateRequest struct {
	Salary       string `json:"salary"`
	Title        string `json:"title"`
	DateOfBegin  string `json:"dateofbegin"`
	DateOfEnd    string `json:"dateofend"`
	Organization string `json:"organization"`
	Description  string `json:"description"`
	TypesOfWork  string `json:"typesofwork"`
}

type RequestCreateInput struct {
	StartPeriod string `json:"startperiod"`
	EndPeriod   string `json:"endperiod"`
	Student     string `json:"student"`
	Description string `json:"description"`
	Vacancy     string `json:"vacancy"`
}

type WorkRequest struct {
	VacancyNumber string `json:"-"`
	Organization  string `json:"Organization"`
	Student       string `json:"Student"`
	Description   string `json:"Description"`
	StartPeriod   string `json:"StartPeriod"`
	EndPeriod     string `json:"EndPeriod"`
}

type CheckAccountInput struct {
	User string `json:"user"`
}

type CheckAccountResponse struct {
	Organization string `json:"Организация"`
}

// ====== Глобальные данные (моки) ======

var (
	mu        sync.RWMutex
	vacancies []Vacancy
	requests  []WorkRequest
	tags      []string

	organizations = map[string]string{
		"49ca8044-cdac-11f0-ae42-38d57ae2c1c1": "Студенческий отряд",
		"f2742040-cdb4-11f0-ae42-38d57ae2c1c1": "Волонтеры ДВФУ",
		"4c09ed30-cdb6-11f0-ae42-38d57ae2c1c1": "ЦПД",
		"76fa74fd-ea4f-11f0-ae6f-38d57ae2c1c1": "CODE WORK",
	}
)

func initData() {
	tags = []string{
		"Наука", "Медицина", "Литература", "Технологии", "Творчество",
		"Программирование", "Алгоритмы", "ICPC", "Помощь пожилым",
		"Общественная польза", "Робототехника", "Кибербезопасность",
		"Маркетинг", "Дизайн", "Финансы",
	}

	// несколько вакансий, включая те, что в примере документа
	vacancies = []Vacancy{
		{
			Organization:   "Волонтеры ДВФУ",
			Description:    "Необходимо доставлять гуманитарную помощь, покупать лекарства для пожилых немобильных людей",
			DateOfBegin:    "2026-01-01T00:00:00",
			DateOfEnd:      "2027-01-01T00:00:00",
			Salary:         50000,
			Title:          "Волонтер",
			DateOfDocument: "2026-01-05T19:40:47",
			TypesOfWork:    []string{"Помощь пожилым", "Медицина"},
			Number:         "000000004",
		},
		{
			Organization:   "Волонтеры ДВФУ",
			Description:    "Сбор мусора на набережной",
			DateOfBegin:    "2026-01-01T00:00:00",
			DateOfEnd:      "2027-01-01T00:00:00",
			Salary:         1000,
			Title:          "Волонтер",
			DateOfDocument: "2026-01-05T19:41:59",
			TypesOfWork:    []string{"Общественная польза"},
			Number:         "000000005",
		},
		{
			Organization:   "ЦПД",
			Description:    "Разработка ПО для прототипа робота-уборщика",
			DateOfBegin:    "2026-01-10T00:00:00",
			DateOfEnd:      "2027-01-10T00:00:00",
			Salary:         120000,
			Title:          "Программист микроконтроллеров",
			DateOfDocument: "2026-01-06T10:15:00",
			TypesOfWork:    []string{"Наука", "Техника", "Программирование", "Алгоритмы", "Робототехника"},
			Number:         "000000001",
		},
		{
			Organization:   "CODE WORK",
			Description:    "Разработка личного кабинета студента (backend на Go)",
			DateOfBegin:    "2026-02-01T00:00:00",
			DateOfEnd:      "2026-08-01T00:00:00",
			Salary:         80000,
			Title:          "Backend-разработчик (Go)",
			DateOfDocument: "2026-01-10T12:00:00",
			TypesOfWork:    []string{"Программирование", "Алгоритмы", "ICPC", "Технологии"},
			Number:         "000000002",
		},
		{
			Organization:   "Студенческий отряд",
			Description:    "Организация и проведение мероприятий для школьников",
			DateOfBegin:    "2026-03-01T00:00:00",
			DateOfEnd:      "2026-06-01T00:00:00",
			Salary:         25000,
			Title:          "Организатор мероприятий",
			DateOfDocument: "2026-01-12T09:30:00",
			TypesOfWork:    []string{"Общественная польза", "Творчество"},
			Number:         "000000003",
		},
		{
			Organization:   "CODE WORK",
			Description:    "Разработка REST API для сервиса подработок",
			DateOfBegin:    "2026-01-20T00:00:00",
			DateOfEnd:      "2026-09-01T00:00:00",
			Salary:         90000,
			Title:          "Разработчик API (Go)",
			DateOfDocument: "2026-01-15T15:45:00",
			TypesOfWork:    []string{"Программирование", "Технологии"},
			Number:         "000000006",
		},
		{
			Organization:   "ЦПД",
			Description:    "Стажировка по анализу данных студентов",
			DateOfBegin:    "2026-04-01T00:00:00",
			DateOfEnd:      "2026-12-01T00:00:00",
			Salary:         60000,
			Title:          "Data Analyst Intern",
			DateOfDocument: "2026-01-18T11:20:00",
			TypesOfWork:    []string{"Наука", "Технологии", "Финансы"},
			Number:         "000000007",
		},
		{
			Organization:   "Волонтеры ДВФУ",
			Description:    "Помощь в организации благотворительного марафона",
			DateOfBegin:    "2026-05-01T00:00:00",
			DateOfEnd:      "2026-06-01T00:00:00",
			Salary:         0,
			Title:          "Волонтер-организатор",
			DateOfDocument: "2026-01-20T08:30:00",
			TypesOfWork:    []string{"Общественная польза", "Творчество"},
			Number:         "000000008",
		},
		{
			Organization:   "CODE WORK",
			Description:    "Разработка бота-помощника для студентов",
			DateOfBegin:    "2026-01-25T00:00:00",
			DateOfEnd:      "2026-11-01T00:00:00",
			Salary:         110000,
			Title:          "Разработчик чат-ботов",
			DateOfDocument: "2026-01-22T14:10:00",
			TypesOfWork:    []string{"Программирование", "Технологии", "Алгоритмы"},
			Number:         "000000009",
		},
		{
			Organization:   "Студенческий отряд",
			Description:    "Настройка и обслуживание компьютерных классов",
			DateOfBegin:    "2026-02-10T00:00:00",
			DateOfEnd:      "2026-07-10T00:00:00",
			Salary:         40000,
			Title:          "IT-специалист",
			DateOfDocument: "2026-01-25T16:00:00",
			TypesOfWork:    []string{"Технологии", "Программирование"},
			Number:         "000000010",
		},
	}

	// стартовые отклики (как в примере + немного своих)
	requests = []WorkRequest{
		{
			VacancyNumber: "000000001",
			Organization:  "CODE WORK",
			Student:       "Иванов Иван Иванович",
			Description:   "Знаю С++ на уровне middle, опыт программирования 6 лет",
			StartPeriod:   "01.01.0001 00:00:00",
			EndPeriod:     "01.01.0001 00:00:00",
		},
		{
			VacancyNumber: "000000001",
			Organization:  "CODE WORK",
			Student:       "Николаев Николай Николаевич",
			Description:   "Изучал С++ месяц",
			StartPeriod:   "02.06.2026 00:00:00",
			EndPeriod:     "01.07.2026 00:00:00",
		},
		{
			VacancyNumber: "000000004",
			Organization:  "Волонтеры ДВФУ",
			Student:       "Петров Петр Петрович",
			Description:   "Есть опыт волонтерской деятельности в хосписе",
			StartPeriod:   "10.02.2026 00:00:00",
			EndPeriod:     "10.05.2026 00:00:00",
		},
	}
}

// ====== Утилиты ======

func formatAPIDateToISO(s string) string {
	// "20260101" -> "2026-01-01T00:00:00"
	t, err := time.Parse("20060102", s)
	if err != nil {
		return ""
	}
	return t.Format("2006-01-02T15:04:05")
}

func formatTo1CDate(s string) string {
	// "20260602" -> "02.06.2026 00:00:00"
	t, err := time.Parse("20060102", s)
	if err != nil {
		return s
	}
	return t.Format("02.01.2006 15:04:05")
}

func nextVacancyNumber() string {
	return fmtNumber(len(vacancies) + 1)
}

func fmtNumber(n int) string {
	return fmtWithZeroes(n, 6)
}

func fmtWithZeroes(n, width int) string {
	s := strconv.Itoa(n)
	if len(s) >= width {
		return s
	}
	return strings.Repeat("0", width-len(s)) + s
}

// ====== Middleware: CORS + Логирование ======

func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}

func withLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("---- %s %s from %s ----", r.Method, r.URL.String(), r.RemoteAddr)

		if r.Body != nil {
			bodyBytes, _ := io.ReadAll(r.Body)
			if len(bodyBytes) > 0 {
				ct := r.Header.Get("Content-Type")
				if strings.Contains(ct, "application/json") {
					var pretty bytes.Buffer
					if err := json.Indent(&pretty, bodyBytes, "", "  "); err == nil {
						log.Printf("Request JSON body:\n%s", pretty.String())
					} else {
						log.Printf("Request body (raw): %s", string(bodyBytes))
					}
				} else {
					log.Printf("Request body: %s", string(bodyBytes))
				}
				r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			}
		}

		next(w, r)
	}
}

// ====== Handlers ======

func handleVacancyCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	var in VacancyCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	salary, err := strconv.Atoi(strings.TrimSpace(in.Salary))
	if err != nil {
		salary = 0
	}

	types := []string{}
	for _, t := range strings.Split(in.TypesOfWork, ",") {
		t = strings.TrimSpace(t)
		if t != "" {
			types = append(types, t)
		}
	}

	now := time.Now().Format("2006-01-02T15:04:05")

	v := Vacancy{
		Organization:   resolveOrgName(in.Organization),
		Description:    in.Description,
		DateOfBegin:    formatAPIDateToISO(in.DateOfBegin),
		DateOfEnd:      formatAPIDateToISO(in.DateOfEnd),
		Salary:         salary,
		Title:          in.Title,
		DateOfDocument: now,
		TypesOfWork:    types,
		Number:         nextVacancyNumber(),
	}

	mu.Lock()
	vacancies = append(vacancies, v)
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(v)
}

func resolveOrgName(id string) string {
	if name, ok := organizations[id]; ok {
		return name
	}
	return id
}

func handleVacancyList(w http.ResponseWriter, r *http.Request) {
	// фильтры: typesofwork, organization, salaryMIN (остальное можно добавить по необходимости)
	q := r.URL.Query()
	filterTypes := q.Get("typesofwork")
	filterOrg := q.Get("organization")
	filterSalaryMinStr := q.Get("salaryMIN")

	var filterTypesSet map[string]struct{}
	if filterTypes != "" {
		filterTypesSet = make(map[string]struct{})
		for _, t := range strings.Split(filterTypes, ",") {
			t = strings.TrimSpace(t)
			if t != "" {
				filterTypesSet[t] = struct{}{}
			}
		}
	}

	var filterOrgsSet map[string]struct{}
	if filterOrg != "" {
		filterOrgsSet = make(map[string]struct{})
		for _, o := range strings.Split(filterOrg, ",") {
			o = strings.TrimSpace(o)
			if o != "" {
				// фильтруем по ID организаций, но в Vacancy у нас имя,
				// поэтому сравним с organizations[id] по имени
				if name, ok := organizations[o]; ok {
					filterOrgsSet[name] = struct{}{}
				}
			}
		}
	}

	var salaryMin int
	if filterSalaryMinStr != "" {
		if v, err := strconv.Atoi(filterSalaryMinStr); err == nil {
			salaryMin = v
		}
	}

	mu.RLock()
	defer mu.RUnlock()

	var result []Vacancy
	for _, v := range vacancies {
		if salaryMin > 0 && v.Salary < salaryMin {
			continue
		}
		if filterOrgsSet != nil {
			if _, ok := filterOrgsSet[v.Organization]; !ok {
				continue
			}
		}
		if filterTypesSet != nil {
			ok := false
			for _, t := range v.TypesOfWork {
				if _, found := filterTypesSet[t]; found {
					ok = true
					break
				}
			}
			if !ok {
				continue
			}
		}
		result = append(result, v)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(result)
}

func handleTags(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(tags)
}

func handleRequestCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	var in RequestCreateInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	orgName := "Неизвестная организация"
	mu.RLock()
	for _, v := range vacancies {
		if v.Number == in.Vacancy {
			orgName = v.Organization
			break
		}
	}
	mu.RUnlock()

	wr := WorkRequest{
		VacancyNumber: in.Vacancy,
		Organization:  orgName,
		Student:       in.Student,
		Description:   in.Description,
		StartPeriod:   formatTo1CDate(in.StartPeriod),
		EndPeriod:     formatTo1CDate(in.EndPeriod),
	}

	mu.Lock()
	requests = append(requests, wr)
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "ok",
		"message":  "request created",
		"vacancy":  in.Vacancy,
		"student":  in.Student,
		"org":      orgName,
		"start":    wr.StartPeriod,
		"end":      wr.EndPeriod,
		"datetime": time.Now().Format(time.RFC3339),
	})
}

func handleRequestList(w http.ResponseWriter, r *http.Request) {
	// по документации: GET + JSON-тело { "vacancy": "000000001" }
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "only GET/POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		Vacancy string `json:"vacancy"`
	}

	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&body)
	}
	if body.Vacancy == "" {
		// можно также брать из query ?vacancy=...
		body.Vacancy = r.URL.Query().Get("vacancy")
	}

	mu.RLock()
	defer mu.RUnlock()

	var filtered []WorkRequest
	for _, wr := range requests {
		if body.Vacancy == "" || wr.VacancyNumber == body.Vacancy {
			filtered = append(filtered, wr)
		}
	}

	resp := make([]interface{}, 0, len(filtered)+1)
	resp = append(resp, map[string]int{"count": len(filtered)})
	for _, wr := range filtered {
		resp = append(resp, wr)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(resp)
}

func handleCheckAccount(w http.ResponseWriter, r *http.Request) {
	// дока: GET + JSON { "user": "ivanov.iv" }
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "only GET/POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var in CheckAccountInput
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&in)
	}
	if in.User == "" {
		in.User = r.URL.Query().Get("user")
	}

	resp := CheckAccountResponse{Organization: ""}

	switch in.User {
	case "ivanov.ii":
		resp.Organization = ""
	case "ivanov.iv":
		resp.Organization = "f2742040-cdb4-11f0-ae42-38d57ae2c1c1"
	default:
		// условно считаем всех остальных студентами без организации
		resp.Organization = ""
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(resp)
}

// ====== main ======

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	initData()

	http.HandleFunc("/JobService/hs/jobservice/vacancy",
		withLogging(withCORS(handleVacancyCreate)))
	http.HandleFunc("/JobService/hs/jobservice/vacancylist/",
		withLogging(withCORS(handleVacancyList)))
	http.HandleFunc("/JobService/hs/jobservice/tags",
		withLogging(withCORS(handleTags)))
	http.HandleFunc("/JobService/hs/jobservice/request",
		withLogging(withCORS(handleRequestCreate)))
	http.HandleFunc("/JobService/hs/jobservice/requestlist",
		withLogging(withCORS(handleRequestList)))
	http.HandleFunc("/JobService/hs/jobservice/checkaccount",
		withLogging(withCORS(handleCheckAccount)))

	// http.HandleFunc("/JobService/hs/jobservice/vacancy",
	// 	withLogging((handleVacancyCreate)))
	// http.HandleFunc("/JobService/hs/jobservice/vacancylist/",
	// 	withLogging((handleVacancyList)))
	// http.HandleFunc("/JobService/hs/jobservice/tags",
	// 	withLogging((handleTags)))
	// http.HandleFunc("/JobService/hs/jobservice/request",
	// 	withLogging((handleRequestCreate)))
	// http.HandleFunc("/JobService/hs/jobservice/requestlist",
	// 	withLogging((handleRequestList)))
	// http.HandleFunc("/JobService/hs/jobservice/checkaccount",
	// 	withLogging((handleCheckAccount)))

	log.Println("Mock JobService listening on :80")
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}
