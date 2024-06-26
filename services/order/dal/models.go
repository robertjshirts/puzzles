package dal

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/puzzles/services/order/gen"
)

type Puzzle struct {
	Id          string         `db:"id"`
	OrderInfoId string         `db:"order_info_id"`
	Name        string         `db:"name"`
	Description string         `db:"description"`
	Price       float64        `db:"price"`
	Type        gen.PuzzleType `db:"puzzle_type"`
}

type ShippingInfo struct {
	Id       string `db:"id"`
	Name     string `db:"name"`
	Address  string `db:"address"`
	AreaCode string `db:"area_code"`
	City     string `db:"city"`
	State    string `db:"state"`
	Country  string `db:"country"`
}

type PaymentInfo struct {
	Id         string `db:"id"`
	CardNumber string `db:"card_number"`
	Expiration string `db:"expiration"`
	Cvv        string `db:"cvv"`
	AreaCode   string `db:"area_code"`
}

type OrderInfo struct {
	Id             string `db:"id"`
	Name           string `db:"name"`
	Status         string `db:"status"`
	ShippingInfoId string `db:"shipping_info_id"`
	PaymentInfoId  string `db:"payment_info_id"`
}

func ToDALModels(order *gen.NewOrderInfo) (*OrderInfo, *ShippingInfo, *PaymentInfo, *[]Puzzle, *gen.Error) {
	if order == nil {
		return nil, nil, nil, nil, &gen.Error{Code: 400, Message: "invalid request"}
	}

	orderId := uuid.New().String()
	paymentId := uuid.New().String()
	shippingId := uuid.New().String()

	var puzzles []Puzzle
	for _, item := range order.Items {
		puzzles = append(puzzles, Puzzle{
			Id:          item.Id,
			OrderInfoId: orderId,
			Name:        item.Name,
			Description: item.Description,
			Price:       item.Price,
			Type:        item.Type,
		})
	}

	return &OrderInfo{
			Id:             orderId,
			Name:           order.Name,
			Status:         "placed",
			ShippingInfoId: shippingId,
			PaymentInfoId:  paymentId,
		}, &ShippingInfo{
			Id:       shippingId,
			Name:     order.ShippingInfo.Name,
			Address:  order.ShippingInfo.Address,
			AreaCode: order.ShippingInfo.AreaCode,
			City:     order.ShippingInfo.City,
			State:    order.ShippingInfo.State,
			Country:  order.ShippingInfo.Country,
		}, &PaymentInfo{
			Id:         paymentId,
			CardNumber: order.PaymentInfo.CardNumber,
			Expiration: order.PaymentInfo.Expiration,
			Cvv:        order.PaymentInfo.Cvv,
			AreaCode:   order.PaymentInfo.AreaCode,
		},
		&puzzles, nil
}

func ToApiModel(order OrderInfo, payment PaymentInfo, shipping ShippingInfo, puzzles []Puzzle) *gen.OrderInfo {
	var items []gen.Puzzle
	for _, item := range puzzles {
		items = append(items, gen.Puzzle{
			Id:          item.Id,
			Name:        item.Name,
			Description: item.Description,
			Price:       item.Price,
			Type:        item.Type,
		})
	}

	fmt.Println("order status " + order.Status)

	return &gen.OrderInfo{
		Id:     order.Id,
		Name:   order.Name,
		Status: gen.OrderStatus(order.Status),
		Items:  items,
		ShippingInfo: gen.ShippingInfo{
			Name:     shipping.Name,
			Address:  shipping.Address,
			AreaCode: shipping.AreaCode,
			City:     shipping.City,
			State:    shipping.State,
			Country:  shipping.Country,
		},
		PaymentInfo: gen.PaymentInfo{
			CardNumber: payment.CardNumber,
			Expiration: payment.Expiration,
			Cvv:        payment.Cvv,
			AreaCode:   payment.AreaCode,
		},
	}
}
