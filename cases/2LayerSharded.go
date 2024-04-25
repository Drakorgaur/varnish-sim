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

// TwoLayerShardedConfig is a configuration for TwoLayerSharded case
// it holds the amount of Varnish proxies and cache size both layers
type TwoLayerShardedConfig struct {
	FirstLayer  LayerConfig `json:"first_layer"`
	SecondLayer LayerConfig `json:"second_layer"`
}

// NewTwoLayerShardedConfig is a helper constructor for TwoLayerShardedConfig
func NewTwoLayerShardedConfig(amount1, cacheSize1, amount2, cacheSize2 int) *TwoLayerShardedConfig {
	return &TwoLayerShardedConfig{
		FirstLayer: LayerConfig{
			Amount:    amount1,
			CacheSize: cacheSize1,
		},
		SecondLayer: LayerConfig{
			Amount:    amount2,
			CacheSize: cacheSize2,
		},
	}
}

// String returns a string representation of TwoLayerShardedConfig.
// CaseConfig's implementation name
func (c *TwoLayerShardedConfig) String() string {
	return "2layer-sharded"
}

// Store stores the TwoLayerShardedConfig to a file
func (c *TwoLayerShardedConfig) Store() error {
	return store(c)
}

// fillVarnishProxies is a helper function to fill a list with Varnish proxies
// [out] proxyList: a list of Varnish proxies
// [in] prefix: a prefix for the Varnish proxy name
// [in] amount: the amount of Varnish proxies to create
// [in] cacheSize: the cache size for each Varnish proxy
func fillVarnishProxies(
	proxyList *[]*model.VarnishProxy,
	prefix string,
	amount int,
	cacheSize int,
) error {
	for i := 0; i < amount; i++ {
		proxy, err := model.NewVarnishProxy(
			fmt.Sprintf("%s-%d", prefix, i),
			cacheSize,
		)
		if err != nil {
			fmt.Println(err)
			return err
		}
		*proxyList = append(*proxyList, proxy)
	}

	return nil
}

// TwoLayerSharded is a case for a two-layer sharded Varnish setup
// first layer has a director that shards requests to the second layer
//
//	has a backend that serves requests
type TwoLayerSharded struct {
	// default backend that serves `all` requests
	backend *model.Backend

	// layers of Varnish proxies
	firstL  []*model.VarnishProxy
	secondL []*model.VarnishProxy

	// configuration for the case
	config TwoLayerShardedConfig
}

// NewTwoLayerSharded is a constructor for TwoLayerSharded
func NewTwoLayerSharded(config TwoLayerShardedConfig) *TwoLayerSharded {
	return &TwoLayerSharded{config: config}
}

// Validate checks if the configuration is valid
// returns an error if one of the fields is invalid
func (c *TwoLayerShardedConfig) Validate() error {
	if c.FirstLayer.Amount < 1 {
		return fmt.Errorf("first layer amount should be greater than 0")
	}
	if c.FirstLayer.CacheSize < 1 {
		return fmt.Errorf("first layer cache size should be greater than 0")
	}
	if c.SecondLayer.Amount < 1 {
		return fmt.Errorf("second layer amount should be greater than 0")
	}
	if c.SecondLayer.CacheSize < 1 {
		return fmt.Errorf("second layer cache size should be greater than 0")
	}

	return nil
}

// Validate checks if the case is valid, validates its configuration
func (t *TwoLayerSharded) Validate() error {
	return t.config.Validate()
}

// SetUp initializes the case and returns a list of VarnishProxy instances
// returns an error if the case cannot be initialized
func (t *TwoLayerSharded) SetUp() ([]*model.VarnishProxy, error) {
	// initialize default backend that serves `all` requests
	t.backend = &model.Backend{Hostname: "default"}

	// fill layers with Varnish proxies
	err := fillVarnishProxies(&t.secondL, "2", t.config.SecondLayer.Amount, t.config.SecondLayer.CacheSize)
	if err != nil {
		return nil, err
	}

	err = fillVarnishProxies(&t.firstL, "1", t.config.FirstLayer.Amount, t.config.FirstLayer.CacheSize)
	if err != nil {
		return nil, err
	}

	// set default backend for the second layer
	for _, varnish := range t.secondL {
		varnish.SetBackend(t.backend)
	}

	// set director distributing requests to the second layer
	for _, varnish := range t.firstL {
		director := model.NewShardDirector()
		for _, secondLayerVarnish := range t.secondL {
			director.AddBackend(secondLayerVarnish)
		}

		varnish.SetDirector(director)
	}

	// return all Varnish proxies that are placed in front.
	// requests will be made to these proxies
	return t.firstL, nil
}

func (t *TwoLayerSharded) Step() error {
	var err error
	for _, varnish := range t.firstL {
		err = WriteStep(varnish)
	}

	for _, varnish := range t.secondL {
		err = WriteStep(varnish)
	}

	return err
}

// PrintResultsCB returns a callback for printing results
func (t *TwoLayerSharded) PrintResultsCB(isJson bool) func() error {
	if isJson {
		return t.PrintResultsJSON
	}
	return t.PrintResultsTable
}

// PrintResultsTable prints the results in a table format
func (t *TwoLayerSharded) PrintResultsTable() error {
	for _, varnish := range t.firstL {
		varnish.PrintResult()
	}

	for _, varnish := range t.secondL {
		varnish.PrintResult()
	}

	return nil
}

// PrintResultsJSON prints the results in a JSON format
func (t *TwoLayerSharded) PrintResultsJSON() error {
	proxies := make([]map[string]interface{}, 0)

	for _, varnish := range t.firstL {
		proxies = append(proxies, varnish.Export())
	}

	for _, varnish := range t.secondL {
		proxies = append(proxies, varnish.Export())
	}

	// add default backend to the list
	proxies = append(proxies, t.backend.Export())

	raw, err := json.MarshalIndent(proxies, "", " ")
	if err != nil {
		return err
	}

	fmt.Println(string(raw))

	return nil
}
