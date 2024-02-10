package test

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mbocek/meet-go/internal/config"
	"github.com/mbocek/meet-go/internal/repository/postgres"
	"github.com/mbocek/meet-go/internal/route/user"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"testing"
)

type GinTestContext struct {
	t        *testing.T
	Ctx      *gin.Context
	Engine   *gin.Engine
	Recorder *httptest.ResponseRecorder
}

func NewGin(t *testing.T, dbURI string) *GinTestContext {
	t.Helper()
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)

	dbRepo, err := postgres.New(context.Background(), config.Postgres{Url: dbURI})
	require.NoError(t, err)

	routes := user.NewRoute(dbRepo)
	routes.RegisterHandlers(r)

	return &GinTestContext{
		t:        t,
		Ctx:      ctx,
		Engine:   r,
		Recorder: w,
	}
}
