package model

import "golang.org/x/exp/constraints"

// Numeric is a constraint for a numeric type
// used in generic types
type Numeric interface {
	constraints.Integer | constraints.Float
}
