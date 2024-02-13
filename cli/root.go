package cli

import (
	"github.com/spf13/cobra"
	"varnish_sim/simulation/providers"
)

// root is the root command for the CLI
var root = &cobra.Command{
	Use:   "vsim",
	Short: "Varnish Cache Simulator",
}

// MinArgCount is minimum amount of arguments required for each the command
// arguments are passed to a provider as source of data for model simulation
const MinArgCount = 1

// setUpRoot sets up the root command
// adds persistent flags to the root command
func setUpRoot() {
	providerFlagUsage := "providers to use.\navailable providers: "
	for _, p := range providers.Providers() {
		providerFlagUsage += p + " "
	}

	// step-interval
	root.PersistentFlags().IntP("step-interval", "", 100, "interval between steps in req count")
	root.PersistentFlags().StringP("provider", "p", "", providerFlagUsage)
	root.PersistentFlags().BoolP("json", "", false, "print json output")
	// not implemented yet
	root.PersistentFlags().StringP("load-balancer", "l", "", "load balancer used to distribute requests for front(edge) proxies")
}

// Run runs the CLI
func Run() error {
	// setUpRoot is called after init, to get filled providers
	setUpRoot()
	return root.Execute()
}
