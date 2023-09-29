package main

import (
	"fmt"
	"time"

	"github.com/GrishaSkurikhin/WB_Task_L0/internal/config"
	"github.com/GrishaSkurikhin/WB_Task_L0/internal/models"
	"github.com/GrishaSkurikhin/WB_Task_L0/internal/storage/postgresql"
	"github.com/google/uuid"
)

func main() {
	order1UUID := uuid.New()
	order1 := models.Order{
		OrderUID:    order1UUID,
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: models.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: models.Payment{
			Transaction:  order1UUID,
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDT:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []models.Item{
			{
				ChrtID:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				Rid:         uuid.New(),
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmID:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		SmID:              99,
		DateCreated:       time.Now(),
	}

	order2UUID := uuid.New()
	order2 := models.Order{
		OrderUID:    order2UUID,
		TrackNumber: "WBILM12345678",
		Entry:       "WBIL",
		Delivery: models.Delivery{
			Name:    "John Doe",
			Phone:   "+1234567890",
			Zip:     "12345",
			City:    "Sampleville",
			Address: "123 Main Street",
			Region:  "Sample Region",
			Email:   "johndoe@example.com",
		},
		Payment: models.Payment{
			Transaction:  order2UUID,
			RequestID:    "987654321",
			Currency:     "EUR",
			Provider:     "paypal",
			Amount:       2500,
			PaymentDT:    1637907727,
			Bank:         "samplebank",
			DeliveryCost: 2000,
			GoodsTotal:   500,
			CustomFee:    100,
		},
		Items: []models.Item{
			{
				ChrtID:      123456,
				TrackNumber: "WBILM12345678",
				Price:       250,
				Rid:         uuid.New(),
				Name:        "Widget",
				Sale:        10,
				Size:        "M",
				TotalPrice:  225,
				NmID:        987654,
				Brand:       "Sample Brand",
				Status:      201,
			},
			{
				ChrtID:      789012,
				TrackNumber: "WBILM12345678",
				Price:       150,
				Rid:         uuid.New(),
				Name:        "Sample Product",
				Sale:        20,
				Size:        "L",
				TotalPrice:  120,
				NmID:        543210,
				Brand:       "Another Brand",
				Status:      200,
			},
		},
		Locale:            "en",
		InternalSignature: "abcdef123456",
		CustomerID:        "johndoe123",
		DeliveryService:   "sampleexpress",
		SmID:              42,
		DateCreated:       time.Now(),
	}

	cfg := config.MustLoad()
	postgres, err := postgresql.New(cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.User, cfg.Storage.Password, cfg.Storage.DBName)
	if err != nil {
		panic(err)
	}
	err = postgres.AddOrder(order1)
	if err != nil {
		panic(err)
	}
	err = postgres.AddOrder(order2)
	if err != nil {
		panic(err)
	}

	orders, err := postgres.GetAllOrders()
	if err != nil {
		panic(err)
	}
	fmt.Println(orders)
}
