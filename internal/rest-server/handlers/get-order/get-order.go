package getorder

import (
	"net/http"

	resp "github.com/GrishaSkurikhin/WB_Task_L0/internal/lib/api/response"
	"github.com/GrishaSkurikhin/WB_Task_L0/internal/lib/logger/sl"
	"github.com/GrishaSkurikhin/WB_Task_L0/internal/models"
	"github.com/GrishaSkurikhin/WB_Task_L0/internal/orders"
	"github.com/go-chi/render"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
	"github.com/google/uuid"
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

		strOrderUUID := r.URL.Query().Get("order_uid")
		orderUUID, err := uuid.Parse(strOrderUUID)
		if err != nil {
			log.Error("invalid order uid", sl.Err(err))
			render.JSON(w, r, resp.Error("invalid order uid"))
			return
		}

		order, err := orders.Get(orderUUID, cache)
		if err != nil {
			log.Error("failed to get order", sl.Err(err))
			render.JSON(w, r, resp.Error("internal error"))
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
