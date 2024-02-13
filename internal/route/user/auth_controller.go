package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mbocek/meet-go/internal/interfaces"
	webError "github.com/mbocek/meet-go/internal/route/error"
	"github.com/rotisserie/eris"
	"net/http"
)

type AuthenticationController struct {
	db  interfaces.DBRepository
	pwd *PasswordService
}

type signIn struct {
	Name     string
	Password string
}

type signUpRequest struct {
	Name     string `json:"name,omitempty" binding:"required,max=100"`
	Surname  string `json:"surname,omitempty" binding:"required,max=100"`
	Email    string `json:"email,omitempty" binding:"required,email,max=255"`
	Password string `json:"password,omitempty" binding:"required"`
}

func (s signUpRequest) buildUser(p *PasswordService) (*User, error) {
	hash, err := p.GenerateHash(s.Password, "")
	if err != nil {
		return nil, err
	}

	return &User{
		Name:         s.Name,
		Surname:      s.Surname,
		Email:        s.Email,
		PasswordHash: hash.HashB64,
		SaltHash:     hash.HashB64,
		Enabled:      false,
	}, nil
}

func NewAuthentication(db interfaces.DBRepository, pwd *PasswordService) *AuthenticationController {
	return &AuthenticationController{db: db, pwd: pwd}
}

func (a *AuthenticationController) Signin(ctx *gin.Context) {
	var signin signIn

	err := ctx.ShouldBindJSON(&signin)
	if err != nil {
		webError.ReturnBadRequestError(ctx, eris.Wrap(err, "cannot read body"))
		return
	}

	var user User

	sql := `SELECT * FROM "user" where email = $1`
	err = a.db.GetDb().Get(&user, sql, signin.Name)
	if webError.IsNoRowsFoundError(err) {
		webError.ReturnAuthenticationError(ctx, eris.New("user doesn't exists"))
		return
	} else if err != nil {
		webError.ReturnInternalServerError(ctx, eris.Wrap(err, "cannot select users"))
		return
	}
	err = a.pwd.Compare(user.PasswordHash, user.SaltHash, signin.Password)
	if err != nil {
		webError.ReturnAuthenticationError(ctx, eris.Wrap(err, "invalid password"))
		return
	}
	if !user.Enabled {
		webError.ReturnAuthenticationError(ctx, eris.New("user is not confirmed"))
		return
	}
	ctx.JSON(http.StatusOK, "")
}

func (a *AuthenticationController) Signup(ctx *gin.Context) {
	var signup signUpRequest

	err := ctx.ShouldBindJSON(&signup)
	if err != nil {
		var e validator.ValidationErrors
		if errors.As(err, &e) {
			webError.ReturnBadRequestValidatorError(ctx, e)
			return
		} else {
			webError.ReturnBadRequestError(ctx, err)
			return

		}
	}

	//check if user exists
	var existUser User
	query := `SELECT * FROM "user" WHERE email = $1`
	err = a.db.GetDb().Get(&existUser, query, signup.Email)
	if err != nil && !webError.IsNoRowsFoundError(err) {
		webError.ReturnInternalServerError(ctx, eris.Wrap(err, "cannot get user"))
		return
	} else if err == nil {
		webError.ReturnConflictError(ctx, eris.New("user already exists"))
		return
	}

	// insert user to database
	user, err := signup.buildUser(a.pwd)
	if err != nil {
		webError.ReturnBadRequestError(ctx, eris.Wrap(err, "cannot set user"))
		return
	}

	query = `INSERT INTO "user"(name, surname, email, password_hash, salt_hash, enabled) 
			  VALUES (:name, :surname, :email, :password_hash, :salt_hash, :enabled)`
	if _, err := a.db.GetDb().NamedExec(query, user); err != nil {
		webError.ReturnInternalServerError(ctx, eris.Wrap(err, "error while saving user to database"))
		return
	}

	ctx.JSON(http.StatusOK, "")
}
