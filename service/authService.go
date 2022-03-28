package service

import(
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/luizarnoldch/REST-based-microservices-API-development-in-Golang-Banking-Auth/domain"
	"github.com/luizarnoldch/REST-based-microservices-API-development-in-Golang-Banking-Auth/dto"
	"github.com/luizarnoldch/REST-based-microservices-API-development-in-Golang-Banking-Lib/logger"
	"github.com/luizarnoldch/REST-based-microservices-API-development-in-Golang-Banking-Lib/errs"
)

type AuthService interface {
	Login()
}

type DefaultAuthService struct {
	repo domain.authRepository
	rolePermissions domain.RolePermissions
}