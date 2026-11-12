package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"vetklinika/backend/internal/middleware"
	"vetklinika/backend/internal/storage"
)

// newTestHandler оборачивает обработчик тем же middleware, что и сервер,
// и подключает пустую БД во временной папке.
func newTestHandler(t *testing.T, h http.HandlerFunc) (http.Handler, *storage.Storage) {
	t.Helper()
	st, err := storage.New(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("storage.New() error = %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return middleware.InjectStorage(st)(h), st
}

func do(h http.Handler, method, target, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestPersons_CreateAndGet(t *testing.T) {
	h, _ := newTestHandler(t, HandlePersons)

	rec := do(h, http.MethodPost, "/api/persons",
		`{"name":"Иван","surname":"Петров","patronymic":"Сергеевич","phone":"9001234567","address":"Пермь"}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("POST status = %d, want %d; body: %s", rec.Code, http.StatusCreated, rec.Body)
	}
	var created struct {
		Id int `json:"id"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil || created.Id == 0 {
		t.Fatalf("POST body = %s, want {\"id\":N}", rec.Body)
	}

	rec = do(h, http.MethodGet, "/api/persons?id=1", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET by id status = %d, want 200", rec.Code)
	}
	var p storage.Person
	if err := json.Unmarshal(rec.Body.Bytes(), &p); err != nil {
		t.Fatalf("GET body is not a person: %s", rec.Body)
	}
	if p.Surname != "петров" {
		t.Errorf("surname = %q, want %q", p.Surname, "петров")
	}

	rec = do(h, http.MethodGet, "/api/persons", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET list status = %d, want 200", rec.Code)
	}
	var list []storage.Person
	if err := json.Unmarshal(rec.Body.Bytes(), &list); err != nil || len(list) != 1 {
		t.Errorf("GET list body = %s, want 1 person", rec.Body)
	}

	rec = do(h, http.MethodGet, "/api/persons?name=петров", "")
	if err := json.Unmarshal(rec.Body.Bytes(), &list); err != nil || len(list) != 1 {
		t.Errorf("search body = %s, want 1 person", rec.Body)
	}
}

func TestPersons_BadRequests(t *testing.T) {
	h, _ := newTestHandler(t, HandlePersons)

	tests := []struct {
		name   string
		method string
		target string
		body   string
		want   int
	}{
		{"невалидный JSON при создании", http.MethodPost, "/api/persons", `{"name":`, http.StatusBadRequest},
		{"пустой id", http.MethodGet, "/api/persons?id=", "", http.StatusBadRequest},
		{"id не число", http.MethodGet, "/api/persons?id=abc", "", http.StatusBadRequest},
		{"клиента нет", http.MethodGet, "/api/persons?id=999", "", http.StatusNotFound},
		{"неподдерживаемый метод", http.MethodPatch, "/api/persons", "", http.StatusMethodNotAllowed},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := do(h, tt.method, tt.target, tt.body)
			if rec.Code != tt.want {
				t.Errorf("status = %d, want %d; body: %s", rec.Code, tt.want, rec.Body)
			}
		})
	}
}

func TestHandlers_NoStorage(t *testing.T) {
	// Без middleware в контексте нет хранилища: обработчик должен
	// вернуть 500, а не упасть с паникой.
	tests := []struct {
		name   string
		h      http.HandlerFunc
		target string
	}{
		{"persons", HandlePersons, "/api/persons"},
		{"types", HandleTypes, "/api/types"},
		{"pets", HandlePets, "/api/pets"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := do(tt.h, http.MethodGet, tt.target, "")
			if rec.Code != http.StatusInternalServerError {
				t.Errorf("status = %d, want 500", rec.Code)
			}
		})
	}
}

func TestTypes_CreateAndList(t *testing.T) {
	h, _ := newTestHandler(t, HandleTypes)

	rec := do(h, http.MethodPost, "/api/types", `{"type":"кошка"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("POST status = %d, want 200; body: %s", rec.Code, rec.Body)
	}
	var created storage.Types
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil || created.Id == 0 || created.Type != "кошка" {
		t.Fatalf("POST body = %s", rec.Body)
	}

	rec = do(h, http.MethodPost, "/api/types", `{"type":"кошка"}`)
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("duplicate POST status = %d, want 500", rec.Code)
	}

	rec = do(h, http.MethodGet, "/api/types", "")
	var list []storage.Types
	if err := json.Unmarshal(rec.Body.Bytes(), &list); err != nil || len(list) != 1 {
		t.Errorf("GET list body = %s, want 1 type", rec.Body)
	}

	rec = do(h, http.MethodDelete, "/api/types", "")
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("DELETE status = %d, want 405", rec.Code)
	}
}

func TestPets_CreateValidation(t *testing.T) {
	h, _ := newTestHandler(t, HandlePets)

	tests := []struct {
		name   string
		target string
		body   string
		want   int
	}{
		{"нет ownerId", "/api/pets", `{"name":"Мурка"}`, http.StatusBadRequest},
		{"ownerId не число", "/api/pets?ownerId=abc", `{"name":"Мурка"}`, http.StatusBadRequest},
		{"невалидный JSON", "/api/pets?ownerId=1", `{"name":`, http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := do(h, http.MethodPost, tt.target, tt.body)
			if rec.Code != tt.want {
				t.Errorf("status = %d, want %d; body: %s", rec.Code, tt.want, rec.Body)
			}
		})
	}
}

func TestPets_Create(t *testing.T) {
	h, st := newTestHandler(t, HandlePets)

	ownerID, err := st.AddPerson(storage.Person{Name: "Иван"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateType(storage.Types{Type: "кошка"}); err != nil {
		t.Fatal(err)
	}

	body := `{"name":"Мурка","weight":4.2,"age":3,"age_write_date":"2026-01-10T00:00:00.000Z","type":1,"allergies":"нет"}`
	rec := do(h, http.MethodPost, "/api/pets?ownerId=1", body)
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body: %s", rec.Code, rec.Body)
	}

	pets, err := st.GetPetsByOwner(ownerID)
	if err != nil {
		t.Fatalf("GetPetsByOwner() error = %v", err)
	}
	if len(*pets) != 1 || (*pets)[0].Name != "Мурка" || (*pets)[0].Type != "кошка" {
		t.Errorf("pets of owner = %+v", *pets)
	}
}
