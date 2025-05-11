package repositoryUsers

import "github.com/xLeSHka/calc/internal/models"

func (r *UsersRepository) Register(person *models.User) (*models.User, error) {
	err := r.DB.Create(person).Error
	if err != nil {
		return nil, err
	}
	return person, nil
}
