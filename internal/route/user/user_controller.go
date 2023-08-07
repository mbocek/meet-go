package user

import (
	"github.com/gin-gonic/gin"
	"github.com/mbocek/meet-go/internal/interfaces"
	error2 "github.com/mbocek/meet-go/internal/route/error"
	"github.com/rotisserie/eris"
	"net/http"
)

type Users []User
type User struct {
	Id      int
	Name    string
	Surname string
	Email   string
}

type Controller struct {
	controller *gin.Context
	db         interfaces.DBRepository
}

func New(db interfaces.DBRepository) *Controller {
	return &Controller{
		db: db,
	}
}

func (c *Controller) GetAll(ctx *gin.Context) {
	users := make(Users, 0)
	sql := `SELECT id, name, surname, email FROM "user"`
	err := c.db.GetDb().Select(&users, sql)
	if error2.IsNoRowsFoundError(err) {
		ctx.JSON(http.StatusOK, users)
		return
	} else if err != nil {
		error2.ReturnInternalServerError(ctx, eris.Wrap(err, "cannot select users"))
		return
	}
	ctx.JSON(http.StatusOK, users)
}
