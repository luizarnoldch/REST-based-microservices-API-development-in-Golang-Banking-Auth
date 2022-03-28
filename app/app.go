package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/luizarnoldch/REST-based-microservices-API-development-in-Golang/tree/master/banking-auth/domain"
	"github.com/luizarnoldch/REST-based-microservices-API-development-in-Golang/tree/master/banking-auth/service"
	"github.com/luizarnoldch/REST-based-microservices-API-development-in-Golang/tree/master/banking-lib/logger"
)

func Start() {
	sanityCheck()
	router := mux.NewRouter()
	authRepository := domain.NewAuthRepository(getDbClient())
	ah := AuthHandler{service.NewLoginService(authRepository, domain.GetRolePermissions(authRepository))}

	router.HandleFunc("/auth/login", ah.Login).Methods(http.MethodPost)
	router.HandleFunc("/auth/register", ah.NotImplementedHandler).Methods(http.MethodPost)
	router.HandleFunc("/auth/verify", ah.Verify).Methods(http.MethodGet)

	address := os.Getenv("SERVER_ADDRESS")
	port := os.Getenv("SERVER_PORT")
	log.Println(fmt.Sprintf("Starting 0Auth server on #{address}:#{port} ..."))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("#{address}:#{port}"), router))
}

func getDbClient() *sqlx.DB {
	driver := "mysql"
	//usuario := "root"
	usuario := os.Getenv("DB_USER")
	//pass := "u1OboD93110614"
	pass := os.Getenv("DB_PASSWD")
	//port := "tcp(localhost:3306)"
	port := os.Getenv("DB_ADDRESS_PORT")
	//table := "banking"
	table := os.Getenv("DB_NAME")

	dataSource := fmt.Sprintf("%s:%s@%s/%s", usuario, pass, port, table)

	client, err := sqlx.Open(driver, dataSource)
	if err != nil {
		panic(err)
	}

	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)
	return client
}

func sanityCheck() {
	envProps := []string{
		"SERVER_ADDRESS",
		"SERVER_PORT",
		"DB_USER",
		"DB_PASSWD",
		"DB_ADDR",
		"DB_PORT",
		"DB_NAME",
	}
	for _, k := range envProps {
		if os.Getenv(k) == "" {
			logger.Error(fmt.Sprintf("Environment variable %s not defined. Terminating application...", k))
		}
	}
}
