package domain

import (
	json2 "encoding/json"
	"github.com/barnettt/banking/exceptions"
	"github.com/barnettt/banking/logger"
	"github.com/jmoiron/sqlx"
	"net/http"
	"net/url"
	"os"
)

type AuthorisationRepositoryDB struct {
	Client *sqlx.DB
}

type AuthorisationRepository interface {
	IsUserAuthorised(token string, currentRoute string, vars map[string]string) (bool, *exceptions.AppError)
}
type AuthResponse struct {
	Code    int
	Message string
}
type RemoteAuthRepository struct {
}

func (repository AuthorisationRepositoryDB) IsUserAuthorised(token string, currentRoute string, vars map[string]string) (bool, *exceptions.AppError) {
	authUrl := buildVerifyUrl(token, currentRoute, vars)
	var res bool
	if response, err := http.Get(authUrl); err != nil {
		logger.Error(err.Error())
		appErr := exceptions.AppError{Code: http.StatusUnauthorized, Message: err.Error()}
		return false, &appErr
	} else {
		if response.StatusCode != 200 {
			var errorString AuthResponse
			// parse as an app error and handle parse error
			err := json2.NewDecoder(response.Body).Decode(&errorString)
			if err != nil {
				logger.Error(err.Error())
				appErr := exceptions.AppError{Code: http.StatusInternalServerError, Message: err.Error()}
				return false, &appErr
			}
			// parse ok return http response error
			logger.Error("Error while decoding authorisation server response : " + errorString.Message)
			appErr := exceptions.AppError{Code: response.StatusCode, Message: errorString.Message}
			return false, &appErr
		} else {
			// parse for as a boolean as status code is 200
			// handle the parse error
			err := json2.NewDecoder(response.Body).Decode(&res)
			if err != nil {
				logger.Error(err.Error())
				appErr := exceptions.AppError{Code: http.StatusInternalServerError, Message: err.Error()}
				return false, &appErr
			}
		}
		return res, nil
	}
	return false, nil
}

func buildVerifyUrl(token string, route string, vars map[string]string) string {
	port := os.Getenv("AUTH_SERVER_PORT")
	host := os.Getenv("SERVER_HOST")
	url := url.URL{Host: host + ":" + port, Path: "auth/verify", Scheme: "http"}
	params := url.Query()
	params.Add("token", token)
	params.Add("operation", route)
	for k, v := range vars {
		params.Add(k, v)
	}
	url.RawQuery = params.Encode()
	logger.Info(url.String())
	return url.String()

}
