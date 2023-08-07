package route

import (
	"github.com/gin-gonic/gin"
	"github.com/mbocek/meet-go/internal/interfaces"
	"github.com/mbocek/meet-go/internal/middleware"
	"github.com/mbocek/meet-go/internal/route/user"
)

type RestService struct {
	userController *user.Controller
}

type HandlerFuncWithDBRepo func(*gin.Context, interfaces.DBRepository)

func New(dbRepo interfaces.DBRepository) *RestService {
	return &RestService{
		userController: user.New(dbRepo),
	}
}

func (r *RestService) RegisterHandlers(e *gin.Engine) {
	e.Use(gin.Recovery()) // recover from panics, send 500 instead!
	e.Use(middleware.RouteLoggerMW())
	apiRoutes := e.Group("/api")

	apiRoutes.GET("/user", r.userController.GetAll)
}
