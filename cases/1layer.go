//  Copyright 2024 Mark Barzali
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0

package cases

import (
	"encoding/json"
	"fmt"
	"varnish_sim/model"
)

// LayerConfig is a helper struct for TwoLayerShardedConfig
// it holds the amount of Varnish proxies and cache size on each layer
type LayerConfig struct {
	Amount    int `json:"amount"`
	CacheSize int `json:"cacheSize"`
}

func (l *LayerConfig) String() string {
	return "LayerConfig"
}

func (l *LayerConfig) Validate() error {
	if l.Amount < 0 {
		return fmt.Errorf("amount must be greater than 0")
	}
	if l.CacheSize < 0 {
		return fmt.Errorf("cache size must be greater than 0")
	}
	return nil
}

func (l *LayerConfig) Store() error {
	return store(l)
}

type OneLayer struct {
	// default backend that serves `all` requests
	backend *model.Backend

	// layers of Varnish proxies
	proxies []*model.VarnishProxy

	// configuration for the case
	config LayerConfig
}

func (o *OneLayer) Step() error {
	var err error

	for _, proxy := range o.proxies {
		err = WriteStep(proxy)
	}

	return err
}

func (o *OneLayer) PrintResultsCB(isJson bool) func() error {
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

func (o *OneLayer) SetUp() ([]*model.VarnishProxy, error) {
	o.backend = &model.Backend{Hostname: "default"}

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
	}
	o.proxies = proxies
	return proxies, nil
}

func (o *OneLayer) Validate() error {
	return o.config.Validate()
}

func NewOneLayer(config LayerConfig) *OneLayer {
	return &OneLayer{config: config}
}
