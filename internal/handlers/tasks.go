package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"tasks-api/internal/models"
	"tasks-api/internal/storage"
)

type LogLevel int

const (
	Debug LogLevel = iota
	Info
)

// logLevel задает текущий минимальный логируемый уровень
var logLevel = Info

func Log(level LogLevel, msg string, args ...interface{}) {
	if level < logLevel {
		return // Пропускаем логи ниже порога
	}
	prefix := ""
	switch level {
	case Debug:
		prefix = "DEBUG"
	case Info:
		prefix = "INFO"
	}

	log.Printf("[%s] %s", prefix, fmt.Sprintf(msg, args...))

}

type Handler struct{ Store storage.Storage }

type ErrorMsg struct {
	Msg string `json:"error"`
}
type StatusMsg struct {
	Msg string `json:"status"`
}

type AllowMsg struct {
	Msg string `json:"Allow"`
}

var allow = AllowMsg{Msg: fmt.Sprintf("%s, %s, %s, %s", http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete)}

func New(s storage.Storage) *Handler { return &Handler{Store: s} }

// /tasks (GET, POST)
func (h *Handler) TasksCollection(w http.ResponseWriter, r *http.Request) {
	// TODO: реализуйте разбор метода, JSON, коды статусов, валидацию
	w.Header().Set("Content-Type", "application/json")

	Log(Info, "%s : %s\n", r.Method, r.RequestURI)

	switch r.Method {
	case http.MethodGet:
		{
			//реализация GET запроса
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(h.Store.List())
		}
	case http.MethodPost:
		{
			body, error := io.ReadAll(r.Body)
			if error != nil {
				Log(Debug, error.Error())
			} else {
				Log(Debug, string(body))
			}
			//реализация POST запроса
			var req models.Task
			//декодируем JSON из запроса
			//если не удалось разобрать JSON
			if err := json.Unmarshal(body, &req); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				Log(Debug, err.Error())
				json.NewEncoder(w).Encode(ErrorMsg{Msg: "Incorrect JSON"})
				return
			}
			//если имя задачи не задано
			if req.Title == "" {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorMsg{Msg: "Task title is empty"})
				return
			}
			//если возникла ошибка при создании
			newTask, err := h.Store.Create(req)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(ErrorMsg{Msg: "Task not created"})
				return
			}
			//если все прошло успешно
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(newTask)
		}
	default:
		{
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(allow)
		}
	}
}

// /tasks/{id} (GET, PUT, DELETE)
func (h *Handler) TaskItem(w http.ResponseWriter, r *http.Request) {
	// TODO: извлечение id, маршрутизация по методу, ошибки
	Log(Info, "%s : %s\n", r.Method, r.RequestURI)
	w.Header().Set("Content-Type", "application/json")
	res := strings.Split(r.URL.Path, "/") //разбиваем URL на элементы
	id_str := res[len(res)-1]             //и берем последний в виде строки
	switch r.Method {
	case http.MethodGet:
		{ //реализация GET запроса
			id, err := strconv.ParseInt(id_str, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorMsg{Msg: "The input value given is not a valid integer"})
				return
			}
			task, exist := h.Store.Get(id)
			if !exist {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(ErrorMsg{Msg: "Wrong id"})
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(task)
		}
	case http.MethodDelete:
		{
			id, err := strconv.ParseInt(id_str, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorMsg{Msg: "The input value given is not a valid integer"})
				return
			}
			err = h.Store.Delete(id)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(ErrorMsg{Msg: err.Error()})
				return
			}
			w.WriteHeader(http.StatusNoContent)
			// json.NewEncoder(w).Encode(StatusMsg{Msg: "Задача удалена"})
		}
	case http.MethodPut:
		{
			body, error := io.ReadAll(r.Body)
			if error != nil {
				Log(Debug, error.Error())
			} else {
				Log(Debug, string(body))
			}
			//если вместо id передаличто-то неподходящее (буквы например)
			id, err := strconv.ParseInt(id_str, 10, 64)
			var req models.Task
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorMsg{Msg: "The input value given is not a valid integer"})
				return
			}
			//если пришел неверный Content-Type
			contentType := r.Header.Get("Content-Type")
			if contentType != "application/json" {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorMsg{Msg: "Wrong Content-Type"})
				return
			}
			//если не удалось разобрать JSON
			if err := json.Unmarshal(body, &req); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorMsg{Msg: "Incorrect JSON"})
				return
			}
			//если имя задачи не задано
			if req.Title == "" {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorMsg{Msg: "Task title is empty"})
				return
			}
			//если все прошло успешно
			req, err = h.Store.Update(id, req)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(ErrorMsg{Msg: err.Error()})
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(StatusMsg{Msg: "Task succesfully updated"})
		}
	default:
		{
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(allow)
		}
	}
}
