package usecases

import (
	"strings"
	"testing"

	"github.com/ButyrinIA/taskapi/internal/core/models"
)

type mockRepo struct {
	tasks map[int]models.Task
}

func (m *mockRepo) GetAll(status string) []models.Task {
	var res []models.Task
	for _, t := range m.tasks {
		if status == "" || t.Status == status {
			res = append(res, t)
		}
	}
	return res
}

func (m *mockRepo) GetByID(id int) (models.Task, bool) {
	t, ok := m.tasks[id]
	return t, ok
}

func (m *mockRepo) Create(task models.Task) models.Task {
	id := len(m.tasks) + 1
	task.ID = id
	m.tasks[id] = task
	return task
}

func TestTaskUsecase(t *testing.T) {
	logChan := make(chan string, 10)
	repo := &mockRepo{tasks: make(map[int]models.Task)}
	uc := NewTaskUsecase(repo, logChan)

	//CreateTask
	task := models.Task{Title: "Test Task", Description: "Desc", Status: "todo"}
	created, err := uc.CreateTask(task)
	if err != nil || created.ID != 1 || created.Title != "Test Task" {
		t.Errorf("Expected no error and ID 1, Title 'Test Task', got %v %d %s", err, created.ID, created.Title)
	}

	//CreateTask with error
	task.Title = ""
	_, err = uc.CreateTask(task)
	if err == nil || err.Error() != "title is required" {
		t.Errorf("Expected error 'title is required', got %v", err)
	}

	//GetTasks
	tasks := uc.GetTasks("todo")
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}
	tasks = uc.GetTasks("done")
	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks, got %d", len(tasks))
	}

	//GetTaskByID
	got, found := uc.GetTaskByID(1)
	if !found || got.Title != "Test Task" {
		t.Errorf("Expected found and Title 'Test Task', got %v %s", found, got.Title)
	}
	_, found = uc.GetTaskByID(999)
	if found {
		t.Error("Expected not found for ID 999")
	}

	// Check logs
	select {
	case msg := <-logChan:
		if !strings.Contains(msg, "Action: CreateTask") && !strings.Contains(msg, "Error: title is required") {
			t.Errorf("Expected log message to contain action or error, got %s", msg)
		}
	default:
	}
}
