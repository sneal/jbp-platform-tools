package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/config"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"os"
	"platform-tools/buildpack"
	"platform-tools/port"
)

func main() {
	err := execute(context.Background())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func execute(ctx context.Context) error {
	bpName := flag.String("buildpack", "", "name of the java buildpack to find apps using")
	portNum := flag.String("port", "", "route port to find apps using")
	flag.Parse()

	conf, err := config.NewFromCFHome(config.SkipTLSValidation())
	if err != nil {
		return err
	}
	cf, err := client.New(conf)
	if err != nil {
		return err
	}

	if *bpName != "" {
		err = executeBuildpack(ctx, cf, *bpName)
		if err != nil {
			return err
		}
	}
	if *portNum != "" {
		err = executePort(ctx, cf, *portNum)
		if err != nil {
			return err
		}
	}
	return nil
}

func executeBuildpack(ctx context.Context, cf *client.Client, bpName string) error {
	fmt.Printf("Apps using buildpack: %s\n", bpName)
	bp := buildpack.New(cf, displayApp)
	return bp.ListApps(ctx, bpName)
}

func executePort(ctx context.Context, cf *client.Client, portNum string) error {
	fmt.Printf("Apps using port: %s\n", portNum)
	p := port.New(cf, displayApp)
	return p.ListApps(ctx, portNum)
}

func displayApp(ctx context.Context, cf *client.Client, app *resource.App) error {
	space, err := cf.Spaces.Get(ctx, app.Relationships.Space.Data.GUID)
	if err != nil {
		return err
	}
	org, err := cf.Organizations.Get(ctx, space.Relationships.Organization.Data.GUID)
	if err != nil {
		return err
	}
	fmt.Printf("  - %s/%s/%s\n", org.Name, space.Name, app.Name)
	return nil
}
