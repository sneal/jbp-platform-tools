package cmd

import (
	"context"
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/spf13/cobra"
	"platform-tools/internal"
)

type routing struct {
	appMatch internal.FoundAppFn
	cf       *client.Client
}

var portNum string
var routeCmd = &cobra.Command{
	Use:   "routing",
	Short: "Operates on routes",
	RunE: func(cmd *cobra.Command, args []string) error {
		cf, err := internal.Client()
		if err != nil {
			return err
		}
		p := &routing{
			appMatch: internal.DisplayApp,
			cf:       cf,
		}
		return p.ListApps(cmd.Context(), portNum)
	},
}

func init() {
	routeCmd.Flags().StringVar(&portNum, "port", "", "port number")
	rootCmd.AddCommand(routeCmd)
}

func (p *routing) ListApps(ctx context.Context, port string) error {
	fmt.Printf("Found the following apps using a route with port %s:\n", port)
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
