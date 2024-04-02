package cases

import (
	"encoding/json"
	"fmt"
	"varnish_sim/model"
)

// LayerConfig is a helper struct for TwoLayerShardedConfig
// it holds the amount of Varnish proxies and cache size on each layer

type OneLayerSharded struct {
	// default backend that serves `all` requests
	backend *model.Backend

	// layers of Varnish proxies
	proxies []*model.VarnishProxy

	// configuration for the case
	config LayerConfig
}

func (o *OneLayerSharded) Step() error {
	var err error

	for _, proxy := range o.proxies {
		err = WriteStep(proxy)
	}

	return err
}

func (o *OneLayerSharded) PrintResultsCB(isJson bool) func() error {
	if isJson {
		return func() error {
			for _, proxy := range o.proxies {
				// Marshal the proxy
				raw, err := json.Marshal(proxy.Export())
				if err != nil {
					return err
				}
				fmt.Println(string(raw))
			}
			raw, err := json.Marshal(o.backend.Export())
			if err != nil {
				return err
			}
			fmt.Println(string(raw))
			return nil
		}
	}
	return func() error {
		for _, proxy := range o.proxies {
			proxy.PrintResult()
		}
		return nil
	}
}

func (o *OneLayerSharded) SetUp() ([]*model.VarnishProxy, error) {
	o.backend = &model.Backend{Hostname: "default"}

	director := model.NewShardDirector()

	proxies := make([]*model.VarnishProxy, 0)
	for i := 0; i < o.config.Amount; i++ {
		proxy, err := model.NewVarnishProxy(
			fmt.Sprintf("proxy-%d", i),
			o.config.CacheSize,
		)
		proxy.SetBackend(o.backend)

		if err != nil {
			return nil, err
		}
		proxies = append(proxies, proxy)

		// add backend to common hash-circle
		director.AddBackend(proxy)

		// set director for the proxy
		proxy.SetDirector(director)
	}
	o.proxies = proxies
	return proxies, nil
}

func (o *OneLayerSharded) Validate() error {
	return o.config.Validate()
}

func NewOneLayerSharded(config LayerConfig) *OneLayerSharded {
	return &OneLayerSharded{config: config}
}
