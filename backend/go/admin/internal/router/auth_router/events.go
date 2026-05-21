package auth_router

import (
	"admin/internal/domains"
	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/middleware"
	"admin/internal/lifecycle"
	"bufio"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

func registerEventRouters(router fiber.Router) {
	router.Get("/events", func(ctx fiber.Ctx) error {
		c := handler.Trans(ctx)
		aesKey, _, err := middleware.DecryptRequest(c)
		if err != nil {
			return err
		}

		ctx.Set("Content-Type", "text/event-stream")
		ctx.Set("Cache-Control", "no-cache")
		ctx.Set("Connection", "keep-alive")
		ctx.Set("Transfer-Encoding", "chunked")
		ctx.Set(domains.HeaderXResponseIsEncrypt, "true")

		return ctx.SendStreamWriter(func(w *bufio.Writer) {
			ticker := time.NewTicker(time.Minute)
			defer ticker.Stop()
			count := 0
			for {
				select {
				case <-lifecycle.ShutdownDone():
					c.L().Debug("sse stream closed by server shutdown")
					return
				case <-ticker.C:
				}
				count++
				plainData, err := json.Marshal(map[string]int{"count": count})
				if err != nil {
					c.L().Error("marshal sse event failed", zap.Error(err))
					return
				}
				encryptedData, err := middleware.EncryptText(string(plainData), aesKey)
				if err != nil {
					c.L().Error("encrypt sse event failed", zap.Error(err))
					return
				}
				eventData, err := json.Marshal(map[string]string{"payload": encryptedData})
				if err != nil {
					c.L().Error("marshal encrypted sse event failed", zap.Error(err))
					return
				}
				if _, err = fmt.Fprintf(w, "event: count\ndata: %s\n\n", eventData); err != nil {
					c.L().Debug("write sse event failed", zap.Error(err))
					return
				}
				if err = w.Flush(); err != nil {
					c.L().Debug("sse client disconnected", zap.Error(err))
					return
				}
			}
		})
	})
}
