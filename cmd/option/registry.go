package option

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/spf13/pflag"
)

// Registry represents the options related to the registry.
type Registry struct {
	RegistryDomain   string
	registryEndpoint string
}

// SetFlag sets the registry domain.
func (r *Registry) SetFlag(registryDomain string) {
	r.RegistryDomain = registryDomain
}

// ApplyFlags applies the flags to the registry options.
func (r *Registry) ApplyFlags(flags *pflag.FlagSet) {
	flags.StringVarP(&r.registryEndpoint, "registry-endpoint", "e", "", "Endpoint of the registry domain (default: registryDomain)")
}

// Parse parses the registry options and sets up the HTTP client.
func (r *Registry) Parse() error {
	if r.registryEndpoint != "" {
		http.DefaultClient.Transport = &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				if after, found := strings.CutPrefix(addr, r.RegistryDomain); found {
					// Resolve registry to endpoint
					return net.Dial(network, r.registryEndpoint+after)
				}
				return net.Dial(network, addr)
			},
		}
	}
	return nil
}
