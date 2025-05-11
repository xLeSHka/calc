package repositoryAgent

import (
	"github.com/xLeSHka/calc/internal/orchestrator/repository"
	"gorm.io/gorm"
)

type AgentRepository struct {
	DB *gorm.DB
}

func New(db *gorm.DB) repository.AgentRepository {
	return &AgentRepository{DB: db}
}
