package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task структура описывающая задачу
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// getTasks - хендлер возвращает всю мапу tasks в формате JSON
func getTasks(w http.ResponseWriter, r *http.Request) {
	// сериализуем данные из мапы tasks
	resp, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		// статус 500 Internal Server Error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// в заголовок записываем тип контента JSON
	w.Header().Set("Content-Type", "application/json")
	// статус 200 OK
	w.WriteHeader(http.StatusOK)
	// записываем сериализованные в JSON данные в тело ответа
	w.Write(resp)
}

// postTask - хендлер принимает задачу в теле запроса и сохраняет ее в мапе tasks
func postTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		// статус 400 Bad Request
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		// статус 400 Bad Request
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tasks[task.ID] = task

	w.Header().Set("Content-Type", "application/json")
	// статус 201 Created
	w.WriteHeader(http.StatusCreated)
}

// getTask - хендлер возвращает задачу с указанным в запросе пути ID, если такая есть в мапе tasks
func getTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	task, ok := tasks[id]
	if !ok {
		// статус 400 Bad Request
		http.Error(w, "Задача "+id+" не найдена", http.StatusBadRequest)
		return
	}

	resp, err := json.MarshalIndent(task, "", "  ")
	if err != nil {
		// статус 400 Bad Request
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// статус 200 OK
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// delTask - хендлер удаляет задачу из мапы tasks по её ID
func delTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	task, ok := tasks[id]
	if !ok {
		// статус 400 Bad Request
		http.Error(w, "Задача "+id+" не найдена", http.StatusBadRequest)
		return
	}

	delete(tasks, task.ID)

	w.Header().Set("Content-Type", "application/json")
	// статус 200 OK
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	r.Get("/tasks", getTasks)
	r.Post("/tasks", postTask)
	r.Get("/tasks/{id}", getTask)
	r.Delete("/tasks/{id}", delTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
