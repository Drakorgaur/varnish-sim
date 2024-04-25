//  Copyright 2024 Mark Barzali
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0

package model

import "golang.org/x/exp/constraints"

// Numeric is a constraint for a numeric type
// used in generic types
type Numeric interface {
	constraints.Integer | constraints.Float
}
