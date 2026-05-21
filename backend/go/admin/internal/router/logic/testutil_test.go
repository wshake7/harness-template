package logic

import (
	"admin/internal/auth"
	"admin/internal/fiberc/handler"
	"admin/internal/services/orm/query"
	"os"
	"testing"

	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	db, err := gorm.Open(sqlite.Open(":memory:"))
	if err == nil {
		query.SetDefault(db)
	}
	os.Exit(m.Run())
}

func newTestCtx(t *testing.T) *handler.Ctx {
	t.Helper()
	ctx := &handler.Ctx{
		SessionInfo: &auth.SessionInfo{
			Id:       1,
			Username: "testuser",
		},
	}
	ctx.SetLogger(zap.NewNop())
	return ctx
}
