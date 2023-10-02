package getordersid

import (
	"net/http"

	resp "github.com/GrishaSkurikhin/WB_Task_L0/internal/lib/api/response"
	"github.com/GrishaSkurikhin/WB_Task_L0/internal/orders"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
)

type Response struct {
	resp.Response
	Ids []string `json:"ids"`
}

func New(log *slog.Logger, cache orders.CacheGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.getordersid.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		ids := orders.GetIDs(cache)
		render.JSON(w, r, ResponseOK(ids))
		log.Info("orders id submitted")
	}
}

func ResponseOK(ids []string) Response {
	return Response{
		Response: resp.OK(),
		Ids:      ids,
	}
}
