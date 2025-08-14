package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ButyrinIA/taskapi/internal/core/models"
)

type mockUsecase struct {
	tasks []models.Task
}

func (m *mockUsecase) GetTasks(status string) []models.Task {
	var res []models.Task
	for _, t := range m.tasks {
		if status == "" || t.Status == status {
			res = append(res, t)
		}
	}
	return res
}

func (m *mockUsecase) GetTaskByID(id int) (models.Task, bool) {
	if id == 1 && len(m.tasks) > 0 {
		return m.tasks[0], true
	}
	return models.Task{}, false
}

func (m *mockUsecase) CreateTask(task models.Task) (models.Task, error) {
	if task.Title == "" {
		return models.Task{}, fmt.Errorf("title is required")
	}
	task.ID = len(m.tasks) + 1
	m.tasks = append(m.tasks, task)
	return task, nil
}

func (m *mockUsecase) LogError(err error) {
}

func TestTaskHandler(t *testing.T) {
	uc := &mockUsecase{}
	handler := NewTaskHandler(uc)

	//GET /tasks
	req := httptest.NewRequest(http.MethodGet, "/tasks?status=todo", nil)
	w := httptest.NewRecorder()
	handler.HandleTasks(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	//POST /tasks (success)
	task := models.Task{Title: "Test Task", Description: "Desc", Status: "todo"}
	body, _ := json.Marshal(task)
	req = httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	w = httptest.NewRecorder()
	handler.HandleTasks(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	//POST /tasks (error: empty title)
	task.Title = ""
	body, _ = json.Marshal(task)
	req = httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	w = httptest.NewRecorder()
	handler.HandleTasks(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	//GET /tasks/1
	req = httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
	w = httptest.NewRecorder()
	handler.HandleTaskByID(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	//GET /tasks/999 (not found)
	req = httptest.NewRequest(http.MethodGet, "/tasks/999", nil)
	w = httptest.NewRecorder()
	handler.HandleTaskByID(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}

	//GET /tasks/invalid (invalid ID)
	req = httptest.NewRequest(http.MethodGet, "/tasks/invalid", nil)
	w = httptest.NewRecorder()
	handler.HandleTaskByID(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}
