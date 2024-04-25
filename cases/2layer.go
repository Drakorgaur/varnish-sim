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

type TwoLayer struct {
	// default backend that serves `all` requests
	backend *model.Backend

	// layers of Varnish proxies
	firstL  []*model.VarnishProxy
	secondL []*model.VarnishProxy

	// configuration for the case
	config TwoLayerShardedConfig
}

func (t *TwoLayer) SetUp() ([]*model.VarnishProxy, error) {
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

	// set respective backends for the first layer
	for i, varnish := range t.firstL {
		varnish.SetBackend(t.secondL[i])
	}

	// return all Varnish proxies that are placed in front.
	// requests will be made to these proxies
	return t.firstL, nil
}

func (t *TwoLayer) Validate() error {
	return t.config.Validate()
}

func (t *TwoLayer) Step() error {
	var err error
	for _, varnish := range t.firstL {
		err = WriteStep(varnish)
	}

	for _, varnish := range t.secondL {
		err = WriteStep(varnish)
	}

	return err
}

func (t *TwoLayer) PrintResultsCB(isJson bool) func() error {
	if isJson {
		return t.PrintResultsJSON
	}
	return t.PrintResultsTable
}

// PrintResultsTable prints the results in a table format
func (t *TwoLayer) PrintResultsTable() error {
	for _, varnish := range t.firstL {
		varnish.PrintResult()
	}

	for _, varnish := range t.secondL {
		varnish.PrintResult()
	}

	return nil
}

// PrintResultsJSON prints the results in a JSON format
func (t *TwoLayer) PrintResultsJSON() error {
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

func NewTwoLayer(config TwoLayerShardedConfig) *TwoLayer {
	return &TwoLayer{
		config: config,
	}
}
