package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"tasks-api/internal/handlers"
	"tasks-api/internal/memory"

	// "tasks-api/internal/storage"
	"tasks-api/internal/models"
)

func main() {
	// TODO: подключите конкретную реализацию (in‑memory) интерфейса Storage
	// var store storage.Storage // = memory.New() // реализуйте сами
	var port = flag.Int("p", 8080, "укажите номер порта (0-65535), по умолчанию используется порт 8080")
	flag.Parse()
	if *port < 0 || *port > 65535 {
		*port = 8080
	}
	store := memory.New()
	//добавляем задачи для тестовых целей
	store.Create(models.Task{Title: "Дефрагментация диска", Done: false})
	store.Create(models.Task{Title: "Сборка мусора", Done: false})
	store.Create(models.Task{Title: "Оптимизация запросов", Done: false})

	h := handlers.New(store)

	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", h.TasksCollection) // GET, POST
	mux.HandleFunc("/tasks/", h.TaskItem)       // GET, PUT, DELETE

	log.Printf("server listening on :%d\n", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), mux); err != nil {
		log.Fatal(err)
	}
}
