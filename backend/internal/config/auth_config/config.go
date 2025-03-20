package authConfig

import (
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/repository/auth"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/repository/database"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/services/auth"
	"github.com/David-Alejandro-Jimenez/sale-watches/pkg/security/security_auth"
)

var (
	UserNameValidator auth.UserNameValidator
	PasswordValidator auth.PasswordValidator
	UserRepo         authRepository.UserRepository
	UserServiceLogin auth.UserServiceLogin
	UserServiceRegister auth.UserServiceRegister
	SaltGenerator securityAuth.Generator
	Hasher securityAuth.Hasher
)


func InitializeHandlers() {
	UserNameValidator = auth.UserNameValidator{}
	PasswordValidator = auth.PasswordValidator{}
	SaltGenerator = securityAuth.RandomSaltGenerator{}
	Hasher = securityAuth.BcryptHasher{}
	UserRepo = authRepository.NewUserRepository(database.DB, SaltGenerator, Hasher)
	UserServiceLogin = auth.NewUsersServiceLogin(UserRepo, &UserNameValidator, &PasswordValidator)
	UserServiceRegister = auth.NewUsersServiceRegister(UserRepo, &UserNameValidator, &PasswordValidator)
}
