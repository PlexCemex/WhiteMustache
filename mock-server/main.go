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
		Description:    "Сбор мусора на набережной, очистка прибрежной полосы от пластика и мусора",
		DateOfBegin:    parseDate("20260101"),
		DateOfEnd:      parseDate("20270101"),
		Salary:         1000,
		Title:          "Волонтер",
		DateOfDocument: time.Date(2026, 1, 5, 19, 41, 59, 0, time.UTC),
		TypesOfWork:    []string{"Общественная польза", "Экология"},
		Number:         "000000005",
	},
	{
		Organization:   "CODE WORK",
		Description:    "Необходимо обучать программированию студентов 1-2 курсов на языке С++, подготовка к олимпиадам",
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
		Description:    "Разработка backend API на Go, опыт обязателен, работа с PostgreSQL и Docker",
		DateOfBegin:    parseDate("20260215"),
		DateOfEnd:      parseDate("20260630"),
		Salary:         150000,
		Title:          "Go разработчик",
		DateOfDocument: time.Date(2026, 1, 7, 14, 20, 0, 0, time.UTC),
		TypesOfWork:    []string{"Программирование", "Go", "API", "Backend"},
		Number:         "000000006",
	},
	{
		Organization:   "DVFU Research Lab",
		Description:    "Помощь в проведении научных исследований в области искусственного интеллекта, обработка данных",
		DateOfBegin:    parseDate("20260201"),
		DateOfEnd:      parseDate("20261031"),
		Salary:         75000,
		Title:          "Научный ассистент",
		DateOfDocument: time.Date(2026, 1, 8, 11, 15, 0, 0, time.UTC),
		TypesOfWork:    []string{"Наука", "Программирование", "Алгоритмы"},
		Number:         "000000007",
	},
	{
		Organization:   "Hospital №1",
		Description:    "Администратор медицинского центра, работа с пациентами и документацией",
		DateOfBegin:    parseDate("20260101"),
		DateOfEnd:      parseDate("20261231"),
		Salary:         30000,
		Title:          "Администратор",
		DateOfDocument: time.Date(2026, 1, 9, 9, 0, 0, 0, time.UTC),
		TypesOfWork:    []string{"Медицина", "Помощь людям"},
		Number:         "000000008",
	},
	{
		Organization:   "Vladivostok Creative Studio",
		Description:    "Работа в команде креативных дизайнеров, создание графического контента для проектов",
		DateOfBegin:    parseDate("20260315"),
		DateOfEnd:      parseDate("20260915"),
		Salary:         45000,
		Title:          "Графический дизайнер",
		DateOfDocument: time.Date(2026, 1, 10, 13, 45, 0, 0, time.UTC),
		TypesOfWork:    []string{"Творчество", "Дизайн", "IT"},
		Number:         "000000009",
	},
	{
		Organization:   "Literature Center DVFU",
		Description:    "Редакция университетского издания, работа с текстами и публикациями",
		DateOfBegin:    parseDate("20260201"),
		DateOfEnd:      parseDate("20260930"),
		Salary:         25000,
		Title:          "Редактор",
		DateOfDocument: time.Date(2026, 1, 11, 10, 20, 0, 0, time.UTC),
		TypesOfWork:    []string{"Литература", "Редакция"},
		Number:         "000000010",
	},
	{
		Organization:   "Tech Startup",
		Description:    "Фронтенд разработчик для веб-приложения, опыт с React и TypeScript приветствуется",
		DateOfBegin:    parseDate("20260301"),
		DateOfEnd:      parseDate("20260831"),
		Salary:         120000,
		Title:          "Frontend Developer",
		DateOfDocument: time.Date(2026, 1, 7, 15, 30, 0, 0, time.UTC),
		TypesOfWork:    []string{"Программирование", "Frontend", "Web"},
		Number:         "000000011",
	},
	{
		Organization:   "ICPC Training Center",
		Description:    "Подготовка студентов к чемпионатам по программированию, тренировки по алгоритмам",
		DateOfBegin:    parseDate("20260101"),
		DateOfEnd:      parseDate("20261231"),
		Salary:         80000,
		Title:          "Тренер ICPC",
		DateOfDocument: time.Date(2026, 1, 9, 14, 0, 0, 0, time.UTC),
		TypesOfWork:    []string{"Программирование", "Алгоритмы", "ICPC", "Обучение"},
		Number:         "000000012",
	},
	{
		Organization:   "Data Science Lab",
		Description:    "Работа с большими данными, машинное обучение, анализ статистики",
		DateOfBegin:    parseDate("20260210"),
		DateOfEnd:      parseDate("20261210"),
		Salary:         160000,
		Title:          "Data Scientist",
		DateOfDocument: time.Date(2026, 1, 8, 16, 45, 0, 0, time.UTC),
		TypesOfWork:    []string{"Программирование", "Наука", "Технологии", "Алгоритмы"},
		Number:         "000000013",
	},
	{
		Organization:   "Green Initiative DVFU",
		Description:    "Экологический проект, уборка парков и посадка деревьев",
		DateOfBegin:    parseDate("20260320"),
		DateOfEnd:      parseDate("20261020"),
		Salary:         5000,
		Title:          "Волонтер Эколог",
		DateOfDocument: time.Date(2026, 1, 10, 11, 30, 0, 0, time.UTC),
		TypesOfWork:    []string{"Помощь пожилым", "Общественная польза", "Экология"},
		Number:         "000000014",
	},
	{
		Organization:   "Mobile Dev Studio",
		Description:    "Разработка мобильных приложений на Flutter и Kotlin, опыт в мобильной разработке",
		DateOfBegin:    parseDate("20260225"),
		DateOfEnd:      parseDate("20260825"),
		Salary:         135000,
		Title:          "Mobile Developer",
		DateOfDocument: time.Date(2026, 1, 9, 10, 15, 0, 0, time.UTC),
		TypesOfWork:    []string{"Программирование", "Mobile", "Технологии"},
		Number:         "000000015",
	},
	{
		Organization:   "Medical Research Institute",
		Description:    "Помощь в медицинских исследованиях, работа с пациентами и документацией",
		DateOfBegin:    parseDate("20260320"),
		DateOfEnd:      parseDate("20261120"),
		Salary:         65000,
		Title:          "Исследовательский ассистент",
		DateOfDocument: time.Date(2026, 1, 11, 14, 50, 0, 0, time.UTC),
		TypesOfWork:    []string{"Медицина", "Наука", "Помощь людям"},
		Number:         "000000016",
	},
	{
		Organization:   "Vladivostok Library",
		Description:    "Каталогизация книг, работа с библиотечной системой, помощь посетителям",
		DateOfBegin:    parseDate("20260201"),
		DateOfEnd:      parseDate("20261231"),
		Salary:         28000,
		Title:          "Библиотекарь",
		DateOfDocument: time.Date(2026, 1, 10, 9, 40, 0, 0, time.UTC),
		TypesOfWork:    []string{"Литература", "Культура"},
		Number:         "000000017",
	},
	{
		Organization:   "IoT Innovations",
		Description:    "Разработка на микроконтроллерах Arduino и Raspberry Pi, встроенные системы",
		DateOfBegin:    parseDate("20260401"),
		DateOfEnd:      parseDate("20261001"),
		Salary:         95000,
		Title:          "Embedded Systems Developer",
		DateOfDocument: time.Date(2026, 1, 8, 12, 20, 0, 0, time.UTC),
		TypesOfWork:    []string{"Программирование", "C++", "IoT", "Технологии"},
		Number:         "000000018",
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
		Number:       "000000001",
		Good:         true,
	},
	{
		Organization: "Tech Startup",
		Student:      "Петров Петр Петрович",
		Description:  "Разработчик с опытом 5 лет, знаю Go, PostgreSQL, Docker",
		StartPeriod:  "15.02.2026 0:00:00",
		EndPeriod:    "30.06.2026 0:00:00",
		Number:       "000000006",
		Good:         false,
	},
	{
		Organization: "DVFU Research Lab",
		Student:      "Сидоров Сергей Сергеевич",
		Description:  "Увлекаюсь AI и машинным обучением, имею опыт работы с Python и TensorFlow",
		StartPeriod:  "10.02.2026 0:00:00",
		EndPeriod:    "15.10.2026 0:00:00",
		Number:       "000000007",
		Good:         true,
	},
	{
		Organization: "Hospital №1",
		Student:      "Кузнецова Елена Викторовна",
		Description:  "Работала администратором в клинике 2 года, ответственная и пунктуальная",
		StartPeriod:  "01.01.2026 0:00:00",
		EndPeriod:    "31.12.2026 0:00:00",
		Number:       "000000008",
		Good:         true,
	},
	{
		Organization: "Vladivostok Creative Studio",
		Student:      "Морозова Анна Дмитриевна",
		Description:  "Дизайнер с портфолио, работала в 3 студиях, знаю Figma, Photoshop, Illustrator",
		StartPeriod:  "15.03.2026 0:00:00",
		EndPeriod:    "15.09.2026 0:00:00",
		Number:       "000000009",
		Good:         true,
	},
	{
		Organization: "Literature Center DVFU",
		Student:      "Соколов Максим Олегович",
		Description:  "Журналист, опыт работы в редакции университетской газеты 1 год",
		StartPeriod:  "01.02.2026 0:00:00",
		EndPeriod:    "30.09.2026 0:00:00",
		Number:       "000000010",
		Good:         false,
	},
	{
		Organization: "Tech Startup",
		Student:      "Волков Игорь Анатольевич",
		Description:  "Frontend разработчик, опыт 4 года, React, TypeScript, Vue.js",
		StartPeriod:  "01.03.2026 0:00:00",
		EndPeriod:    "31.08.2026 0:00:00",
		Number:       "000000011",
		Good:         false,
	},
	{
		Organization: "ICPC Training Center",
		Student:      "Лебедева Валентина Сергеевна",
		Description:  "Чемпионка регионального чемпионата по программированию, готова обучать",
		StartPeriod:  "01.01.2026 0:00:00",
		EndPeriod:    "31.12.2026 0:00:00",
		Number:       "000000012",
		Good:         true,
	},
	{
		Organization: "Data Science Lab",
		Student:      "Романов Константин Вячеславович",
		Description:  "Data Scientist с опытом 6 лет, Python, R, SQL, работал в крупных проектах",
		StartPeriod:  "10.02.2026 0:00:00",
		EndPeriod:    "10.12.2026 0:00:00",
		Number:       "000000013",
		Good:         true,
	},
	{
		Organization: "Green Initiative DVFU",
		Student:      "Никитина Ольга Ивановна",
		Description:  "Люблю природу, активно участвую в экологических акциях",
		StartPeriod:  "20.03.2026 0:00:00",
		EndPeriod:    "20.10.2026 0:00:00",
		Number:       "000000014",
		Good:         false,
	},
	{
		Organization: "Mobile Dev Studio",
		Student:      "Федоров Виталий Федорович",
		Description:  "Мобильный разработчик, опыт 3 года с Flutter и Kotlin, несколько приложений в AppStore",
		StartPeriod:  "25.02.2026 0:00:00",
		EndPeriod:    "25.08.2026 0:00:00",
		Number:       "000000015",
		Good:         true,
	},
	{
		Organization: "Medical Research Institute",
		Student:      "Смирнова Дарья Павловна",
		Description:  "Студентка медицинского факультета, хочу помогать в исследованиях",
		StartPeriod:  "20.03.2026 0:00:00",
		EndPeriod:    "20.11.2026 0:00:00",
		Number:       "000000016",
		Good:         false,
	},
	{
		Organization: "Vladivostok Library",
		Student:      "Голубев Артём Александрович",
		Description:  "Библиотекарь с опытом 4 года, знаю системы каталогизации, люблю работать с людьми",
		StartPeriod:  "01.02.2026 0:00:00",
		EndPeriod:    "31.12.2026 0:00:00",
		Number:       "000000017",
		Good:         true,
	},
	{
		Organization: "IoT Innovations",
		Student:      "Карпов Денис Игоревич",
		Description:  "Embedded системщик, опыт 5 лет, Arduino, Raspberry Pi, C++, Python",
		StartPeriod:  "01.04.2026 0:00:00",
		EndPeriod:    "01.10.2026 0:00:00",
		Number:       "000000018",
		Good:         false,
	},
}

var notifies = []Notify{
	{
		Text:            "Уважаемый Николаев Николай Николаевич! \n Одобрена ваша заявка по вакансии на должность Учитель по программированию на С++. \n Сообщение от руководителя: Подходите в кабинет C315 14.01.2026 с 13 до 14",
		Date:            time.Date(2026, 1, 11, 0, 0, 0, 0, time.UTC),
		NumberOfRequest: "000000001",
	},
	{
		Text:            "Уважаемый Иванов Иван Иванович! \n Спасибо за участие. К сожалению, мы выбрали другого кандидата. Удачи в поиске!",
		Date:            time.Date(2026, 1, 12, 0, 0, 0, 0, time.UTC),
		NumberOfRequest: "000000004",
	},
	{
		Text:            "Уважаемый Сидоров Сергей Сергеевич! \n Одобрена ваша заявка на позицию Научный ассистент. \n Встреча с руководителем: 15.02.2026 в 10:00 в офисе ДВФУ, кабинет 405",
		Date:            time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC),
		NumberOfRequest: "000000007",
	},
	{
		Text:            "Уважаемая Кузнецова Елена Викторовна! \n Вы приняты на должность Администратора. Начало работы: 01.02.2026. Явитесь в 09:00 с документами.",
		Date:            time.Date(2026, 1, 9, 0, 0, 0, 0, time.UTC),
		NumberOfRequest: "000000008",
	},
	{
		Text:            "Уважаемая Морозова Анна Дмитриевна! \n Одобрена ваша заявка на должность Графический дизайнер. \n Первое совещание команды: 16.03.2026 в 14:00",
		Date:            time.Date(2026, 1, 12, 0, 0, 0, 0, time.UTC),
		NumberOfRequest: "000000009",
	},
	{
		Text:            "Уважаемый Романов Константин Вячеславович! \n Поздравляем! Вы выбраны на должность Data Scientist. Контракт будет отправлен на почту.",
		Date:            time.Date(2026, 1, 11, 0, 0, 0, 0, time.UTC),
		NumberOfRequest: "000000013",
	},
	{
		Text:            "Уважаемый Федоров Виталий Федорович! \n Одобрена заявка на должность Mobile Developer. Собеседование в офисе: 01.03.2026 в 15:00",
		Date:            time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC),
		NumberOfRequest: "000000015",
	},
	{
		Text:            "Уважаемый Голубев Артём Александрович! \n Принято решение об одобрении вашей заявки. Начало работы: 01.03.2026",
		Date:            time.Date(2026, 1, 12, 0, 0, 0, 0, time.UTC),
		NumberOfRequest: "000000017",
	},
}

var tags = []string{
	"Наука",
	"Медицина",
	"Литература",
	"Технологии",
	"Творчество",
	"Программирование",
	"Алгоритмы",
	"ICPC",
	"Помощь пожилым",
	"Общественная польза",
	"Backend",
	"Frontend",
	"Go",
	"API",
	"Обучение",
	"C++",
	"Python",
	"Web",
	"Mobile",
	"IoT",
	"Дизайн",
	"Редакция",
	"Культура",
	"Экология",
}

var accountsDB = map[string]Account{
	"ivanov.ii": {
		Organization: "",
		Student:      "123-694-775 67",
	},
	"ivanov.iv": {
		Organization: "f2742040-cdb4-11f0-ae42-38d57ae2c1c1",
		Student:      "",
	},
	"petrov.pp": {
		Organization: "4c09ed30-cdb6-11f0-ae42-38d57ae2c1c1",
		Student:      "",
	},
	"smirnova.dp": {
		Organization: "",
		Student:      "234-567-890 12",
	},
	"kuznetsova.ev": {
		Organization: "7a8b9c0d-1e2f-3a4b-5c6d-7e8f9a0b1c2d",
		Student:      "",
	},
	"volkov.ia": {
		Organization: "",
		Student:      "345-678-901 23",
	},
	"lebedeva.vs": {
		Organization: "",
		Student:      "456-789-012 34",
	},
	"romanov.kv": {
		Organization: "e1f2a3b4-c5d6-7e8f-9a0b-1c2d3e4f5a6b",
		Student:      "",
	},
}

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
	if account, exists := accountsDB[user]; exists {
		response = account
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
		if v.Number == numberOfRequest {
			result = append(result, v)
			break
		}
	}

	if len(result) == 0 {
		result = vacancies[:1]
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
	mux.HandleFunc("/JobService/hs/jobservice/mynotify/", getNotifications)
	mux.HandleFunc("/JobService/hs/jobservice/vacancyfromnotify/", getVacancyFromNotify)
	mux.HandleFunc("/JobService/hs/jobservice/closevacancy/", closeVacancy)

	handler := corsMiddleware(mux)

	log.Fatal(http.ListenAndServe(":80", handler))
}
