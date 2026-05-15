package lib

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// SolaceDriver implemented over MQTT (Solace supports MQTT). DSN should be a broker URL
// e.g. tcp://broker:1883 or ssl://broker:8883. Optional query param clientid can set client id.

type SolaceDriver struct {
	client mqtt.Client
}

func (s *SolaceDriver) Open(dsn string) error {
	if dsn == "" { return errors.New("empty dsn") }
	u, err := url.Parse(dsn)
	if err != nil { return err }

	opts := mqtt.NewClientOptions()
	// use raw URL as broker
	opts.AddBroker(u.Scheme + "://" + u.Host + u.Path)

	if u.User != nil {
		if pw, ok := u.User.Password(); ok {
			opts.SetUsername(u.User.Username())
			opts.SetPassword(pw)
		} else {
			opts.SetUsername(u.User.Username())
		}
	}

	q := u.Query()
	clientID := q.Get("clientid")
	if clientID == "" {
		clientID = "db-bench-solace-" + strings.ReplaceAll(time.Now().Format(time.RFC3339Nano), ":", "-")
	}
	opts.SetClientID(clientID)
	// default clean session true
	opts.SetCleanSession(true)

	c := mqtt.NewClient(opts)
	token := c.Connect()
	if !token.WaitTimeout(5 * time.Second) {
		return errors.New("mqtt connect timeout")
	}
	if token.Error() != nil { return token.Error() }
	s.client = c
	return nil
}

func (s *SolaceDriver) Close() error {
	if s.client != nil && s.client.IsConnected() {
		s.client.Disconnect(250)
	}
	return nil
}

// RunOp publishes to topic. query format: "topic|message". If message omitted uses timestamp.
func (s *SolaceDriver) RunOp(ctx context.Context, query string) (time.Duration, error) {
	if s.client == nil { return 0, errors.New("client not open") }
	parts := strings.SplitN(query, "|", 2)
	if len(parts) == 0 || parts[0] == "" { return 0, errors.New("missing topic") }
	topic := parts[0]
	msg := time.Now().Format(time.RFC3339Nano)
	if len(parts) == 2 { msg = parts[1] }
	start := time.Now()
	token := s.client.Publish(topic, 0, false, msg)
	if !token.WaitTimeout(5 * time.Second) {
		return time.Since(start), errors.New("publish timeout")
	}
	return time.Since(start), token.Error()
}

func (s *SolaceDriver) HealthCheck(ctx context.Context) error {
	if s.client == nil { return errors.New("client not open") }
	if s.client.IsConnected() { return nil }
	return errors.New("mqtt client not connected")
}

func (s *SolaceDriver) Query(ctx context.Context, query string) ([]map[string]interface{}, error) {
	return nil, errors.New("query not supported for solace driver")
}

func init() { Register("solace", func() Driver { return &SolaceDriver{} }) }
