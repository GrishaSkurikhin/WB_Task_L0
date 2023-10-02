package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/GrishaSkurikhin/WB_Task_L0/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	ordersView      = "orders_view"
	ordersTable     = "orders"
	deliveriesTable = "deliveries"
	paymentsTable   = "payments"
	itemsTable      = "items"
)

type OrderStorage struct {
	*sqlx.DB
}

func New(host, port, user, password, name string) (*OrderStorage, error) {
	const op = "storage.postgresql.New"

	dataSource := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, name)
	conn, err := sqlx.Connect("postgres", dataSource)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: ping failed: %w", op, err)
	}
	return &OrderStorage{conn}, nil
}

func (s *OrderStorage) Disconnect(ctx context.Context) error {
	const op = "storage.postgresql.Disconnect"

	err := s.Close()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *OrderStorage) GetAllOrders() ([]models.Order, error) {
	const op = "storage.postgresql.GetAllOrders"

	stmt, err := s.Preparex(fmt.Sprintf("SELECT * FROM %s", ordersView))
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	rows := []viewRowStruct{}
	err = stmt.Select(&rows)
	if err != nil {
		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	orders, err := rowsToOrders(rows)
	if err != nil {
		return nil, fmt.Errorf("%s: scan statement: %w", op, err)
	}

	return orders, nil
}

func (s *OrderStorage) AddOrder(order models.Order) error {
	const op = "storage.postgresql.AddOrder"

	tx, err := s.Begin()
	if err != nil {
		return fmt.Errorf("%s: begin transaction: %w", op, err)
	}

	addToDeliveries, err := tx.Prepare(
		fmt.Sprintf(`INSERT INTO %s 
		("name", phone, zip, city, address, region, email) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING delivery_id`, deliveriesTable),
	)
	if err != nil {
		return fmt.Errorf("%s: prepare addToDeliveries statement: %w", op, err)
	}
	defer addToDeliveries.Close()

	addToPayments, err := tx.Prepare(
		fmt.Sprintf(`INSERT INTO %s 
		("transaction", request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`, paymentsTable),
	)
	if err != nil {
		return fmt.Errorf("%s: prepare addToPayments statement: %w", op, err)
	}
	defer addToPayments.Close()

	addToOrders, err := tx.Prepare(
		fmt.Sprintf(`INSERT INTO %s 
		(order_uid, track_number, entry, delivery_id, locale, internal_signature, customer_id, delivery_service, sm_id, date_created) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`, ordersTable),
	)
	if err != nil {
		return fmt.Errorf("%s: prepare addToOrders statement: %w", op, err)
	}
	defer addToOrders.Close()

	addToItems, err := tx.Prepare(
		fmt.Sprintf(`INSERT INTO %s 
		(chrt_id, track_number, price, rid, "name", sale, "size", total_price, nm_id, brand, status) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`, itemsTable),
	)
	if err != nil {
		return fmt.Errorf("%s: prepare addToItems statement: %w", op, err)
	}
	defer addToItems.Close()

	var deliveryID int
	err = addToDeliveries.QueryRow(order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City,
		order.Delivery.Address, order.Delivery.Region, order.Delivery.Email).Scan(&deliveryID)
	if err != nil {
		tx_err := tx.Rollback()
		if tx_err != nil {
			return fmt.Errorf("%s: rollback transaction: %w", op, err)
		}
		return fmt.Errorf("%s: execute addToDeliveries statement: %w", op, err)
	}

	_, err = addToPayments.Exec(order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount,
		time.Unix(int64(order.Payment.PaymentDT), 0), order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		tx_err := tx.Rollback()
		if tx_err != nil {
			return fmt.Errorf("%s: rollback transaction: %w", op, err)
		}
		return fmt.Errorf("%s: execute addToPayments statement: %w", op, err)
	}

	_, err = addToOrders.Exec(order.OrderUID, order.TrackNumber, order.Entry, deliveryID, order.Locale, order.InternalSignature,
		order.CustomerID, order.DeliveryService, order.SmID, order.DateCreated)
	if err != nil {
		tx_err := tx.Rollback()
		if tx_err != nil {
			return fmt.Errorf("%s: rollback transaction: %w", op, err)
		}
		return fmt.Errorf("%s: execute addToDeliveries statement: %w", op, err)
	}

	for _, item := range order.Items {
		_, err = addToItems.Exec(item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale,
			item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			tx_err := tx.Rollback()
			if tx_err != nil {
				return fmt.Errorf("%s: rollback transaction: %w", op, err)
			}
			return fmt.Errorf("%s: execute addToItems statement: %w", op, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s: commit transaction: %w", op, err)
	}

	return nil
}

func rowsToOrders(rows []viewRowStruct) ([]models.Order, error) {
	pointersOrders := make([]*models.Order, 0, len(rows))
	orderMap := make(map[uuid.UUID]*models.Order, len(rows))

	for _, row := range rows {
		if _, isExist := orderMap[row.OrderUID]; !isExist {
			order := models.Order{
				OrderUID:    row.OrderUID,
				TrackNumber: row.TrackNumber,
				Entry:       row.Entry,
				Delivery: models.Delivery{
					Name:    row.DeliveryName,
					Phone:   row.DeliveryPhone,
					Zip:     row.DeliveryZip,
					City:    row.DeliveryCity,
					Address: row.DeliveryAddress,
					Region:  row.DeliveryRegion,
					Email:   row.DeliveryEmail,
				},
				Payment: models.Payment{
					Transaction:  row.PaymentTransaction,
					RequestID:    row.PaymentRequestID,
					Currency:     row.PaymentCurrency,
					Provider:     row.PaymentProvider,
					Amount:       row.PaymentAmount,
					PaymentDT:    int(row.PaymentDT.Unix()),
					Bank:         row.PaymentBank,
					DeliveryCost: row.PaymentDeliveryCost,
					GoodsTotal:   row.PaymentGoodsTotal,
					CustomFee:    row.PaymentCustomFee,
				},
				Items:             []models.Item{},
				Locale:            row.Locale,
				InternalSignature: row.InternalSignature,
				CustomerID:        row.CustomerID,
				DeliveryService:   row.DeliveryService,
				SmID:              row.SmID,
				DateCreated:       row.DateCreated,
			}
			orderMap[order.OrderUID] = &order
			pointersOrders = append(pointersOrders, &order)
		}
		orderMap[row.OrderUID].Items = append(orderMap[row.OrderUID].Items, models.Item{
			ChrtID:      row.ItemChrtID,
			TrackNumber: row.TrackNumber,
			Price:       row.ItemPrice,
			Rid:         row.ItemRid,
			Name:        row.ItemName,
			Sale:        row.ItemSale,
			Size:        row.ItemSize,
			TotalPrice:  row.ItemTotalPrice,
			NmID:        row.ItemNmID,
			Brand:       row.ItemBrand,
			Status:      row.ItemStatus,
		})
	}
	orders := make([]models.Order, len(pointersOrders))
	for i, ptr := range pointersOrders {
		orders[i] = *ptr
	}
	return orders, nil
}

type viewRowStruct struct {
	OrderUID            uuid.UUID `db:"order_uid"`
	TrackNumber         string    `db:"track_number"`
	Entry               string    `db:"entry"`
	DeliveryName        string    `db:"delivery_name"`
	DeliveryPhone       string    `db:"delivery_phone"`
	DeliveryZip         string    `db:"delivery_zip"`
	DeliveryCity        string    `db:"delivery_city"`
	DeliveryAddress     string    `db:"delivery_address"`
	DeliveryRegion      string    `db:"delivery_region"`
	DeliveryEmail       string    `db:"delivery_email"`
	PaymentTransaction  uuid.UUID `db:"payment_transaction"`
	PaymentRequestID    string    `db:"payment_request_id"`
	PaymentCurrency     string    `db:"payment_currency"`
	PaymentProvider     string    `db:"payment_provider"`
	PaymentAmount       int       `db:"payment_amount"`
	PaymentDT           time.Time `db:"payment_dt"`
	PaymentBank         string    `db:"payment_bank"`
	PaymentDeliveryCost int       `db:"payment_delivery_cost"`
	PaymentGoodsTotal   int       `db:"payment_goods_total"`
	PaymentCustomFee    int       `db:"payment_custom_fee"`
	ItemChrtID          int       `db:"item_chrt_id"`
	ItemTrackNumber     string    `db:"item_track_number"`
	ItemPrice           int       `db:"item_price"`
	ItemRid             uuid.UUID `db:"item_rid"`
	ItemName            string    `db:"item_name"`
	ItemSale            int       `db:"item_sale"`
	ItemSize            string    `db:"item_size"`
	ItemTotalPrice      int       `db:"item_total_price"`
	ItemNmID            int       `db:"item_nm_id"`
	ItemBrand           string    `db:"item_brand"`
	ItemStatus          int       `db:"item_status"`
	Locale              string    `db:"locale"`
	InternalSignature   string    `db:"internal_signature"`
	CustomerID          string    `db:"customer_id"`
	DeliveryService     string    `db:"delivery_service"`
	SmID                int       `db:"sm_id"`
	DateCreated         time.Time `db:"date_created"`
}
