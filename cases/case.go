package cases

import (
	"encoding/json"
	"fmt"
	"os"
	"varnish_sim/model"
)

// Case is an interface for a simulation case
type Case interface {
	// SetUp initializes the case and returns a list of VarnishProxy instances
	SetUp() ([]*model.VarnishProxy, error)

	// Validate checks if the case is valid
	Validate() error

	Step() error

	PrintResultsCB(bool) func() error
}

// CaseConfig is an interface for a configuration of a simulation case
type CaseConfig interface {
	// String returns the name of the case
	String() string

	// Validate checks if the configuration is valid
	Validate() error

	// Store stores the configuration to a file
	Store() error
}

// store marshals to json the configuration and stores it to a file
func store(c CaseConfig) error {
	// Marshal the config
	raw, err := json.Marshal(c)
	if err != nil {
		return err
	}

	// Store the config
	fd, err := os.Create(c.String() + ".json")
	if err != nil {
		return err
	}

	_, err = fd.Write(raw)
	if err != nil {
		return err
	}

	return nil
}

type StepConfig struct {
	StepInterval int `json:"step_interval"`
}

func WriteStep(v *model.VarnishProxy) error {
	f, err := os.OpenFile(fmt.Sprintf("steps/%s.step", v.Hostname()), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(v.Step() + "\n")
	return err
}
