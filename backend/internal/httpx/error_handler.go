package httpx

import (
	"log/slog"

	"github.com/labstack/echo/v4"
)

func HTTPErrorHandler(l *slog.Logger) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		ctx := c.Request().Context()

		if ae, ok := err.(*AppError); ok {
			if ae.Status >= 500 {
				l.ErrorContext(ctx, "server_error",
					"code", ae.Code, "err", ae.Err,
					"path", c.Path(), "method", c.Request().Method, "status", ae.Status,
				)
			} else {
				l.WarnContext(ctx, "client_error",
					"code", ae.Code, "err", ae.Err,
					"path", c.Path(), "method", c.Request().Method, "status", ae.Status,
				)
			}
			_ = c.JSON(ae.Status, map[string]string{
				"code": ae.Code, "message": ae.Message,
			})
			return
		}

		l.ErrorContext(ctx, "unexpected_error",
			"err", err, "path", c.Path(), "method", c.Request().Method,
		)
		_ = c.JSON(500, map[string]string{
			"code": "InternalError", "message": "サーバでエラーが発生しました",
		})
	}
}
