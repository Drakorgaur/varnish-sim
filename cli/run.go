package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"varnish_sim/cases"
	"varnish_sim/simulation"
)

func init() {
	fillCasesCmd()
}

// fillCasesCmd fills the root command with subcommands for cases
func fillCasesCmd() {
	root.AddCommand(TwoLayerShardedCmd())
	root.AddCommand(OneLayerCmd())
}

// TwoLayerShardedCmd returns a command for the two-layer sharded case
func TwoLayerShardedCmd() *cobra.Command {
	firstAmount := 0
	firstCacheSize := 0
	secondAmount := 0
	secondCacheSize := 0

	cmd := &cobra.Command{
		Use:     "2layer-sharded",
		Aliases: []string{"2lsh"},
		Short:   "Two-layer sharded case",
		Long:    "Simulation case with two-layer sharded Varnish proxies",
		Args:    cobra.MinimumNArgs(MinArgCount),
		RunE: func(cmd *cobra.Command, args []string) error {
			providerFlag := root.Flag("provider")
			if providerFlag == nil {
				return fmt.Errorf("provider flag is not set")
			}

			twoLayerSharded := cases.NewTwoLayerSharded(
				*cases.NewTwoLayerShardedConfig(firstAmount, firstCacheSize, secondAmount, secondCacheSize),
			)

			err := twoLayerSharded.Validate()
			if err != nil {
				return err
			}

			frontProxies, err := twoLayerSharded.SetUp()
			if err != nil {
				return err
			}

			jsonFlag := root.Flag("json")
			json := jsonFlag.Value.String() == "true"

			return simulation.Run(
				frontProxies,
				args,
				nil,
				providerFlag.Value.String(),
				twoLayerSharded.PrintResultsCB(json),
			)
		},
	}

	cmd.Flags().IntVarP(&firstAmount, "first-amount", "f", 0, "Amount of Varnish proxies in the first layer")
	cmd.Flags().IntVarP(&firstCacheSize, "first-cache-size", "F", 0, "Cache size of Varnish proxies in the first layer")
	cmd.Flags().IntVarP(&secondAmount, "second-amount", "s", 0, "Amount of Varnish proxies in the second layer")
	cmd.Flags().IntVarP(&secondCacheSize, "second-cache-size", "S", 0, "Cache size of Varnish proxies in the second layer")

	return cmd
}

func OneLayerCmd() *cobra.Command {
	amount := 0
	cacheSize := 0

	cmd := &cobra.Command{
		Use:     "1layer",
		Aliases: []string{"1l"},
		Short:   "One-layer case",
		Long:    "Simulation case with one-layer Varnish proxies",
		Args:    cobra.MinimumNArgs(MinArgCount),
		RunE: func(cmd *cobra.Command, args []string) error {
			providerFlag := root.Flag("provider")
			if providerFlag == nil {
				return fmt.Errorf("provider flag is not set")
			}

			oneLayer := cases.NewOneLayer(
				cases.LayerConfig{
					amount,
					cacheSize,
				},
			)

			err := oneLayer.Validate()
			if err != nil {
				return err
			}

			frontProxies, err := oneLayer.SetUp()
			if err != nil {
				return err
			}

			jsonFlag := root.Flag("json")
			json := jsonFlag.Value.String() == "true"

			return simulation.Run(
				frontProxies,
				args,
				nil,
				providerFlag.Value.String(),
				oneLayer.PrintResultsCB(json),
			)
		},
	}

	cmd.Flags().IntVarP(&amount, "amount", "a", 0, "Amount of Varnish proxies")
	cmd.Flags().IntVarP(&cacheSize, "cache-size", "c", 0, "Cache size of Varnish proxies")

	return cmd
}
