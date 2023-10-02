package nats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/GrishaSkurikhin/WB_Task_L0/internal/lib/logger/sl"
	"github.com/GrishaSkurikhin/WB_Task_L0/internal/models"
	"github.com/GrishaSkurikhin/WB_Task_L0/internal/orders"
	"github.com/nats-io/stan.go"
	"golang.org/x/exp/slog"
)

type nats struct {
	stan.Conn
}

func New(host, port, stanClusterID, clientID string) (*nats, error) {
	const op = "broker.nats.New"

	sc, err := stan.Connect(stanClusterID, clientID, stan.NatsURL(fmt.Sprintf("%s:%s", host, port)))
	if err != nil {
		return nil, fmt.Errorf("%s: failed to connect: %w", op, err)
	}
	return &nats{sc}, nil
}

func (n *nats) Disconnect(ctx context.Context) error {
	const op = "broker.nats.Disconnect"

	err := n.Close()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (n *nats) SubscribeOrderChannel(log *slog.Logger, cache orders.CacheAdder, storage orders.StorageAdder, subject string) error {
	const op = "broker.nats.SubscribeChannel"

	_, err := n.Subscribe(subject, addFromChannel(log, cache, storage),
		stan.StartWithLastReceived(), stan.SetManualAckMode())

	if err != nil {
		return fmt.Errorf("%s: failed to subscribe: %w", op, err)
	}
	return nil
}

func addFromChannel(log *slog.Logger, cache orders.CacheAdder, storage orders.StorageAdder) stan.MsgHandler {
	return func(msg *stan.Msg) {
		const op = "broker.nats.addFromChannel"

		log = log.With(
			slog.String("op", op),
		)

		var order models.Order
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Error("failed to unmarshal order json", sl.Err(err))
			return
		}

		err = orders.Add(order, cache, storage)
		if err != nil {
			log.Error("failed to add order to storage", sl.Err(err))
			return
		}

		if err := msg.Ack(); err != nil {
			log.Error("failed to acknowledges message", sl.Err(err))
			return
		}
		log.Info("order has been received and saved")
	}
}
