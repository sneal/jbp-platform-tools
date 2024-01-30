package internal

import (
	"context"
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type FoundAppFn func(context.Context, *client.Client, *resource.App) error

func DisplayApp(ctx context.Context, cf *client.Client, app *resource.App) error {
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
