package user_test

import (
	"context"
	_ "embed"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"github.com/mbocek/meet-go/db"
	apiError "github.com/mbocek/meet-go/internal/route/error"
	"github.com/mbocek/meet-go/internal/route/user"
	fixtures "github.com/mbocek/meet-go/internal/route/user/fixtures"
	"github.com/mbocek/meet-go/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func loginAssertions(t *testing.T, r *httptest.ResponseRecorder) {
	assert.Equal(t, http.StatusOK, r.Code)
}
func wrongPasswordAssertions(t *testing.T, r *httptest.ResponseRecorder) {
	assert.Equal(t, http.StatusUnauthorized, r.Code)
}

func notEnabledPasswordAssertions(t *testing.T, r *httptest.ResponseRecorder) {
	assert.Equal(t, http.StatusUnauthorized, r.Code)
}

func signupAssertions(t *testing.T, r *httptest.ResponseRecorder, db *sqlx.DB) {
	assert.Equal(t, http.StatusOK, r.Code)
	var actualUser user.User
	err := db.Get(&actualUser, `SELECT * FROM "user" WHERE email = $1`, "test@test.com")
	require.NoError(t, err)
	assert.Equal(t, "Name", actualUser.Name)
	assert.Equal(t, "Surname", actualUser.Surname)
	assert.Equal(t, "test@test.com", actualUser.Email)
}

func signupWrongEmailAssertions(t *testing.T, r *httptest.ResponseRecorder, db *sqlx.DB) {
	assert.Equal(t, http.StatusBadRequest, r.Code)
	var e apiError.ApiError
	err := json.Unmarshal(r.Body.Bytes(), &e)
	assert.Nil(t, err)
	assert.Contains(t, "validation error", e.Error)
	assert.Contains(t, "Email", e.Validation[0].Field)
	assert.Contains(t, "email", e.Validation[0].Tag)
}

func userAlreadyExistsAssertions(t *testing.T, r *httptest.ResponseRecorder, db *sqlx.DB) {
	assert.Equal(t, http.StatusConflict, r.Code)
	var e apiError.ApiError
	err := json.Unmarshal(r.Body.Bytes(), &e)
	assert.Nil(t, err)
	assert.Equal(t, "user already exists", e.Error)
}

var (
	authMigrationTable = "auth_migration"
)

func TestSignin(t *testing.T) {
	p := test.NewPostgres(t, context.Background())
	defer p.Close()

	testData := []struct {
		name       string
		path       string
		request    string
		assertions func(*testing.T, *httptest.ResponseRecorder)
	}{
		{name: "Signin", path: "/api/v1/auth/signin", request: `{"name":"test@email.com", "password": "test"}`, assertions: loginAssertions},
		{name: "Wrong password", path: "/api/v1/auth/signin", request: `{"name":"test@email.com", "password": "test1"}`, assertions: wrongPasswordAssertions},
		{name: "Not enabled", path: "/api/v1/auth/signin", request: `{"name":"test2@email.com", "password": "test"}`, assertions: notEnabledPasswordAssertions},
	}

	test.Migrate(t, p.Url(), db.Migrations, nil)
	test.Migrate(t, p.Url(), fixtures.UserMigration, &authMigrationTable)

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			gin := test.NewGin(t, p.Url())

			r, err := http.NewRequestWithContext(gin.Ctx, http.MethodPost, data.path, strings.NewReader(data.request))
			assert.Nil(t, err)
			gin.Engine.ServeHTTP(gin.Recorder, r)
			data.assertions(t, gin.Recorder)
		})
	}
}

func TestSignup(t *testing.T) {
	p := test.NewPostgres(t, context.Background())
	defer p.Close()

	testData := []struct {
		name       string
		path       string
		request    string
		assertions func(*testing.T, *httptest.ResponseRecorder, *sqlx.DB)
	}{
		{name: "Signup", path: "/api/v1/auth/signup", request: `{"name":"Name", "surname":"Surname", "email":"test@test.com", "password": "test"}`, assertions: signupAssertions},
		{name: "Signup wrong email", path: "/api/v1/auth/signup", request: `{"name":"Name", "surname":"Surname", "email":"testtest.com", "password": "test"}`, assertions: signupWrongEmailAssertions},
		{name: "Existing user", path: "/api/v1/auth/signup", request: `{"name":"Name", "surname":"Surname", "email":"test@email.com", "password": "test"}`, assertions: userAlreadyExistsAssertions},
	}

	test.Migrate(t, p.Url(), db.Migrations, nil)
	test.Migrate(t, p.Url(), fixtures.UserMigration, &authMigrationTable)

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			gin := test.NewGin(t, p.Url())

			r, err := http.NewRequestWithContext(gin.Ctx, http.MethodPost, data.path, strings.NewReader(data.request))
			assert.Nil(t, err)
			gin.Engine.ServeHTTP(gin.Recorder, r)
			data.assertions(t, gin.Recorder, gin.DB)
		})
	}
}
