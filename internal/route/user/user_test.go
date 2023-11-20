package user_test

import (
	"context"
	"github.com/mbocek/meet-go/db"
	"github.com/mbocek/meet-go/internal/test"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func emptyUsers(t *testing.T, r *httptest.ResponseRecorder) {
	assert.Equal(t, http.StatusOK, r.Code)
	assert.JSONEq(t, `[]`, r.Body.String())
}

func TestUsers(t *testing.T) {
	p := test.NewPostgres(t, context.Background())
	defer p.Close()

	testData := []struct {
		name       string
		path       string
		migrate    test.MigrateTest
		assertions func(*testing.T, *httptest.ResponseRecorder)
	}{
		{name: "EmptyUsers", path: "/api/v1/users", assertions: emptyUsers},
	}

	test.Migrate(t, p.Url(), db.Migrations, nil)

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
