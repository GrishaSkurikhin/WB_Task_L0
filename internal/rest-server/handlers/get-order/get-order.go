package getorder

import (
	"net/http"

	customerrors "github.com/GrishaSkurikhin/WB_Task_L0/internal/custom-errors"
	resp "github.com/GrishaSkurikhin/WB_Task_L0/internal/lib/api/response"
	"github.com/GrishaSkurikhin/WB_Task_L0/internal/lib/logger/sl"
	"github.com/GrishaSkurikhin/WB_Task_L0/internal/models"
	"github.com/GrishaSkurikhin/WB_Task_L0/internal/orders"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
)

type Response struct {
	resp.Response
	Order models.Order `json:"order"`
}

func New(log *slog.Logger, cache orders.CacheGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.getorder.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		orderUUID := r.URL.Query().Get("order_uid")
		order, err := orders.Get(orderUUID, cache)
		if err != nil {
			switch specificErr := err.(type) {
			case customerrors.OrderNotFound:
				log.Error("failed to get order", sl.Err(specificErr))
				render.JSON(w, r, resp.Error("order not found"))
			case customerrors.WrongID:
				log.Error("failed to get order", sl.Err(specificErr))
				render.JSON(w, r, resp.Error("wrong order id"))
			default:
				log.Error("failed to get order", sl.Err(specificErr))
				render.JSON(w, r, resp.Error("internal error"))
			}
			return
		}

		render.JSON(w, r, ResponseOK(order))
		log.Info("order found and submitted")
	}
}

func ResponseOK(order models.Order) Response {
	return Response{
		Response: resp.OK(),
		Order:    order,
	}
}
