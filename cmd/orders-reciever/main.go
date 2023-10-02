package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/GrishaSkurikhin/WB_Task_L0/internal/config"
	"github.com/nats-io/stan.go"
)

const (
	natsClientID = "2"
)

func main() {
	ordersFiles := []string{"order1.json", "order2.json"}
	path := "cmd/orders-reciever/"

	ordersByte := make([][]byte, 0, len(ordersFiles))
	for _, orderFile := range ordersFiles {
		file, err := ioutil.ReadFile(path + orderFile)
		if err != nil {
			log.Fatalf("Error opening file %s: %v", orderFile, err)
		}

		ordersByte = append(ordersByte, file)
	}

	cfg := config.MustLoad()
	sc, err := stan.Connect(cfg.Nats.Cluster, natsClientID)
	if err != nil {
		log.Fatalf("Error connect to nats-streaming: %v", err)
	}
	defer sc.Close()

	for _, order := range ordersByte {
		err = sc.Publish(cfg.Nats.SubjectOrderAdd, order)
		if err != nil {
			log.Fatalf("Error sent message: %v", err)
		}
	}

	fmt.Printf("All messages sent")
}
