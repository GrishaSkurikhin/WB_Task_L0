package mapcache

import (
	customerrors "github.com/GrishaSkurikhin/WB_Task_L0/internal/custom-errors"
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
	if _, isExist := mc.m[order.OrderUID.String()]; isExist {
		return customerrors.OrderAlreadyExist{}
	} else {
		mc.m[order.OrderUID.String()] = order
		return nil
	}
}

func (mc *MapCache) GetOrder(orderUUID uuid.UUID) (models.Order, error) {
	if order, isExist := mc.m[orderUUID.String()]; isExist {
		return order, nil
	} else {
		return models.Order{}, customerrors.OrderNotFound{}
	}
}

func (mc *MapCache) GetOrdersID() []string {
	ordersID := make([]string, 0, len(mc.m))
	for id := range mc.m {
		ordersID = append(ordersID, id)
	}
	return ordersID
}
