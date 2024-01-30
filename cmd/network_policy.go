package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"platform-tools/internal"
	"strconv"
)

const networkPolicyPath = "/networking/v1/external/policies"

type networkPolicy struct {
	appMatch      internal.FoundAppFn
	cf            *client.Client
	organizations map[string]*resource.Organization
	spaces        map[string]*resource.Space
}

type PolicyResult struct {
	TotalPolicies int `json:"total_policies"`
	Policies      []struct {
		Source struct {
			ID string `json:"id"`
		} `json:"source"`
		Destination struct {
			ID       string `json:"id"`
			Protocol string `json:"protocol"`
			Ports    struct {
				Start int `json:"start"`
				End   int `json:"end"`
			} `json:"ports"`
		} `json:"destination"`
	} `json:"policies"`
}

var policyPortNum string
var networkPolicyCmd = &cobra.Command{
	Use:   "network-policy",
	Short: "Operates on network policies",
	RunE: func(cmd *cobra.Command, args []string) error {
		cf, err := internal.Client()
		if err != nil {
			return err
		}
		np := &networkPolicy{
			appMatch:      internal.DisplayApp,
			cf:            cf,
			organizations: make(map[string]*resource.Organization),
			spaces:        make(map[string]*resource.Space),
		}
		targetPort, err := strconv.Atoi(policyPortNum)
		if err != nil {
			return err
		}
		return np.ListApps(cmd.Context(), targetPort)
	},
}

func init() {
	networkPolicyCmd.Flags().StringVar(&policyPortNum, "port", "", "port number")
	rootCmd.AddCommand(networkPolicyCmd)
}

func (p *networkPolicy) ListApps(ctx context.Context, targetPort int) error {

	req, err := http.NewRequestWithContext(ctx, "GET", p.cf.ApiURL(networkPolicyPath), nil)
	if err != nil {
		return fmt.Errorf("error creating GET request for %s: %w", networkPolicyPath, err)
	}
	resp, err := p.cf.ExecuteAuthRequest(req)
	if err != nil {
		return fmt.Errorf("error executing GET request for %s: %w", networkPolicyPath, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var result PolicyResult
	if resp.Body == nil || resp.StatusCode == http.StatusNoContent {
		return nil
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil && err != io.EOF {
		return fmt.Errorf("error decoding response JSON: %w", err)
	}

	fmt.Printf("Found the following apps using a network policy targeting port %d:\n", targetPort)
	for _, policy := range result.Policies {
		if targetPort >= policy.Destination.Ports.Start && targetPort <= policy.Destination.Ports.End {
			sa, aErr := p.cf.Applications.Get(ctx, policy.Source.ID)
			if aErr != nil {
				return aErr
			}
			da, aErr := p.cf.Applications.Get(ctx, policy.Destination.ID)
			if aErr != nil {
				return aErr
			}

			san, aErr := p.formatAppName(ctx, sa)
			if aErr != nil {
				return aErr
			}
			dan, aErr := p.formatAppName(ctx, da)
			if aErr != nil {
				return aErr
			}
			fmt.Printf("  - %s -> %s\n", san, dan)
		}
	}

	return nil
}

func (p *networkPolicy) formatAppName(ctx context.Context, app *resource.App) (string, error) {
	space, ok := p.spaces[app.Relationships.Space.Data.GUID]
	if !ok {
		var err error
		space, err = p.cf.Spaces.Get(ctx, app.Relationships.Space.Data.GUID)
		if err != nil {
			return "", err
		}
		p.spaces[app.Relationships.Space.Data.GUID] = space
	}
	org, ok := p.organizations[space.Relationships.Organization.Data.GUID]
	if !ok {
		var err error
		org, err = p.cf.Organizations.Get(ctx, space.Relationships.Organization.Data.GUID)
		if err != nil {
			return "", err
		}
		p.organizations[space.Relationships.Organization.Data.GUID] = org
	}
	return fmt.Sprintf("%s/%s/%s", org.Name, space.Name, app.Name), nil
}
