package notify

import (
	go_notify "github.com/dzonint/go-notify/go"
	kafka_notify "github.com/dzonint/go-notify/kafka"
	"go.uber.org/zap"
)

type Notifier interface {
	Post(event string, data interface{}) error
	Start(event string, outputChan chan interface{})
	Stop(event string, outputChan chan interface{}) error
	Close()
}

func NewGoNotifier(l *zap.Logger) Notifier {
	return go_notify.NewGoNotify(l)
}

func NewKafkaNotifier(address, topic, groupID string, l *zap.Logger) Notifier {
	return kafka_notify.NewKafkaNotify(address, topic, groupID, l)
}
