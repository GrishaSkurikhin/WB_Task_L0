package mapcache

import (
	"fmt"

	"github.com/GrishaSkurikhin/WB_Task_L0/internal/models"
	"github.com/google/uuid"
)

type MapCache struct {
	m map[uuid.UUID]models.Order
}

func (mc *MapCache) AddOrder(order models.Order) error {
	mc.m[order.OrderUID] = order
	return nil
}

func (mc *MapCache) GetOrder(orderUUID uuid.UUID) (models.Order, error) {
	if order, isExist := mc.m[orderUUID]; isExist {
		return order, nil
	} else {
		return models.Order{}, fmt.Errorf("cant find order with specified uid")
	}
}