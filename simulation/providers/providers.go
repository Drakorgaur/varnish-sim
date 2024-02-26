package providers

import (
	"os"
	"strconv"
	"strings"
)

const (
	// VsimFrmtUrlPosEnvName is the name of the environment variable that
	// holds the position of the URL in the line passed to (default) formatter
	VsimFrmtUrlPosEnvName = "VSIM_FRMT_URL_POS"
	// VsimFrmtSizePosEnvName is the name of the environment variable that
	// holds the position of the Size of request in the line passed to (default) formatter
	VsimFrmtSizePosEnvName = "VSIM_FRMT_SIZE_POS"
)

// providers is a slice of strings that holds the names of the providers
// it is just for CLI to show the available ones.
var providers []string

func init() {
	providers = make([]string, 0)
	providers = append(providers, fileProviderName)
}

// Request is a struct that holds the URL and the Size of the request
// passed for simulation
// Request is struct that are passed into channel connecting provider and simulation
type Request struct {
	Url  string
	Size int
}

// Providers returns the list of available providers
func Providers() []string {
	return providers
}

// Provider is an interface that defines the methods that a provider
// should implement
// Provider is used to generate requests for simulation
type Provider interface {
	// Channel returns a channel that will be used to pass requests
	Channel() <-chan *Request

	// SetFormatter sets the function that will be used to format the line
	// passed to the provider
	// Formatter of line should get a line with a break line at the end and return
	// the URL and the Size of the request
	SetFormatter(func(string) (string, int))

	// String returns the name of the provider
	String() string
}

// defaultFormatter is the default function that formats a line provided
// and extracts from it the URL and the Size of the request
// that will be passed for simulation
// Note: cuts the last character of the line, which is assumed to be a newline character
func defaultFormatter(line string) (string, int) {
	urlPos := 1
	urlPosEnv := os.Getenv(VsimFrmtUrlPosEnvName)
	if urlPosEnv != "" {
		urlPos, _ = strconv.Atoi(urlPosEnv)
	}

	sizePos := 0
	sizePosEnv := os.Getenv(VsimFrmtSizePosEnvName)
	if sizePosEnv != "" {
		if v, err := strconv.Atoi(sizePosEnv); err == nil {
			sizePos = v
		}
	}

	sep := " "
	sepEnv := os.Getenv("VSIM_FRMT_SEP")
	if sepEnv != "" {
		sep = sepEnv
	}

	// parse the line
	split := strings.Split(line[:len(line)-1], sep)
	size, err := strconv.Atoi(split[sizePos])
	if err != nil {
		// default fallback
		return split[urlPos], 1000
	}
	return split[urlPos], size
}

// NewProviderByName returns a new provider that is specified by the name
// arg is passed to the provider to specify the source of the requests
func NewProviderByName(providerName string, arg []string) Provider {
	// use request-provider to generate requests
	// prob via a channel
	switch providerName {
	case "file":
		return &FileProvider{Files: arg}
	default:
		// use default-provider
		// TODO: implement something
	}
	return nil
}
