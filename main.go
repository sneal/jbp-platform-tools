package main

import (
	"context"
	"flag"
	"fmt"
	"golang.org/x/exp/slices"
	"os"

	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/config"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

func main() {
	name := flag.String("name", "java_buildpack_offline", "name of the java buildpack to find apps using")
	flag.Parse()
	err := execute(*name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func execute(bpName string) error {
	ctx := context.Background()
	conf, err := config.NewFromCFHome(config.SkipTLSValidation())
	if err != nil {
		return err
	}
	cf, err := client.New(conf)
	if err != nil {
		return err
	}

	var appsUsingBuildpack []string
	opts := &client.AppListOptions{
		LifecycleType: resource.LifecycleBuildpack,
	}
	for {
		apps, pager, err := cf.Applications.List(ctx, opts)
		if err != nil {
			return err
		}
		for _, app := range apps {
			if appUsesBuildpack(ctx, app, bpName) {
				n, err := appDisplayName(ctx, cf, app)
				if err != nil {
					return err
				}
				appsUsingBuildpack = append(appsUsingBuildpack, n)
				fmt.Println(n)
			}
		}
		if !pager.HasNextPage() {
			break
		}
		pager.NextPage(opts)
	}

	fmt.Printf("Found %d apps using the buildpack %s", len(appsUsingBuildpack), bpName)
	return nil
}

func appUsesBuildpack(ctx context.Context, app *resource.App, bpName string) bool {
	return slices.IndexFunc(app.Lifecycle.BuildpackData.Buildpacks, func(bp string) bool { return bp == bpName }) > -1
}

func appDisplayName(ctx context.Context, cf *client.Client, app *resource.App) (string, error) {
	space, err := cf.Spaces.Get(ctx, app.Relationships.Space.Data.GUID)
	if err != nil {
		return "", err
	}
	org, err := cf.Organizations.Get(ctx, space.Relationships.Organization.Data.GUID)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s/%s", org.Name, space.Name, app.Name), nil
}
