package model

import "fmt"

// WebInterface is a representation of a web server, web accelerator, or any other web service
type WebInterface interface {
	// Get returns the size of the object in bytes
	Get(string, int) int

	// String returns the name of the web interface
	String() string
}

// Backend is representation of leaf-node in simulation topology
// If you need more complex business logic of Backend request processing,
// create a new representation and implement a WebInterface.
type Backend struct {
	Hostname string
}

// Get interface WebInterface for Backend
// we do not need to use url, as no logic is implemented
func (b *Backend) Get(_ string, size int) int {
	// TODO: make more complex
	return size
}

// String returns the name of the web interface
func (b *Backend) String() string {
	return b.Hostname
}

// VarnishProxy is a representation of a Varnish proxy
type VarnishProxy struct {
	cache    *CacheStorage[string, int] // cache [request URI, object bytes]
	director Director
	hostname string

	backend WebInterface

	// metrics
	cacheMetric   CacheMetric
	routingMetric RoutingMetric[int]
}

func (v *VarnishProxy) TableData() (name string, rows [][]string) {
	name = "VarnishProxy"
	rows = append(rows, []string{"Hostname", v.hostname})
	rows = append(rows, []string{"Cache Size", fmt.Sprintf("%d", v.cache.Size())})
	rows = append(rows, []string{"Cache Used", fmt.Sprintf("%d", v.cache.stored)})
	rows = append(rows, []string{"Routes To", fmt.Sprintf("%s", generateRoutesTo(v))})

	cacheMetric := v.cacheMetric.ExportType()
	rows = append(rows, []string{"Cache hit", fmt.Sprintf("%f", cacheMetric["hit"])})
	rows = append(rows, []string{"Cache miss", fmt.Sprintf("%f", cacheMetric["miss"])})
	rows = append(rows, []string{"CHR", fmt.Sprintf("%f", cacheMetric["hit"]/(cacheMetric["hit"]+cacheMetric["miss"]))})

	for k, v := range v.routingMetric {
		rows = append(rows, []string{fmt.Sprintf("-> %s", k.String()), fmt.Sprintf("%d", v)})
	}

	return
}

func (v *VarnishProxy) Hostname() string {
	return v.hostname
}

func (v *VarnishProxy) StepHeader() string {
	return v.cacheMetric.StepHeader()
}

func (v *VarnishProxy) Step() string {
	return v.cacheMetric.Step()
}

func (v *VarnishProxy) Export() map[string]interface{} {
	self := make(map[string]interface{})
	self["cache"] = v.cacheMetric.ExportType()
	self["routing"] = v.routingMetric.ExportType()
	self["cache_size"] = v.cache.Size()
	self["cache_used"] = v.cache.stored
	self["routes_to"] = generateRoutesTo(v)

	export := make(map[string]interface{})
	export[v.hostname] = self

	return export
}

func generateRoutesTo(v *VarnishProxy) []string {
	routesTo := make([]string, 0)
	if v.director != nil {
		for _, v := range v.director.Backends() {
			routesTo = append(routesTo, fmt.Sprintf("%s", v.String()))
		}
	}
	if v.backend != nil {
		routesTo = append(routesTo, fmt.Sprintf("%s", v.backend.String()))
	}
	return routesTo
}

// NewVarnishProxy
// argument Hostname
// size - storage's size in bytes
func NewVarnishProxy(hostname string, size int) (*VarnishProxy, error) {
	storage, err := NewCacheStorage[string, int](size)
	if err != nil {
		return nil, err
	}

	proxy := VarnishProxy{
		hostname: hostname,
		cache:    storage,
	}

	proxy.initializeMetrics()

	return &proxy, nil
}

func (v *VarnishProxy) initializeMetrics() {
	v.routingMetric = make(map[WebInterface]int)
}

func (v *VarnishProxy) SetDirector(d Director) *VarnishProxy {
	v.director = d
	return v
}

func (v *VarnishProxy) SetBackend(w WebInterface) *VarnishProxy {
	v.backend = w
	return v
}

func (v *VarnishProxy) CacheSize() int {
	return v.cache.Size()
}

// String interface webInterface
func (v *VarnishProxy) String() string {
	return v.hostname
}

// Get interface webInterface
// req - request URI
// size - object size in bytes
func (v *VarnishProxy) Get(req string, size int) int {
	// callback OnRequest
	// try to get from Cache
	obj, ok := v.cache.Get(req)
	if ok {
		v.cacheMetric.Hit()

		return obj
	} else {
		v.cacheMetric.Miss()
	}

	if v.director != nil {
		// director based on its internal logic selects a backend
		backend := v.director.GetBackend(req)
		v.routingMetric[backend]++

		artifactSize := backend.Get(req, size)

		// cache the result
		v.cache.Store(req, artifactSize)

		return artifactSize
	}

	if v.backend != nil {
		artifactSize := v.backend.Get(req, size)

		// cache the result
		v.cache.cache.Add(req, artifactSize)

		return artifactSize
	}

	return 0
}

func (v *VarnishProxy) PrintResult() {
	PrintTable(v)
}
