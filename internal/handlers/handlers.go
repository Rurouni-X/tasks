package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"tasks/internal/database"
	"tasks/internal/models"
)

type Handlers struct {
	store *database.TaskStore
}

func NewHandlers(store *database.TaskStore) *Handlers {
	return &Handlers{store: store}
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, map[string]string{"error": message})
}

func getID(w http.ResponseWriter, r *http.Request) (int, bool) {
	path := strings.Split(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/")
	id, err := strconv.Atoi(path[0])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "некорректный ID задачи")
		return 0, false
	}
	return id, true
}

func (h *Handlers) GetAllTasks(w http.ResponseWriter, r *http.Request) {

	task, err := h.store.GetAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка получения задач")
		return
	}
	respondWithJSON(w, http.StatusOK, task)
}

func (h *Handlers) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	id, ok := getID(w, r)

	if !ok {
		return
	}

	task, err := h.store.GetById(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, task)
}

func (h *Handlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.CreateTaskInput

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		respondWithError(w, http.StatusBadRequest, "некорректный запрос")
		return
	}

	if strings.TrimSpace(task.Title) == "" {
		respondWithError(w, http.StatusBadRequest, "Пустой заголовок")
		return
	}

	result, err := h.store.Create(task)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, result)
}

func (h *Handlers) UpdateTask(w http.ResponseWriter, r *http.Request) {

	id, ok := getID(w, r)

	if !ok {
		return
	}

	var input models.UpdateTaskInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "некорректный запрос")
		return
	}

	if input.Title != nil && strings.TrimSpace(*input.Title) == "" {
		respondWithError(w, http.StatusBadRequest, "Заголовок обязателен")
		return
	}

	task, err := h.store.Update(id, input)

	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, task)
}

func (h *Handlers) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, ok := getID(w, r)

	if !ok {
		return
	}

	err := h.store.Delete(id)

	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusNoContent, map[string]string{"message": "success"})
}
