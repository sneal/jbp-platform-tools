package buildpack

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"golang.org/x/exp/slices"
	"platform-tools/internal"
)

type Buildpack struct {
	appMatch internal.FoundAppFn
	cf       *client.Client
}

func New(cf *client.Client, appMatchCallback internal.FoundAppFn) *Buildpack {
	return &Buildpack{
		appMatch: appMatchCallback,
		cf:       cf,
	}
}

func (bp *Buildpack) ListApps(ctx context.Context, bpName string) error {
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
