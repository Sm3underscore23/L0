package main

import (
	"flag"
	"log"
	"os"

	"test_produser/config"
	"test_produser/producer"
	"test_produser/usecase/order"
)

var (
	configPath = flag.String("config-path", "", "Path to config file")
	ordersNum  = flag.Int("order-num", 0, "Number of orders to create")
)

func init() {
	err := os.Chdir("..")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(os.Getwd())
}

func main() {
	flag.Parse()

	mainCfg, err := config.NewConfig(*configPath)
	if err != nil {
		log.Fatalf("loading config error: %s", err)
	}

	producer, err := producer.New(mainCfg.Kafka.BrokerList)
	if err != nil {
		log.Fatalf("kafka service failed: %s", err)
	}

	orderUC := order.New(producer, mainCfg.Kafka.Topic, *ordersNum)

	orderUC.SendOrderData()
}
