package dal

import (
	"github.com/google/uuid"

	"github.com/puzzles/services/orders/gen"
)

type Puzzle struct {
	Id          string         `db:"id"`
	OrderInfoId string         `db:"order_info_id"`
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
	Id             string          `db:"id"`
	Name           string          `db:"name"`
	Status         gen.OrderStatus `db:"status"`
	ShippingInfoId string          `db:"shipping_info_id"`
	PaymentInfoId  string          `db:"payment_info_id"`
}

func ToDALModels(order *gen.NewOrderInfo) (*OrderInfo, *ShippingInfo, *PaymentInfo, *[]Puzzle) {
	if order == nil {
		return nil, nil, nil, nil
	}

	orderId := uuid.New().String()
	paymentId := uuid.New().String()
	shippingId := uuid.New().String()

	var puzzles []Puzzle
	for _, item := range order.Items {
		puzzles = append(puzzles, Puzzle{
			Id:          item.Id,
			OrderInfoId: orderId,
			Description: item.Description,
			Price:       item.Price,
			Type:        item.Type,
		})
	}

	return &OrderInfo{
			Id:             orderId,
			Name:           order.Name,
			Status:         gen.OrderStatus("placed"),
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
		&puzzles
}

func ToApiModel(order OrderInfo, shipping ShippingInfo, payment PaymentInfo, puzzles []Puzzle) *gen.OrderInfo {
	var items []gen.Puzzle
	for _, item := range puzzles {
		items = append(items, gen.Puzzle{
			Id:          item.Id,
			Description: item.Description,
			Price:       item.Price,
			Type:        item.Type,
		})
	}

	return &gen.OrderInfo{
		Id:     order.Id,
		Name:   order.Name,
		Status: order.Status,
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
