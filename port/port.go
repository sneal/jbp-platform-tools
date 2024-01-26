package port

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"jbp-platform-tools/internal"
)

type Port struct {
	appMatch internal.FoundAppFn
	cf       *client.Client
}

func New(cf *client.Client, appMatchCallback internal.FoundAppFn) *Port {
	return &Port{
		appMatch: appMatchCallback,
		cf:       cf,
	}
}

func (p *Port) ListApps(ctx context.Context, port string) error {
	opts := &client.RouteListOptions{}
	opts.Ports.EqualTo(port)
	for {
		routes, pager, err := p.cf.Routes.List(ctx, opts)
		if err != nil {
			return err
		}
		for _, route := range routes {
			for _, d := range route.Destinations {
				if d.App.GUID != nil {
					app, appErr := p.cf.Applications.Get(ctx, *d.App.GUID)
					if appErr != nil {
						return err
					}
					matchErr := p.appMatch(ctx, p.cf, app)
					if matchErr != nil {
						return matchErr
					}
				}
			}
		}
		if !pager.HasNextPage() {
			break
		}
		pager.NextPage(opts)
	}
	return nil
}
