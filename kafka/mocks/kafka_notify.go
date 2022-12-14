package mocks

type kafkaNotify struct{}

var Start func(event string, outputChan chan interface{})

func NewKafkaNotify() *kafkaNotify {
	return &kafkaNotify{}
}

func (k *kafkaNotify) Start(event string, outputChan chan interface{}) {}
func (k *kafkaNotify) Close()                                          {}

var Post func(
	event string,
	data interface{},
) error

func (k *kafkaNotify) Post(
	event string,
	data interface{},
) error {
	return Post(event, data)
}

var Stop func(
	event string,
	outputChan chan interface{},
) error

func (k *kafkaNotify) Stop(
	event string,
	outputChan chan interface{},
) error {
	return Stop(event, outputChan)
}

var Listen func(
	outputChan chan interface{},
) error

func (k *kafkaNotify) Listen(
	outputChan chan interface{},
) error {
	return Listen(outputChan)
}
