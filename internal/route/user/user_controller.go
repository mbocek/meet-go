package user

import (
	"github.com/gin-gonic/gin"
	"github.com/mbocek/meet-go/internal/interfaces"
	"github.com/mbocek/meet-go/internal/route/error"
	"github.com/rotisserie/eris"
	"net/http"
)

type Users []User
type User struct {
	Id           int    `db:"id"`
	Name         string `db:"name"`
	Surname      string `db:"surname"`
	Email        string `db:"email"`
	PasswordHash string `db:"password_hash"`
	SaltHash     string `db:"salt_hash"`
}

type UserController struct {
	controller *gin.Context
	db         interfaces.DBRepository
}

func NewUserController(db interfaces.DBRepository) *UserController {
	return &UserController{
		db: db,
	}
}

func (c *UserController) GetAll(ctx *gin.Context) {
	users := make(Users, 0)
	sql := `SELECT id, name, surname, email FROM "user"`
	err := c.db.GetDb().Select(&users, sql)
	if error.IsNoRowsFoundError(err) {
		ctx.JSON(http.StatusOK, users)
		return
	} else if err != nil {
		error.ReturnInternalServerError(ctx, eris.Wrap(err, "cannot select users"))
		return
	}
	ctx.JSON(http.StatusOK, users)
}
