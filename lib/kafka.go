package lib

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

// KafkaDriver produces simple messages to a topic. DSN is comma-separated broker list.
type KafkaDriver struct{
	producer sarama.SyncProducer
}

func (k *KafkaDriver) Open(dsn string) error {
	if dsn == "" { return errors.New("empty dsn") }
	brokers := strings.Split(dsn, ",")
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.RequiredAcks = sarama.WaitForLocal
	p, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil { return err }
	k.producer = p
	return nil
}

func (k *KafkaDriver) Close() error { if k.producer!=nil { return k.producer.Close() }; return nil }

// RunOp expects query in format "topic|message". If message omitted uses timestamp.
func (k *KafkaDriver) RunOp(ctx context.Context, query string) (time.Duration, error) {
	if k.producer==nil { return 0, errors.New("producer not open") }
	parts := strings.SplitN(query, "|", 2)
	if len(parts)==0 || parts[0]=="" { return 0, errors.New("missing topic") }
	topic := parts[0]
	msg := "ping"
	if len(parts)==2 { msg = parts[1] }
	m := &sarama.ProducerMessage{Topic: topic, Value: sarama.StringEncoder(msg)}
	start := time.Now()
	_, _, err := k.producer.SendMessage(m)
	return time.Since(start), err
}

func (k *KafkaDriver) HealthCheck(ctx context.Context) error { if k.producer==nil { return errors.New("producer not open") }; return nil }

func (k *KafkaDriver) Query(ctx context.Context, query string) ([]map[string]interface{}, error) {
	return nil, errors.New("query not supported for kafka driver")
}

func init() { Register("kafka", func() Driver { return &KafkaDriver{} }) }
