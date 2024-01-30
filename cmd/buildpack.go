package cmd

import (
	"context"
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
	"platform-tools/internal"
)

type buildpack struct {
	appMatch internal.FoundAppFn
	cf       *client.Client
}

var bpName string
var buildpackCmd = &cobra.Command{
	Use:   "buildpack",
	Short: "Operates on buildpacks",
	RunE: func(cmd *cobra.Command, args []string) error {
		cf, err := internal.Client()
		if err != nil {
			return err
		}
		bp := &buildpack{
			appMatch: internal.DisplayApp,
			cf:       cf,
		}
		return bp.ListApps(cmd.Context(), bpName)
	},
}

func init() {
	buildpackCmd.Flags().StringVar(&bpName, "name", "", "buildpack name")
	rootCmd.AddCommand(buildpackCmd)
}

func (bp *buildpack) ListApps(ctx context.Context, bpName string) error {
	fmt.Printf("Found the following apps using buildpacks %s:\n", bpName)
	opts := &client.AppListOptions{
		LifecycleType: resource.LifecycleBuildpack,
	}
	for {
		apps, pager, err := bp.cf.Applications.List(ctx, opts)
		if err != nil {
			return err
		}
		for _, app := range apps {
			if appUsesBuildpack(app, bpName) {
				err = bp.appMatch(ctx, bp.cf, app)
				if err != nil {
					return nil
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

func appUsesBuildpack(app *resource.App, bpName string) bool {
	return slices.IndexFunc(app.Lifecycle.BuildpackData.Buildpacks, func(bp string) bool { return bp == bpName }) > -1
}
