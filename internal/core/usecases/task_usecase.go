package usecases

import (
	"fmt"

	"github.com/ButyrinIA/taskapi/internal/core/models"
)

type TaskRepository interface {
	GetAll(status string) []models.Task
	GetByID(id int) (models.Task, bool)
	Create(task models.Task) models.Task
}

type TaskUsecase interface {
	GetTasks(status string) []models.Task
	GetTaskByID(id int) (models.Task, bool)
	CreateTask(task models.Task) (models.Task, error)
	LogError(err error)
}

type taskUsecase struct {
	repo    TaskRepository
	logChan chan<- string
}

func NewTaskUsecase(repo TaskRepository, logChan chan<- string) TaskUsecase {
	return &taskUsecase{repo: repo, logChan: logChan}
}

func (u *taskUsecase) GetTasks(status string) []models.Task {
	tasks := u.repo.GetAll(status)
	u.logAsync(fmt.Sprintf("Action: GetTasks, Filter: %s, Count: %d", status, len(tasks)))
	return tasks
}

func (u *taskUsecase) GetTaskByID(id int) (models.Task, bool) {
	task, found := u.repo.GetByID(id)
	if found {
		u.logAsync(fmt.Sprintf("Action: GetTaskByID, ID: %d", id))
	} else {
		u.LogError(fmt.Errorf("task with ID %d not found", id))
	}
	return task, found
}

func (u *taskUsecase) CreateTask(task models.Task) (models.Task, error) {
	if task.Title == "" {
		err := fmt.Errorf("title is required")
		u.LogError(err)
		return models.Task{}, err
	}
	created := u.repo.Create(task)
	u.logAsync(fmt.Sprintf("Action: CreateTask, ID: %d, Title: %s", created.ID, created.Title))
	return created, nil
}

func (u *taskUsecase) logAsync(msg string) {
	go func() {
		u.logChan <- msg
	}()
}

func (u *taskUsecase) LogError(err error) {
	go func() {
		u.logChan <- fmt.Sprintf("Error: %v", err)
	}()
}
