package user_test

import (
	"context"
	_ "embed"
	"github.com/mbocek/meet-go/db"
	fixtures "github.com/mbocek/meet-go/internal/route/user/fixtures"
	"github.com/mbocek/meet-go/internal/test"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func loginAssertions(t *testing.T, r *httptest.ResponseRecorder) {
	assert.Equal(t, http.StatusOK, r.Code)
}

var (
	authMigrationTable = "auth_migration"
)

func TestLogin(t *testing.T) {
	p := test.NewPostgres(t, context.Background())
	defer p.Close()

	testData := []struct {
		name       string
		path       string
		request    string
		assertions func(*testing.T, *httptest.ResponseRecorder)
	}{
		{name: "Login", path: "/api/v1/login", request: `{"name":"test", "password": "test"}`, assertions: loginAssertions},
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
