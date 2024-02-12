package simulation

import (
	"fmt"
	"varnish_sim/model"
	"varnish_sim/simulation/providers"
)

// Run starts the simulation
// arg is an argument for provider. For example, a path to a file (for file-provider)
func Run(
	proxies []*model.VarnishProxy,
	args []string,
	formatter func(string) (string, int),
	providerName string,
	simulationEndCb func() error,
) error {
	// use directors to distribute requests
	director := model.NewRoundRobinDirector()
	for _, proxy := range proxies {
		director.AddBackend(proxy)
	}

	//
	provider := providers.NewProviderByName(providerName, args)
	if provider == nil {
		return fmt.Errorf("provider %s not found", providerName)
	}
	provider.SetFormatter(formatter)

	ch := provider.Channel()

	// start the simulation
	for req := range ch {
		if req == nil {
			break
		}
		b := director.GetBackend(req.Url)
		b.Get(req.Url, req.Size)
	}

	return simulationEndCb()
}
