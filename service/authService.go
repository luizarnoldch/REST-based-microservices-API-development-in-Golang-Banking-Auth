package service

import(
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/luizarnoldch/REST-based-microservices-API-development-in-Golang/tree/master/banking-auth/domain"
	"github.com/luizarnoldch/REST-based-microservices-API-development-in-Golang/tree/master/banking-auth/dto"
	"github.com/luizarnoldch/REST-based-microservices-API-development-in-Golang/tree/master/banking-lib/logger"
	"github.com/luizarnoldch/REST-based-microservices-API-development-in-Golang/tree/master/banking-lib/errs"
)

type AuthService interface {
	Login()
}

type DefaultAuthService struct {
	repo domain.authRepository
	rolePermissions domain.RolePermissions
}