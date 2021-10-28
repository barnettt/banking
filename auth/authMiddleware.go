package auth

import (
	"github.com/barnettt/banking/domain"
	"github.com/barnettt/banking/exceptions"
	"github.com/barnettt/banking/logger"
	"github.com/barnettt/banking/util"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

type AuthorisationMiddleware struct {
	Repository domain.AuthorisationRepository
}

func (authMiddleware AuthorisationMiddleware) AuthorisationHandler() func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			currentRouteVars := mux.Vars(request)
			currentRoute := mux.CurrentRoute(request)
			authHeader := request.Header.Get("Authorization")
			contentType := request.Header.Get("Content-Type")
			if authHeader != "" {
				token := getTokenFromHeader(authHeader)

				isAuthorised, appErr := authMiddleware.Repository.IsUserAuthorised(token, currentRoute.GetName(), currentRouteVars)
				if isAuthorised {
					next.ServeHTTP(response, request)
				} else if appErr != nil {
					appErr := exceptions.AppError{Code: http.StatusUnauthorized, Message: appErr.Message}
					util.WriteResponse(response, http.StatusUnauthorized, appErr, contentType)
					return
				}
			} else {
				appErr := exceptions.AppError{Code: http.StatusForbidden, Message: "User not authenticated"}
				util.WriteResponse(response, http.StatusForbidden, appErr, contentType)
			}
		})
	}

}

func NewAuthorisationMiddleware(repo domain.AuthorisationRepository) AuthorisationMiddleware {
	return AuthorisationMiddleware{Repository: repo}
}
func getTokenFromHeader(header string) string {
	jwtToken := strings.Split(header, " ")
	logger.Info(jwtToken[0] + "   " + jwtToken[1])
	return jwtToken[1]
}
