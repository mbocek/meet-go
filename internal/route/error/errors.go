package error

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/mbocek/meet-go/internal/middleware"
	"net/http"
)

type ApiError struct {
	Error string `json:"error"`
}

func IsNoRowsFoundError(err error) bool {
	return err == sql.ErrNoRows
}

func ReturnNotFoundError(c *gin.Context, err error) {
	log := middleware.GetRequestLogger(c)
	log.Error().Stack().Err(err).Int("httpStatusCode", http.StatusNotFound).Msg("not found")
	c.JSON(http.StatusNotFound, ApiError{Error: err.Error()})
}

func ReturnInternalServerError(c *gin.Context, err error) {
	log := middleware.GetRequestLogger(c)
	log.Error().Stack().Err(err).Int("httpStatusCode", http.StatusInternalServerError).Msg("internal server error")
	c.JSON(http.StatusInternalServerError, ApiError{Error: err.Error()})
}
