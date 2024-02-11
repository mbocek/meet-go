package user

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mbocek/meet-go/internal/interfaces"
	"github.com/mbocek/meet-go/internal/middleware"
)

type RestService struct {
	userController *UserController
	authController *AuthenticationController
}

func NewRoute(dbRepo interfaces.DBRepository) *RestService {
	pwdService := NewPasswordService(3, 32, 64*1024, 2, 32)

	return &RestService{
		userController: NewUserController(dbRepo),
		authController: NewAuthentication(dbRepo, pwdService),
	}
}

func (r *RestService) RegisterHandlers(e *gin.Engine) {
	e.Use(gin.Recovery()) // recover from panics, send 500 instead!
	e.Use(middleware.RouteLoggerMW())
	e.Use(cors.Default())
	apiRoutes := e.Group("/api/v1")

	apiRoutes.GET("/users", r.userController.GetAll)
	apiRoutes.POST("/signin", r.authController.Login)
}
