package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"tasks-api/internal/models"
	"tasks-api/internal/storage"
)

type Handler struct{ Store storage.Storage }

type ErrorMsg struct {
	Msg string `json:"error"`
}
type StatusMsg struct {
	Msg string `json:"status"`
}

func New(s storage.Storage) *Handler { return &Handler{Store: s} }

// /tasks (GET, POST)
func (h *Handler) TasksCollection(w http.ResponseWriter, r *http.Request) {
	// TODO: реализуйте разбор метода, JSON, коды статусов, валидацию
	w.Header().Set("Content-Type", "application/json")

	log.Printf("%s : %s\n", r.Method, r.RequestURI)
	switch r.Method {
	case http.MethodGet:
		{
			//реализация GET запроса
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(h.Store.List())
		}
	case http.MethodPost:
		{
			//реализация POST запроса
			var req models.Task
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorMsg{Msg: "Некорректный JSON"})
				return
			}
			_, err := h.Store.Create(req)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(ErrorMsg{Msg: "Задача не создана"})
				return
			}
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(StatusMsg{Msg: "Задача успешно создана"})
		}
	default:
		{
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(ErrorMsg{Msg: "Метод не поддерживается"})
		}
	}
}

// /tasks/{id} (GET, PUT, DELETE)
func (h *Handler) TaskItem(w http.ResponseWriter, r *http.Request) {
	// TODO: извлечение id, маршрутизация по методу, ошибки
	log.Printf("%s : %s\n", r.Method, r.RequestURI)
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		{ //реализация GET запроса
			id, err := strconv.Atoi(r.URL.Query().Get("id"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorMsg{Msg: "id не целое число"})
				return
			}
			task, exist := h.Store.Get(id)
			if !exist {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(ErrorMsg{Msg: "Не существующий id"})
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(task)
		}
	case http.MethodDelete:
		{
			id, err := strconv.Atoi(r.URL.Query().Get("id"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorMsg{Msg: "id не целое число"})
				return
			}
			err = h.Store.Delete(id)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(ErrorMsg{Msg: err.Error()})
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(StatusMsg{Msg: "Задача удалена"})
		}
	case http.MethodPut:
		{
			id, err := strconv.Atoi(r.URL.Query().Get("id"))
			var req models.Task
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorMsg{Msg: "id не целое число"})
				return
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorMsg{Msg: "Некорректный JSON"})
				return
			}
			req, err = h.Store.Update(id, req)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(ErrorMsg{Msg: err.Error()})
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(StatusMsg{Msg: "Задача обновлена"})
		}
	default:
		{
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(ErrorMsg{Msg: "Метод не поддерживается"})
		}
	}
}
