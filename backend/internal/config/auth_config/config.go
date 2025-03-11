package authConfig

import (
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/repository/auth"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/repository/database"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/services/auth"
)

var (
	UserNameValidator auth.UserNameValidator
	PasswordValidator auth.PasswordValidator
	UserRepo         authRepository.UserRepository
	UserServiceLogin auth.UserServiceLogin
	UserServiceRegister auth.UserServiceRegister
)


func InitializeHandlers() {
	UserRepo = authRepository.NewUserRepository(database.DB)
	UserServiceLogin = auth.NewUsersServiceLogin(UserRepo, &UserNameValidator, &PasswordValidator)
	UserServiceRegister = auth.NewUsersServiceRegister(UserRepo, &UserNameValidator, &PasswordValidator)
}
