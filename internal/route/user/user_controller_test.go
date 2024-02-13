package user_test

import (
	"context"
	_ "embed"
	"encoding/json"
	"github.com/mbocek/meet-go/db"
	"github.com/mbocek/meet-go/internal/route/user"
	fixtures "github.com/mbocek/meet-go/internal/route/user/fixtures"
	"github.com/mbocek/meet-go/internal/test"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func oneUserAssertions(t *testing.T, r *httptest.ResponseRecorder) {
	assert.Equal(t, http.StatusOK, r.Code)
	var u []user.User
	err := json.Unmarshal(r.Body.Bytes(), &u)
	assert.Nil(t, err)
	assert.Greater(t, len(u), 1)
}

var (
	userMigrationTable = "user_migration"

	//go:embed fixtures/test/one-user.json
	oneUser string
)

func TestUsers(t *testing.T) {
	p := test.NewPostgres(t, context.Background())
	defer p.Close()

	testData := []struct {
		name       string
		path       string
		migrate    test.MigrateTest
		assertions func(*testing.T, *httptest.ResponseRecorder)
	}{
		{name: "One user", path: "/api/v1/users", assertions: oneUserAssertions},
	}

	test.Migrate(t, p.Url(), db.Migrations, nil)
	test.Migrate(t, p.Url(), fixtures.UserMigration, &userMigrationTable)

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			gin := test.NewGin(t, p.Url())

			r, err := http.NewRequestWithContext(gin.Ctx, http.MethodGet, data.path, nil)
			assert.Nil(t, err)
			gin.Engine.ServeHTTP(gin.Recorder, r)
			data.assertions(t, gin.Recorder)
		})
	}
}
