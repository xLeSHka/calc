package repositoryUsers

import (
	"github.com/xLeSHka/calc/internal/orchestrator/repository"
	"gorm.io/gorm"
)

type UsersRepository struct {
	DB *gorm.DB
}

func New(db *gorm.DB) repository.UsersRepository {
	return &UsersRepository{DB: db}
}
