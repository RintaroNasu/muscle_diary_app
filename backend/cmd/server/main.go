package main

import (
	"log/slog"

	"github.com/RintaroNasu/muscle_diary_app/cmd/migrate"
	"github.com/RintaroNasu/muscle_diary_app/internal/db"
	"github.com/RintaroNasu/muscle_diary_app/internal/httpx"
	"github.com/RintaroNasu/muscle_diary_app/internal/logging"
	"github.com/RintaroNasu/muscle_diary_app/routes"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	logger := logging.New()
	slog.SetDefault(logger)

	e.HTTPErrorHandler = httpx.HTTPErrorHandler(logger)

	// データベースに接続
	conn, err := db.New()
	if err != nil {
		e.Logger.Fatal("Failed to connect to database: ", err)
	}

	// データベースを閉じる
	defer func() {
		if err := db.CloseDB(conn); err != nil {
			e.Logger.Error("Failed to close database: ", err)
		}
	}()

	// マイグレーション
	if err := migrate.Migrate(conn); err != nil {
		e.Logger.Fatal("Failed to migrate database: ", err)
	}

	// seedの追加
	if err := db.Seed(conn); err != nil {
		e.Logger.Fatal("Failed to seed database: ", err)
	}

	// ルーティング
	routes.Register(e, conn)

	e.Logger.Fatal(e.Start(":8080"))
}
