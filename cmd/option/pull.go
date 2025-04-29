package option

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/spf13/pflag"
)

// Pull represents the options for the pull command.
type Pull struct {
	Token
	RegistryDomain   string
	registryEndpoint string
}

// SetFlag sets the registry domain and token mode for the pull option.
func (p *Pull) SetFlag(registryDomain string, tokenMode string) {
	p.RegistryDomain = registryDomain
	p.Token.SetFlag(tokenMode)
}

// ApplyFlags applies the flags for the pull command.
func (p *Pull) ApplyFlags(flags *pflag.FlagSet) {
	flags.StringVarP(&p.registryEndpoint, "registry-endpoint", "e", "", "Endpoint of the registry domain (default: registryDomain)")
}

// Parse parses the pull options and returns the appropriate token.
func (p *Pull) Parse() error {
	if err := p.Token.Parse(p.RegistryDomain); err != nil {
		return err
	}
	if p.registryEndpoint != "" {
		http.DefaultClient.Transport = &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				if after, found := strings.CutPrefix(addr, p.RegistryDomain); found {
					// Resolve registry to endpoint
					return net.Dial(network, p.registryEndpoint+after)
				}
				return net.Dial(network, addr)
			},
		}
	}
	return nil
}
