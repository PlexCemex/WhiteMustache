package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Experience struct {
	Title   string `json:"title"`
	Company string `json:"company"`
	Period  string `json:"period"`
}

type Education struct {
	Degree string `json:"degree"`
	School string `json:"school"`
	Year   string `json:"year"`
}

type Candidate struct {
	ID         int          `json:"id"`
	Name       string       `json:"name"`
	Position   string       `json:"position"`
	Location   string       `json:"location"`
	Email      string       `json:"email"`
	Phone      string       `json:"phone"`
	Avatar     string       `json:"avatar"`
	Rating     float64      `json:"rating"`
	About      string       `json:"about"`
	Skills     []string     `json:"skills"`
	Experience []Experience `json:"experience"`
	Education  []Education  `json:"education"`
}

type Job struct {
	ID           int      `json:"id"`
	Title        string   `json:"title"`
	Company      string   `json:"company"`
	Salary       string   `json:"salary"`
	Level        string   `json:"level"`
	Type         string   `json:"type"`
	Location     string   `json:"location"`
	Description  string   `json:"description"`
	Requirements []string `json:"requirements"`
	About        string   `json:"about"`
}

var mockJobs = []Job{
	{
		ID:           1,
		Title:        "Senior Go Developer",
		Company:      "TechCorp",
		Salary:       "200,000 - 300,000 ₽",
		Level:        "Senior",
		Type:         "Full-time",
		Location:     "Владивосток",
		Description:  "Ищем опытного Go разработчика для работы над микросервисами",
		Requirements: []string{"Go", "PostgreSQL", "Docker", "REST API"},
		About:        "TechCorp - IT компания с офисом во Владивостоке",
	},
	{
		ID:           2,
		Title:        "Backend Developer",
		Company:      "StartupXYZ",
		Salary:       "150,000 - 200,000 ₽",
		Level:        "Middle",
		Type:         "Full-time",
		Location:     "Москва (удаленно)",
		Description:  "Разработка API и работа с БД",
		Requirements: []string{"Go", "SQL", "Git"},
		About:        "Молодой стартап в сфере финтеха",
	},
}

var mockCandidates = []Candidate{
	{
		ID:       1,
		Name:     "Иван Иванов",
		Position: "Senior Go Developer",
		Location: "Владивосток",
		Email:    "ivan@example.com",
		Phone:    "+7 (999) 123-45-67",
		Avatar:   "👨‍💼",
		Rating:   4.8,
		About:    "Опытный разработчик с 8 годами стажа",
		Skills:   []string{"Go", "PostgreSQL", "Docker", "Kubernetes"},
		Experience: []Experience{
			{
				Title:   "Senior Developer",
				Company: "TechCorp",
				Period:  "2020 - настоящее",
			},
			{
				Title:   "Middle Developer",
				Company: "OtherCorp",
				Period:  "2017 - 2020",
			},
		},
		Education: []Education{
			{
				Degree: "Магистратура",
				School: "ДВФУ",
				Year:   "2015 - 2017",
			},
		},
	},
	{
		ID:       2,
		Name:     "Мария Петрова",
		Position: "Go Developer",
		Location: "Москва",
		Email:    "maria@example.com",
		Phone:    "+7 (999) 987-65-43",
		Avatar:   "👩‍💼",
		Rating:   4.5,
		About:    "Разработчик с опытом в микросервисах",
		Skills:   []string{"Go", "gRPC", "PostgreSQL", "Docker"},
		Experience: []Experience{
			{
				Title:   "Backend Developer",
				Company: "StartupXYZ",
				Period:  "2019 - настоящее",
			},
		},
		Education: []Education{
			{
				Degree: "Бакалавр",
				School: "МГУ",
				Year:   "2016 - 2020",
			},
		},
	},
}

func jsonResponse(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func logRequest(r *http.Request) {
	log.Printf("[%s] %s %s", r.Method, r.RequestURI, r.RemoteAddr)
}

func handleJobs(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	jsonResponse(w, http.StatusOK, mockJobs)
}

func handleCandidates(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	switch r.Method {
	case http.MethodGet:
		jsonResponse(w, http.StatusOK, mockCandidates)
	case http.MethodPost:
		var candidate Candidate
		if err := json.NewDecoder(r.Body).Decode(&candidate); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		candidate.ID = len(mockCandidates) + 1
		mockCandidates = append(mockCandidates, candidate)
		jsonResponse(w, http.StatusCreated, candidate)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/api/jobs", handleJobs)
	http.HandleFunc("/api/candidates", handleCandidates)
	http.HandleFunc("/api/", handleOptions) // Для preflight запросов

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
