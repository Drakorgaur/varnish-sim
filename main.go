//  Copyright 2024 Mark Barzali
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0

package main

import "varnish_sim/cli"

func main() {
	err := cli.Run()
	if err != nil {
		print(err.Error())
	}
}
