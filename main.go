package main

import "varnish_sim/cli"

func main() {
	err := cli.Run()
	if err != nil {
		print(err.Error())
	}
}
