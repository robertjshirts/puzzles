package handler

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/puzzles/services/order/dal"
	"github.com/puzzles/services/order/gen"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

type OrderHandler struct {
	db    *dal.SQLDal
	ch    *amqp.Channel
	conn  *amqp.Connection
	queue string
}

type QueueMessage struct {
	OrderID string `json:"order_id"`
}

func NewOrderHandler(user string, pass string, host string, port int, db *dal.SQLDal) *OrderHandler {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", user, pass, host, port))
	if err != nil {
		log.Fatal(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	queue, err := ch.QueueDeclare(
		"order_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	return &OrderHandler{
		conn:  conn,
		ch:    ch,
		queue: queue.Name,
		db:    db,
	}
}

func (h *OrderHandler) Close() {
	h.ch.Close()
	h.conn.Close()
}

func (h *OrderHandler) CheckHealth(c *gin.Context) {
	c.Status(200)
}

func (h *OrderHandler) CreateNewOrder(c *gin.Context) {
	var orderJSON gen.NewOrderInfo
	if err := c.BindJSON(&orderJSON); err != nil {
		c.JSON(400, gen.Error{Code: 400, Message: "invalid request"})
		return
	}

	orderInfo, shippingInfo, paymentInfo, puzzles, genErr := dal.ToDALModels(&orderJSON)
	if genErr != nil {
		c.JSON(genErr.Code, genErr)
		return
	}

	genErr = h.db.CreateOrder(*orderInfo, *paymentInfo, *shippingInfo, *puzzles)
	if genErr != nil {
		c.JSON(genErr.Code, genErr)
		return
	}

	body, err := json.Marshal(QueueMessage{OrderID: orderInfo.Id})
	if err != nil {
		c.JSON(500, gen.Error{Code: 500, Message: "There was an issue marshalling the rabbitmq message"})
		return
	}

	err = h.ch.PublishWithContext(c, "", h.queue, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	if err != nil {
		c.JSON(500, gen.Error{Code: 500, Message: "There was an issue publishing the message to the queue"})
		return
	}

	apiOrderModel := dal.ToApiModel(*orderInfo, *paymentInfo, *shippingInfo, *puzzles)
	if apiOrderModel == nil {
		c.JSON(500, gen.Error{Code: 500, Message: "There was an issue converting from DAL models to API Models"})
		return
	}

	c.JSON(201, apiOrderModel)
}

func (h *OrderHandler) DeleteOrder(c *gin.Context, id string) {
	genErr := h.db.DeleteOrder(id)
	if genErr != nil {
		c.JSON(genErr.Code, genErr)
		return
	}
	c.Status(204)
}

func (h *OrderHandler) GetOrder(c *gin.Context, id string) {
	orderInfo, payment, shipping, puzzles, genErr := h.db.GetOrder(id)
	if genErr != nil {
		c.JSON(genErr.Code, genErr)
		return
	}

	apiOrderModel := dal.ToApiModel(*orderInfo, *payment, *shipping, *puzzles)
	c.JSON(200, *apiOrderModel)
}
