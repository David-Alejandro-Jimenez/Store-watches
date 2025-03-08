package authConfig

import (
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/repository/auth"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/repository/database"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/services/auth"
)

var (
	UserRepo         authRepository.UserRepository
	UserServiceLogin auth.UserServiceLogin
	UserServiceRegister auth.UserServiceRegister
)


func InitializeHandlers() {
	UserRepo = authRepository.NewUserRepository(database.DB)
	UserServiceLogin = auth.NewUsersServiceLogin(UserRepo)
	UserServiceRegister = auth.NewUsersService(UserRepo)
}
