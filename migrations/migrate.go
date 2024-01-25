package migrations

import (
	"embed"
	"fmt"

	"github.com/1boombacks1/testTaskInfotecs/internal/config"
	log "github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var embedMigrations embed.FS

func init() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config err: %s", err)
	}
	fmt.Println(cfg.Postgres.DSN + "?sslmode=disable")
	db, err := goose.OpenDBWithDriver("postgres", cfg.Postgres.DSN+"?sslmode=disable")
	if err != nil {
		log.Fatal("Не удалось подключиться к бд:", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal("Не удалось закрыть бд", err)
		}
	}()

	goose.SetBaseFS(embedMigrations)
	if err = goose.SetDialect("postgres"); err != nil {
		log.Fatal("Не удалось поставить диалект", err)
	}

	if err = goose.Up(db, "."); err != nil {
		log.Fatal("Не удалось поднять миграции\n", err)
	}

	log.Info("Миграции успешно поднялись!")
}
