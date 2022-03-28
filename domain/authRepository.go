package domain

import(
	"database/sql"
	"github.com/luizarnoldch/REST-based-microservices-API-development-in-Golang-Banking-Lib/logger"
	"github.com/luizarnoldch/REST-based-microservices-API-development-in-Golang-Banking-Lib/errs"
	"github.com/jmoiron/sqlx"
)

type AuthRepository struct {
	FinBy(username string,password string) (*Login, *errs.AppError)
	GenerateAndSaveRefreshTokenStore(authToken AuthToken) (string,*errs.AppError)
	RefreshTokenExists(refreshToken string) *errs.AppError
}

type AuthRepositoryDb struct {
	client *sqlx.DB
}

func (d AuthRepositoryDb) RefreshTokenExists(refreshToken string) *errs.AppError{
	sqlSelect := "SELECt refresh_token FROM refresh_token_store WHERE refresh_token = ?"
	var token string
	err := d.client.Get(&token,sqlSelect,refreshToken)
	if err != nil{
		if err == sql.ErrNoRows {
			return errs.NewAuthenticationError("refresh token not registered in the store")
		} else {
			logger.Error("Unexpected database error: " + err.Error())
			return errs.NewUnexpectedError("unexpected database error")
		}
	}
	return nil
}

func (d AuthRepositoryDb) GenerateAndSaveRefreshTokenStore (authToken AuthToken) (string, *errs.AppError) {
	var appErr *errs.AppError
	var refreshToken string
	if refreshToken, appErr = authToken.newRefreshToken(); appErr != nil{
		return "",appErr
	}

	sqlInsert := "INSERT INTO refresh_token_store SET refresh_token = ?"
	_, err := d.client.Exec(sqlInsert, refreshToken)
	if err != nil {
		logger.Error("Unexpected database error: " + err.Error())
		return "", errs.NewUnexpectedError("unexpected database error")
	}
	return refreshToken, nil
}

func (d AuthRepositoryDb) FindBy(username, password string) (*Login, *errs.AppError){
	var login Login
	sqlVerify := `
	SELECT 
		username,
		u.customer.id,
		role,
		group_concat(a.account.id) as account_numbers
	FROM
		users as u
	LEFT JOIN accounts as a
		ON a.customer_id = u.customer.id
	WHERE
		username = ? AND password = ?
	GROUP BY a.customer_id
	`
	err := d.client.Get(&login, sqlVerify,username,password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewAuthenticationError("invalid credentials")
		} else {
			logger.Error("Error while verifying request from database: " + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected database error")
		}
	}
	return &login,nil
}

func NewAuthRepository(client *sql.DB) AuthRepository {
	return AuthRepository{client}
}