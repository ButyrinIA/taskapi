package repositories

import (
	"testing"

	"github.com/ButyrinIA/taskapi/internal/core/models"
)

func TestInMemoryTaskRepository(t *testing.T) {
	repo := NewInMemoryTaskRepository()

	//Create
	task := models.Task{Title: "Test Task", Description: "Description", Status: "todo"}
	created := repo.Create(task)
	if created.ID != 1 || created.Title != "Test Task" {
		t.Errorf("Expected ID 1 and Title 'Test Task', got %d %s", created.ID, created.Title)
	}

	//GetByID
	got, found := repo.GetByID(1)
	if !found || got.Title != "Test Task" {
		t.Errorf("Expected found and Title 'Test Task', got %v %s", found, got.Title)
	}
	_, found = repo.GetByID(999)
	if found {
		t.Error("Expected not found for ID 999")
	}

	//GetAll
	all := repo.GetAll("")
	if len(all) != 1 {
		t.Errorf("Expected 1 task, got %d", len(all))
	}

	//GetAll with filter
	filtered := repo.GetAll("todo")
	if len(filtered) != 1 {
		t.Errorf("Expected 1 task for status 'todo', got %d", len(filtered))
	}
	filtered = repo.GetAll("done")
	if len(filtered) != 0 {
		t.Errorf("Expected 0 tasks for status 'done', got %d", len(filtered))
	}
}
