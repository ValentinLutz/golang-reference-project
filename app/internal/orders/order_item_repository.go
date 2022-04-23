package orders

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type OrderItemRepository struct {
	logger *zerolog.Logger
	db     *sqlx.DB
}

func NewOrderItemRepository(logger *zerolog.Logger, database *sqlx.DB) *OrderItemRepository {
	return &OrderItemRepository{logger: logger, db: database}
}

func (orderItemRepository *OrderItemRepository) FindAll() ([]OrderItemEntity, error) {
	rows, err := orderItemRepository.db.Query("SELECT id, order_id, creation_date, item_name FROM golang_reference_project.order_item")
	if err != nil {
		return nil, err
	}

	var orderItemEntities []OrderItemEntity
	for rows.Next() {
		var orderItemEntity OrderItemEntity

		err := rows.Scan(&orderItemEntity.Id, &orderItemEntity.OrderId, &orderItemEntity.CreationDate, &orderItemEntity.Name)
		if err != nil {
			return nil, err
		}

		orderItemEntities = append(orderItemEntities, orderItemEntity)
	}

	return orderItemEntities, nil
}

func (orderItemRepository *OrderItemRepository) SaveAll(orderItemEntities []OrderItemEntity) {
	_, err := orderItemRepository.db.NamedExec(
		`INSERT INTO golang_reference_project.order_item (order_id, creation_date, item_name) VALUES (:order_id, :creation_date, :item_name)`, orderItemEntities)
	if err != nil {
		orderItemRepository.logger.Error().
			Err(err).
			Msg("Failed to save order item entities into order item table")
	}
}
