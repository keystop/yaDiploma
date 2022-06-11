package database

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"

	"github.com/keystop/yaDiploma/internal/config"
	"github.com/keystop/yaDiploma/internal/models"
	"github.com/keystop/yaDiploma/pkg/logger"
)

type serverDB struct {
	*sql.DB
}

//CheckDBConnection trying connect to db.
func (s *serverDB) CheckDBConnection(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	err := s.PingContext(ctx)
	if err != nil {
		logger.Error("Ошибка подключения к БД", err)
	}
}

// makeMigrations start here for autotests
// func (s *serverDB) makeMigrations() {
// 	p := "Миграции базы данных:"
// 	logger.Info(p, "Старт")
// 	// setup database
// 	logger.Info(p, "Применение миграций")
// 	if err := goose.Up(s.DB, "../../db/migrations"); err != nil {
// 		logger.Error(p, err)
// 	}
// 	logger.Info(p, "Завершение") // run app
// }

func OpenDBConnect() models.ServerDB {
	s := new(serverDB)
	ctx := context.Background()
	db, err := sql.Open("postgres", config.Cfg.DBConnString())
	if err != nil {
		logger.Error("Ошибка подключения к БД", err)
	}
	s.DB = db
	s.CheckDBConnection(ctx)
	s.createTables(ctx)

	// s.makeMigrations()
	return s
}

func (s *serverDB) Close() {
	s.DB.Close()
}
