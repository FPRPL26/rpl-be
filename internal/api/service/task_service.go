package service

import (
	"context"

	"github.com/FPRPL26/rpl-be/internal/api/repository"
	"github.com/FPRPL26/rpl-be/internal/dto"
	"github.com/FPRPL26/rpl-be/internal/entity"
	"github.com/FPRPL26/rpl-be/internal/pkg/meta"
)

type (
	TaskService interface {
		Create(ctx context.Context, req dto.CreateTaskRequest) (entity.Task, error)
		GetAll(ctx context.Context, metaReq meta.Meta) ([]entity.Task, meta.Meta, error)
		GetById(ctx context.Context, taskId string) (entity.Task, error)
		Update(ctx context.Context, taskId string, req dto.UpdateTaskRequest) (entity.Task, error)
		Delete(ctx context.Context, taskId string) error
	}

	taskService struct {
		taskRepo  repository.TaskRepository
		mediaRepo repository.MediaAssetRepository
	}
)

func NewTask(taskRepo repository.TaskRepository, mediaRepo repository.MediaAssetRepository) TaskService {
	return &taskService{taskRepo, mediaRepo}
}

func (s *taskService) Create(ctx context.Context, req dto.CreateTaskRequest) (entity.Task, error) {
	taskCreateResult, err := s.taskRepo.Create(ctx, nil, entity.Task{
		PhotoUrl:    req.PhotoUrl,
		Description: req.Description,
		Deadline:    req.Deadline,
		Status:      entity.TaskStatus(req.Status),
	})
	if err != nil {
		return entity.Task{}, err
	}

	if req.PhotoUrl != nil {
		if err := s.mediaRepo.MarkAsUsed(*req.PhotoUrl); err != nil {
			return taskCreateResult, err
		}
	}

	return taskCreateResult, nil
}

func (s *taskService) GetAll(ctx context.Context, metaReq meta.Meta) ([]entity.Task, meta.Meta, error) {
	return s.taskRepo.GetAll(ctx, nil, metaReq)
}

func (s *taskService) GetById(ctx context.Context, taskId string) (entity.Task, error) {
	return s.taskRepo.GetById(ctx, nil, taskId)
}

func (s *taskService) Update(ctx context.Context, taskId string, req dto.UpdateTaskRequest) (entity.Task, error) {
	task, err := s.taskRepo.GetById(ctx, nil, taskId)
	if err != nil {
		return entity.Task{}, err
	}

	var oldPhotoURL string
	if task.PhotoUrl != nil {
		oldPhotoURL = *task.PhotoUrl
	}

	task.PhotoUrl = req.PhotoUrl
	task.Description = req.Description
	task.Deadline = req.Deadline
	task.Status = entity.TaskStatus(req.Status)

	updateTaskResult, err := s.taskRepo.Update(ctx, nil, task)
	if err != nil {
		return entity.Task{}, err
	}

	if req.PhotoUrl != nil && *req.PhotoUrl != oldPhotoURL {
		s.mediaRepo.MarkAsUnused(oldPhotoURL)
		s.mediaRepo.MarkAsUsed(*req.PhotoUrl)
	} else if req.PhotoUrl == nil && oldPhotoURL != "" {
		s.mediaRepo.MarkAsUnused(oldPhotoURL)
	}

	return updateTaskResult, nil
}

func (s *taskService) Delete(ctx context.Context, taskId string) error {
	task, err := s.taskRepo.GetById(ctx, nil, taskId)
	if err != nil {
		return err
	}

	if err := s.taskRepo.Delete(ctx, nil, task); err != nil {
		return err
	}

	if task.PhotoUrl != nil {
		s.mediaRepo.MarkAsUnused(*task.PhotoUrl)
	}
	return nil
}
