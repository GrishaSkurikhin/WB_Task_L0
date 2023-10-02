package orders

import (
	"fmt"

	"github.com/GrishaSkurikhin/WB_Task_L0/internal/models"
	"github.com/google/uuid"
)

type CacheAdder interface {
	AddOrder(order models.Order) error
}

type StorageAdder interface {
	AddOrder(order models.Order) error
}

func Add(order models.Order, cache CacheAdder, strg StorageAdder) error {
	const op = "orders.Add"

	// TODO: additional order logic

	err := cache.AddOrder(order)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = strg.AddOrder(order)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

type CacheGetter interface {
	GetOrder(orderUUID uuid.UUID) (models.Order, error)
	GetOrdersID() []string
}

func Get(orderUUID uuid.UUID, cache CacheGetter) (models.Order, error) {
	const op = "orders.Get"

	order, err := cache.GetOrder(orderUUID)
	if err != nil {
		return models.Order{}, fmt.Errorf("%s: %w", op, err)
	}
	return order, nil
}

func GetIDs(cache CacheGetter) []string {
	return cache.GetOrdersID()
}

type StorageGetter interface {
	GetAllOrders() ([]models.Order, error)
}

func AddAllToCache(cache CacheAdder, strg StorageGetter) error {
	const op = "orders.AddAllToCache"

	orders, err := strg.GetAllOrders()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, order := range orders {
		err := cache.AddOrder(order)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}
