package main

import (
	"github.com/RintaroNasu/muscle_diary_app/cmd/migrate"
	"github.com/RintaroNasu/muscle_diary_app/internal/db"
	"github.com/RintaroNasu/muscle_diary_app/routes"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

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

	// ルーティング
	routes.Register(e, conn)

	e.Logger.Fatal(e.Start(":8080"))
}
