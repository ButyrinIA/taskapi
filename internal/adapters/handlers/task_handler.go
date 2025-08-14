package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/ButyrinIA/taskapi/internal/core/models"
	"github.com/ButyrinIA/taskapi/internal/core/usecases"
)

type TaskHandler struct {
	taskUsecase usecases.TaskUsecase
}

func NewTaskHandler(taskUsecase usecases.TaskUsecase) *TaskHandler {
	return &TaskHandler{taskUsecase: taskUsecase}
}

func (h *TaskHandler) HandleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		status := r.URL.Query().Get("status")
		tasks := h.taskUsecase.GetTasks(status)
		if err := json.NewEncoder(w).Encode(tasks); err != nil {
			h.taskUsecase.LogError(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		var task models.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			h.taskUsecase.LogError(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		createdTask, err := h.taskUsecase.CreateTask(task)
		if err != nil {
			h.taskUsecase.LogError(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(createdTask); err != nil {
			h.taskUsecase.LogError(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *TaskHandler) HandleTaskByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(path)
	if err != nil {
		h.taskUsecase.LogError(err)
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, found := h.taskUsecase.GetTaskByID(id)
	if !found {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	if err := json.NewEncoder(w).Encode(task); err != nil {
		h.taskUsecase.LogError(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
