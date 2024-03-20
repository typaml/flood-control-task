package database

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

const (
	K = 5  //Максимальное кол-во вызовов метода Check в течение N секунд
	N = 60 //Последние N секунд вызовов метода Check
)

type SQliteFloodControl struct {
	Db *sql.DB
}

func NewSQliteFloodControl(dbPath string, ctx context.Context) (*SQliteFloodControl, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS flood_control (id INTEGER NOT NULL PRIMARY KEY, user_id INTEGER, call_time INTEGER)")
	if err != nil {
		return nil, err
	}
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	return &SQliteFloodControl{Db: db}, nil
}

func (fc *SQliteFloodControl) Check(ctx context.Context, userID int64) (bool, error) {
	timeRequest := time.Now().Unix()
	defer fc.Db.Close()

	// Удаляем записи старше N+1 секунд
	_, err := fc.Db.ExecContext(ctx, "DELETE FROM flood_control WHERE call_time < ?", timeRequest-N+1)
	if err != nil {
		log.Printf("cannot delete records %v", err)
	}
	// Проверяем кол-во вызовов за последние N секунд для конкретного user_id
	rows, err := fc.Db.QueryContext(ctx, "SELECT COUNT(*) FROM flood_control WHERE user_id = ? AND call_time > ?", userID, timeRequest-N)
	if err != nil {
		return false, err
	}

	// Проверяем, что кол-во вызовов не превышает K
	var count int
	if rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return false, err
		}
	}
	rows.Close()

	if count > K {
		return false, nil
	}

	// Добавляем запись в таблицу
	_, err = fc.Db.ExecContext(ctx, "INSERT INTO flood_control (user_id, call_time) VALUES (?, ?)", userID, timeRequest)
	if err != nil {
		return false, err
	}
	return true, nil
}
