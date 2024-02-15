package db

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Expression struct {
	ID          uint `gorm:"primaryKey"`
	Expression  string
	Result      string
	Status      string `gorm:"default:'in_progress'"`
	CreatedAt   time.Time
	EvaluatedAt time.Time
}

func ConnectToPostgreSQL(host, port, user, password, dbname string) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		user,
		password,
		dbname,
		host,
		port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
