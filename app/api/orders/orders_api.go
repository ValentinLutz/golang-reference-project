package orders

import (
	"app/api/responses"
	"app/internal"
	"app/internal/orders"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"net/http"
)

type API struct {
	logger              *zerolog.Logger
	db                  *sqlx.DB
	orderRepository     *orders.OrderRepository
	orderItemRepository *orders.OrderItemRepository
	config              *internal.Config
}

func NewAPI(logger *zerolog.Logger, db *sqlx.DB, config *internal.Config) *API {
	return &API{
		logger:              logger,
		db:                  db,
		orderRepository:     orders.NewOrderRepository(logger, db),
		orderItemRepository: orders.NewOrderItemRepository(logger, db),
		config:              config,
	}
}

func (orderApi *API) RegisterHandlers(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/api/orders", orderApi.getOrders)
	router.HandlerFunc(http.MethodPost, "/api/orders", orderApi.postOrder)
	router.HandlerFunc(http.MethodGet, "/api/orders/:orderId", orderApi.getOrder)
}

func (orderApi *API) getOrders(responseWriter http.ResponseWriter, request *http.Request) {
	orderEntities, err := orderApi.orderRepository.FindAll()
	if err != nil {
		responses.Error(responseWriter, request, http.StatusInternalServerError, 100, err.Error())
	}
	orderItemEntities, err := orderApi.orderItemRepository.FindAll()
	if err != nil {
		responses.Error(responseWriter, request, http.StatusInternalServerError, 100, err.Error())
	}
	var ordersResponse OrdersResponse
	for _, order := range orderEntities {
		for _, orderItem := range orderItemEntities {
			if order.Id == orderItem.OrderId {
				order.Items = append(order.Items, orderItem)
				//sliceLen := len(orderItemEntities) - 1
				//orderItemEntities[i] = orderItemEntities[sliceLen]
				//orderItemEntities = orderItemEntities[:sliceLen]
			}
		}
		ordersResponse = append(ordersResponse, FromOrderEntity(&order))
	}

	responses.StatusOK(responseWriter, request, &ordersResponse)
}

func (orderApi *API) postOrder(responseWriter http.ResponseWriter, request *http.Request) {
	orderRequest, err := FromJSON(request.Body)
	if err != nil {
		responses.Error(responseWriter, request, http.StatusBadRequest, 200, err.Error())
		return
	}
	orderEntity := orderRequest.ToOrderEntity(orderApi.config.Region, orderApi.config.Environment)
	orderApi.orderRepository.Save(&orderEntity)
	orderApi.orderItemRepository.SaveAll(orderEntity.Items)

	responses.StatusCreated(responseWriter, request, nil)
}

func (orderApi *API) getOrder(responseWriter http.ResponseWriter, request *http.Request) {
	params := httprouter.ParamsFromContext(request.Context())
	orderId := orders.OrderId(params.ByName("orderId"))
	orderEntity, err := orderApi.orderRepository.FindById(orderId)
	if err != nil {
		responses.Error(responseWriter, request, http.StatusNotFound, 300, err.Error())
		return
	}
	orderItemEntities, err := orderApi.orderItemRepository.FindAllByOrderId(orderId)
	if err != nil {
		responses.Error(responseWriter, request, http.StatusInternalServerError, 100, err.Error())
	}
	orderEntity.Items = orderItemEntities

	entity := FromOrderEntity(&orderEntity)
	responses.StatusOK(responseWriter, request, &entity)
}
