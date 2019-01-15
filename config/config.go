package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/cloudfoundry/sonde-go/events"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Nozzel    *NozzelConfig
	WaveFront *WaveFrontConfig
}

type NozzelConfig struct {
	APIURL                 string `required:"true" envconfig:"api_url"`
	Username               string `required:"true"`
	Password               string `required:"true"`
	FirehoseSubscriptionID string `required:"true" envconfig:"firehose_subscription_id"`
	SkipSSL                bool   `default:"false" envconfig:"skip_ssl"`

	SelectedEvents []events.Envelope_EventType `ignored:"true"`
}

type WaveFrontConfig struct {
	URL           string `envconfig:"URL"`
	Token         string `envconfig:"API_TOKEN"`
	ProxyAddr     string `envconfig:"PROXY_ADDR"`
	ProxyPort     int    `envconfig:"PROXY_PORT"`
	FlushInterval int    `required:"true" envconfig:"FLUSH_INTERVAL"`
	Prefix        string `required:"true" envconfig:"PREFIX"`
}

var defaultEvents = []events.Envelope_EventType{
	events.Envelope_ValueMetric,
	events.Envelope_CounterEvent,
}

func Parse() (*Config, error) {
	nozzelConfig := &NozzelConfig{}
	err := envconfig.Process("nozzle", nozzelConfig)
	if err != nil {
		return nil, err
	}

	selectedEvents, err := parseSelectedEvents()
	if err != nil {
		return nil, err
	}
	nozzelConfig.SelectedEvents = selectedEvents

	wavefrontConfig := &WaveFrontConfig{}
	err = envconfig.Process("wavefront", wavefrontConfig)
	if err != nil {
		return nil, err
	}

	config := &Config{Nozzel: nozzelConfig, WaveFront: wavefrontConfig}
	return config, nil
}

func parseSelectedEvents() ([]events.Envelope_EventType, error) {
	envValue := os.Getenv("NOZZLE_SELECTED_EVENTS")
	if envValue == "" {
		return defaultEvents, nil
	} else {
		selectedEvents := []events.Envelope_EventType{}

		for _, envValueSplit := range strings.Split(envValue, ",") {
			envValueSlitTrimmed := strings.TrimSpace(envValueSplit)
			val, found := events.Envelope_EventType_value[envValueSlitTrimmed]
			if found {
				selectedEvents = append(selectedEvents, events.Envelope_EventType(val))
			} else {
				return nil, errors.New(fmt.Sprintf("[%s] is not a valid event type", envValueSlitTrimmed))
			}
		}
		return selectedEvents, nil
	}
}