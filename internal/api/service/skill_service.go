package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/FPRPL26/rpl-be/internal/api/repository"
	"github.com/FPRPL26/rpl-be/internal/dto"
	myerror "github.com/FPRPL26/rpl-be/internal/pkg/error"
	"github.com/FPRPL26/rpl-be/internal/pkg/meta"
	"gorm.io/gorm"
)

type (
	SkillService interface {
		GetAll(ctx context.Context, metaReq meta.Meta) ([]dto.SkillResponse, meta.Meta, error)
		GetById(ctx context.Context, id int64) (dto.SkillResponse, error)
	}

	skillService struct {
		skillRepo repository.SkillRepository
	}
)

func NewSkillService(skillRepo repository.SkillRepository) SkillService {
	return &skillService{skillRepo}
}

func (s *skillService) GetAll(ctx context.Context, metaReq meta.Meta) ([]dto.SkillResponse, meta.Meta, error) {
	skills, metaRes, err := s.skillRepo.GetAll(ctx, nil, metaReq)
	if err != nil {
		return nil, metaRes, err
	}

	res := make([]dto.SkillResponse, 0, len(skills))
	for _, sk := range skills {
		res = append(res, dto.SkillResponse{
			ID:   sk.ID,
			Name: sk.Name,
		})
	}

	return res, metaRes, nil
}

func (s *skillService) GetById(ctx context.Context, id int64) (dto.SkillResponse, error) {
	sk, err := s.skillRepo.GetById(ctx, nil, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.SkillResponse{}, myerror.New("skill not found", http.StatusNotFound)
		}
		return dto.SkillResponse{}, err
	}

	return dto.SkillResponse{
		ID:   sk.ID,
		Name: sk.Name,
	}, nil
}
