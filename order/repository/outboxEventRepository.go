package repository

import (
	"context"
	"fmt"
	"order/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OutboxEventRepository struct {
	db *gorm.DB
}

func NewOutboxEventRepository(db *gorm.DB) *OutboxEventRepository {
	return &OutboxEventRepository{
		db: db,
	}
}

func (r *OutboxEventRepository) WithTx(tx *gorm.DB) *OutboxEventRepository {
	return &OutboxEventRepository{
		db: tx,
	}
}

func (r *OutboxEventRepository) CreateOutBoxEvent(ctx context.Context, event model.OutboxEvent) error {
	err := r.db.WithContext(ctx).Create(&event).Error
	if err != nil {
		return fmt.Errorf("error in creating outboxEvent %w", err)
	}
	return nil
}

func (r *OutboxEventRepository) GetAllNotPublishOutBoxEvent(ctx context.Context) ([]model.OutboxEvent, error) {
	var event []model.OutboxEvent
	err := r.db.WithContext(ctx).Where("published=?", false).Order("created_at ASC").Find(&event).Error
	if err != nil {
		return nil, fmt.Errorf(
			"error getting unpublished outbox events: %w",
			err,
		)
	}

	return event, nil
}

func (r *OutboxEventRepository) GetOutBoxEventById(ctx context.Context, id uuid.UUID) (*model.OutboxEvent, error) {
	var event model.OutboxEvent
	err := r.db.WithContext(ctx).Where("id=?", id).First(&event).Error
	if err != nil {
		return nil, fmt.Errorf("error in getting outboxEvent by ID %w ", err)
	}

	return &event, nil
}

func (r *OutboxEventRepository) UpdateBoxEventPublishById(
	ctx context.Context,
	id uuid.UUID,
) error {

	now := time.Now()

	result := r.db.WithContext(ctx).
		Model(&model.OutboxEvent{}).
		Where("id = ? AND published = ?", id, false).
		Updates(map[string]interface{}{
			"published":    true,
			"published_at": now,
		})

	if result.Error != nil {
		return fmt.Errorf(
			"error updating outbox event: %w",
			result.Error,
		)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf(
			"outbox event not found or already published",
		)
	}

	return nil
}
