package queue

import (
	"encoding/json"
	"hezzltask5/internal/models"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

const bufCap = 10

type LogJournal interface {
	Upload([]models.Event) error
}

type NatsQueue struct {
	conn *nats.Conn
	buf  []models.Event
	LogJournal
}

func NewNatsQueue(c *nats.Conn, log LogJournal) *NatsQueue {
	return &NatsQueue{conn: c, buf: make([]models.Event, 0, 10), LogJournal: log}
}

func NewNatsconn(addr string) (*nats.Conn, error) {
	nc, err := nats.Connect(addr)
	if err != nil {
		return nil, err
	}
	return nc, nil
}

func (nq *NatsQueue) Run() {
	go nq.listen()
}

func (nq *NatsQueue) AddToQueue(m models.Event) {
	data, err := json.Marshal(m)
	if err != nil {
		logrus.Printf("queue level: addtoqueue: %v", err)
	}
	err = nq.conn.Publish("logs", data)
	if err != nil {
		logrus.Printf("queue level: addtoqueue: %v", err)
	}
}

func (nq *NatsQueue) AddTobuf(data []byte) {
	var Event models.Event
	err := json.Unmarshal(data, &Event)
	if err != nil {
		logrus.Printf("queue level: AddTobuf: %v", err)
	}
	if len(nq.buf) == bufCap {
		if err := nq.LogJournal.Upload(nq.buf); err != nil {
			logrus.Printf("queue level: AddTobuf: %v", err)
		}
		nq.buf = nq.buf[:0]
	}
	nq.buf = append(nq.buf, Event)

}

func (nq *NatsQueue) listen() {
	nq.conn.Subscribe("logs", func(m *nats.Msg) {
		nq.AddTobuf(m.Data)
	})
}
