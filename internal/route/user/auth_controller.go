package user

import (
	"github.com/gin-gonic/gin"
	"github.com/mbocek/meet-go/internal/interfaces"
	"github.com/mbocek/meet-go/internal/route/error"
	"github.com/rotisserie/eris"
	"net/http"
)

type AuthenticationController struct {
	db  interfaces.DBRepository
	pwd *PasswordService
}

func NewAuthentication(db interfaces.DBRepository, pwd *PasswordService) *AuthenticationController {
	return &AuthenticationController{db: db, pwd: pwd}
}

func (a *AuthenticationController) Login(ctx *gin.Context) {
	login := struct {
		Name     string
		Password string
	}{}

	err := ctx.BindJSON(&login)
	if err != nil {
		error.ReturnBadRequestError(ctx, eris.Wrap(err, "cannot read body"))
		return
	}

	var user User

	sql := `SELECT * FROM "user" where name = $1`
	err = a.db.GetDb().Get(&user, sql, login.Name)
	if error.IsNoRowsFoundError(err) {
		error.ReturnAuthenticationError(ctx, eris.New("user doesn't exists"))
		return
	} else if err != nil {
		error.ReturnInternalServerError(ctx, eris.Wrap(err, "cannot select users"))
		return
	}
	err = a.pwd.Compare(user.PasswordHash, user.SaltHash, login.Password)
	if err != nil {
		error.ReturnAuthenticationError(ctx, eris.Wrap(err, "invalid password"))
		return
	}
	ctx.JSON(http.StatusOK, nil)
}
