package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/moov-io/achgateway/internal/service"
	"github.com/moov-io/base/log"
	"gocloud.dev/pubsub"
	"gocloud.dev/pubsub/kafkapubsub"
	"time"
)

func CreateKafkaSubscription(logger log.Logger, cfg *service.KafkaConfig) (*pubsub.Subscription, error) {
	config := kafkapubsub.MinimalConfig()
	config.Version = minKafkaVersion
	config.Net.TLS.Enable = cfg.TLS

	config.Net.SASL.Enable = cfg.Key != ""
	config.Net.SASL.Mechanism = sarama.SASLMechanism("PLAIN")
	config.Net.SASL.User = cfg.Key
	config.Net.SASL.Password = cfg.Secret

	// AutoCommit in Sarama refers to "automated publishing of consumer offsets
	// to the broker" rather than a Kafka broker's meaning of "commit consumer
	// offsets on read" which leads to "at-most-once" delivery.
	config.Consumer.Offsets.AutoCommit.Enable = cfg.AutoCommit

	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.IsolationLevel = sarama.ReadCommitted

	logger.Info().
		Set("tls", log.Bool(cfg.TLS)).
		Set("group", log.String(cfg.Group)).
		Set("sasl.enable", log.Bool(config.Net.SASL.Enable)).
		Set("sasl.user", log.String(cfg.Key)).
		Set("topic", log.String(cfg.Topic)).
		Log("setting up kafka subscription")

	return kafkapubsub.OpenSubscription(cfg.Brokers, config, cfg.Group, []string{cfg.Topic}, &kafkapubsub.SubscriptionOptions{
		WaitForJoin: 10 * time.Second,
	})
}
