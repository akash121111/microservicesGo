package database

import (
	"fmt"
	"order/config"
	"order/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabaseCongig(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", cfg.POSTGRES_HOST, cfg.POSTGRES_USER, cfg.POSTGRES_PASSWORD, cfg.POSTGRES_DB, cfg.POSTGRES_PORT)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, err
}

func MigrateDB(db *gorm.DB) error {
	return db.AutoMigrate(&model.Order{}, &model.OrderItem{}, &model.OutboxEvent{})
}
