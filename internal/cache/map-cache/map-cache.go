package mapcache

import (
	"fmt"

	"github.com/GrishaSkurikhin/WB_Task_L0/internal/models"
	"github.com/google/uuid"
)

type MapCache struct {
	m map[string]models.Order
}

func New() *MapCache {
	return &MapCache{
		m: make(map[string]models.Order),
	}
}

func (mc *MapCache) AddOrder(order models.Order) error {
	const op = "mapcache.AddOrder"

	if _, isExist := mc.m[order.OrderUID.String()]; isExist {
		return fmt.Errorf("%s: %s", op, "order is already in cache")
	} else {
		mc.m[order.OrderUID.String()] = order
		return nil
	}
}

func (mc *MapCache) GetOrder(orderUUID uuid.UUID) (models.Order, error) {
	const op = "mapcache.GetOrder"

	if order, isExist := mc.m[orderUUID.String()]; isExist {
		return order, nil
	} else {
		return models.Order{}, fmt.Errorf("%s: %s", op, "cant find order with specified uid")
	}
}
