package order

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"test_produser/entity"
	"test_produser/producer"
	"test_produser/usecase"
)

type orderSender struct {
	producer.Producer
	topicName string
	ordersNum int
}

func New(producer producer.Producer, topicName string, ordersNum int) usecase.Order {
	return &orderSender{
		Producer:  producer,
		topicName: topicName,
		ordersNum: ordersNum,
	}
}

func (uc *orderSender) SendOrderData() error {
	for range uc.ordersNum {

		trackNumber := gofakeit.ProductUPC()
		itemsNum := gofakeit.Number(1, 4)
		items := make([]entity.Item, 0, itemsNum)
		var goodsTotal int

		for i := range itemsNum {
			itemPrice := gofakeit.Number(50, 10000)
			sale := gofakeit.Number(1, itemPrice)

			items = append(items, entity.Item{
				ChrtID:      gofakeit.Number(1, 10000000000),
				TrackNumber: trackNumber,
				Price:       itemPrice,
				Rid:         gofakeit.UUID(),
				Name:        gofakeit.ProductName(),
				Sale:        sale,
				Size:        strconv.Itoa(gofakeit.Number(1, 50)),
				TotalPrice:  itemPrice - sale,
				NmID:        gofakeit.Number(1, 10000000000),
				Brand:       gofakeit.Company(),
				Status:      gofakeit.RandomInt([]int{200, 201, 202, 203, 204}),
			})

			goodsTotal += items[i].TotalPrice
		}

		uid := uuid.NewString()
		deliveryCost := gofakeit.Number(100, 5000)
		customFee := gofakeit.Number(0, 1000)

		order := entity.OrderInfo{
			OrderUID:    entity.OrderUID(uid),
			TrackNumber: trackNumber,
			Entry:       gofakeit.Word(),
			Delivery: entity.Delivery{
				Name:    gofakeit.FirstName() + " " + gofakeit.LastName(),
				Phone:   gofakeit.Phone(),
				Zip:     gofakeit.Zip(),
				City:    gofakeit.City(),
				Address: gofakeit.Address().Street,
				Region:  gofakeit.Address().State,
				Email:   gofakeit.Email(),
			},
			Payment: entity.Payment{
				Transaction:  uid,
				RequestID:    "",
				Currency:     gofakeit.CurrencyShort(),
				Provider:     "wbpay",
				Amount:       goodsTotal + deliveryCost + customFee,
				PaymentDT:    int(time.Now().Unix()),
				Bank:         gofakeit.BankName(),
				DeliveryCost: deliveryCost,
				GoodsTotal:   goodsTotal,
				CustomFee:    customFee,
			},
			Items:             items,
			Locale:            gofakeit.LanguageBCP(),
			InternalSignature: "",
			CustomerID:        gofakeit.UUID(),
			DeliveryService:   gofakeit.Word(),
			ShardKey:          strconv.Itoa(gofakeit.Number(1, 10)),
			SmID:              gofakeit.Number(1, 99),
			DateCreated:       time.Now(),
			OofShard:          strconv.Itoa(gofakeit.Number(1, 10)),
		}

		data, err := json.Marshal(order)
		if err != nil {
			return err
		}

		msg := &sarama.ProducerMessage{
			Topic: uc.topicName,
			Value: sarama.ByteEncoder(data),
		}

		if _, _, err = uc.SendMessage(msg); err != nil {
			return err
		}

		log.Println(string(data))
	}

	
	return nil
}
