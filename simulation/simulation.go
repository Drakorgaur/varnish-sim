package simulation

import (
	"fmt"
	"os"
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
	stepInterval int,
	stepRegister func() error,
) error {
	// create dir steps if it does not exist
	if err := os.Mkdir("steps", 0755); err != nil && !os.IsExist(err) {
		return err
	}

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
	cnt := 0
	for req := range ch {
		if req == nil {
			break
		}
		b := director.GetBackend(req.Url)
		b.Get(req.Url, req.Size)
		cnt++

		if cnt%stepInterval == 0 {
			// print statistics
			if err := stepRegister(); err != nil {
				fmt.Println(err)
			}
		}
	}

	return simulationEndCb()
}
