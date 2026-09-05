package repository

import (
	"context"
	"saga/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProcessEventRepository struct {
	db *gorm.DB
}

func NewProcessEventRepository(db *gorm.DB) *ProcessEventRepository {
	return &ProcessEventRepository{
		db: db,
	}
}
func (r *ProcessEventRepository) IsProcessed(ctx context.Context, eventID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.ProcessedEvent{}).Where("event_id=?", eventID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil

}

func (r *ProcessEventRepository) MarkProcessed(ctx context.Context, eventID uuid.UUID) error {
	event := model.ProcessedEvent{
		EventID:     eventID,
		ProcessedAt: time.Now(),
	}

	return r.db.WithContext(ctx).Create(&event).Error

}
