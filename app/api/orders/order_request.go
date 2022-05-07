package orders

import (
	"app/internal/config"
	"app/internal/orders"
	"encoding/json"
	"io"
	"math/rand"
	"strconv"
	"time"
)

func FromJSON(reader io.Reader) (OrderRequest, error) {
	decoder := json.NewDecoder(reader)
	var order OrderRequest
	err := decoder.Decode(&order)
	if err != nil {
		return OrderRequest{}, err
	}
	return order, nil
}

func (order *OrderRequest) ToOrderEntity(region config.Region, environment config.Environment) orders.OrderEntity {
	creationDate := time.Now()
	orderId := orders.GenerateOrderId(region, environment, creationDate, strconv.Itoa(rand.Int()))

	var orderItems []orders.OrderItemEntity
	for _, item := range order.Items {
		orderItems = append(orderItems, item.ToOrderItemEntity(orderId, creationDate))
	}

	return orders.OrderEntity{
		Id:           orderId,
		Workflow:     "default_workflow",
		CreationDate: creationDate,
		Status:       orders.OrderPlaced,
		Items:        orderItems,
	}
}
