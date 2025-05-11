package repositoryUsers

import "github.com/xLeSHka/calc/internal/models"

func (r *UsersRepository) Login(login string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("login = ?", login).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
