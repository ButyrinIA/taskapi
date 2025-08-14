package repositories

import (
	"sync"

	"github.com/ButyrinIA/taskapi/internal/core/models"
)

type InMemoryTaskRepository struct {
	mu     sync.Mutex
	tasks  map[int]models.Task
	nextID int
}

func NewInMemoryTaskRepository() *InMemoryTaskRepository {
	return &InMemoryTaskRepository{
		tasks:  make(map[int]models.Task),
		nextID: 1,
	}
}

func (r *InMemoryTaskRepository) GetAll(status string) []models.Task {
	r.mu.Lock()
	defer r.mu.Unlock()

	var result []models.Task
	for _, task := range r.tasks {
		if status == "" || task.Status == status {
			result = append(result, task)
		}
	}
	return result
}

func (r *InMemoryTaskRepository) GetByID(id int) (models.Task, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	task, found := r.tasks[id]
	return task, found
}

func (r *InMemoryTaskRepository) Create(task models.Task) models.Task {
	r.mu.Lock()
	defer r.mu.Unlock()

	task.ID = r.nextID
	r.nextID++
	r.tasks[task.ID] = task
	return task
}
