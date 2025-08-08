package env

import (
	"fmt"
	"github.com/Gustcat/auth/internal/config"
	"github.com/caarlos0/env/v10"
	"net"
	"os"
)

type initialProducerConfig struct {
	KafkaBrokersCount int `env:"KAFKA_BROKERS_COUNT" envDefault:"1"`
	BrokerAddrs       []string
}

var _ config.KafkaProducerConfig = (*initialProducerConfig)(nil)

func NewKafkaProducerConfig() (*initialProducerConfig, error) {
	kp := &initialProducerConfig{}
	if err := env.Parse(kp); err != nil {
		return nil, fmt.Errorf("loading config for kafka_common producer from env is failed: %w", err)
	}

	return kp, nil
}

func (kpg *initialProducerConfig) Brokers() ([]string, error) {
	for i := 1; i <= kpg.KafkaBrokersCount; i++ {
		host := os.Getenv(fmt.Sprintf("KAFKA_BROKER_%d_HOST", i))
		port := os.Getenv(fmt.Sprintf("KAFKA_BROKER_%d_PORT", i))

		if host == "" || port == "" {
			continue
		}

		broker := net.JoinHostPort(host, port)
		kpg.BrokerAddrs = append(kpg.BrokerAddrs, broker)
	}

	if len(kpg.BrokerAddrs) == 0 {
		return nil, fmt.Errorf("no address for Kafka broker")
	}

	return kpg.BrokerAddrs, nil
}
