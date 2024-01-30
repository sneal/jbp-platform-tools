package internal

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/config"
)

func Client() (*client.Client, error) {
	conf, err := config.NewFromCFHome(config.SkipTLSValidation())
	if err != nil {
		return nil, err
	}
	cf, err := client.New(conf)
	if err != nil {
		return nil, err
	}
	return cf, nil
}
