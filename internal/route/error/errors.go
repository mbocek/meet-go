package error

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mbocek/meet-go/internal/middleware"
	"net/http"
)

type ApiError struct {
	Error      string                      `json:"error"`
	Validation []ApiRequestValidationError `json:"validation"`
}

type ApiRequestValidationError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
}

func IsNoRowsFoundError(err error) bool {
	return err == sql.ErrNoRows
}

func ReturnNotFoundError(c *gin.Context, err error) {
	log := middleware.GetRequestLogger(c)
	log.Error().Stack().Err(err).Int("httpStatusCode", http.StatusNotFound).Msg("not found")
	c.AbortWithStatusJSON(http.StatusNotFound, ApiError{Error: err.Error()})
}

func ReturnAuthenticationError(c *gin.Context, err error) {
	log := middleware.GetRequestLogger(c)
	log.Error().Stack().Err(err).Int("httpStatusCode", http.StatusUnauthorized).Msg("authentication")
	c.AbortWithStatusJSON(http.StatusUnauthorized, ApiError{})
}

func ReturnBadRequestError(c *gin.Context, err error) {
	log := middleware.GetRequestLogger(c)
	log.Error().Stack().Err(err).Int("httpStatusCode", http.StatusBadRequest).Msg("bad request")
	c.AbortWithStatusJSON(http.StatusBadRequest, ApiError{Error: err.Error()})
}

func ReturnBadRequestValidatorError(c *gin.Context, err validator.ValidationErrors) {
	log := middleware.GetRequestLogger(c)
	log.Error().Stack().Err(err).Int("httpStatusCode", http.StatusBadRequest).Msg("bad request")
	apiError := ApiError{Error: "validation error"}
	for _, e := range err {
		apiError.Validation = append(apiError.Validation, ApiRequestValidationError{
			Field: e.Field(),
			Tag:   e.Tag(),
		})
	}
	c.AbortWithStatusJSON(http.StatusBadRequest, apiError)
}

func ReturnInternalServerError(c *gin.Context, err error) {
	log := middleware.GetRequestLogger(c)
	log.Error().Stack().Err(err).Int("httpStatusCode", http.StatusInternalServerError).Msg("internal server error")
	c.AbortWithStatusJSON(http.StatusInternalServerError, ApiError{Error: err.Error()})
}

func ReturnConflictError(c *gin.Context, err error) {
	log := middleware.GetRequestLogger(c)
	log.Error().Stack().Err(err).Int("httpStatusCode", http.StatusConflict).Msg("conflict")
	c.AbortWithStatusJSON(http.StatusConflict, ApiError{Error: err.Error()})
}
