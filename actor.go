package main

import (
	"fmt"
	"os"

	"github.com/data-harvesters/goapify"
)

type input struct {
	*goapify.ProxyConfigurationOptions `json:"proxyConfiguration"`
	*goapify.UrlConfigurations         `json:"startUrls"`

	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

func main() {
	a := goapify.NewActor(
		os.Getenv("APIFY_DEFAULT_KEY_VALUE_STORE_ID"),
		os.Getenv("APIFY_TOKEN"),
		os.Getenv("APIFY_DEFAULT_DATASET_ID"),
	)

	i := new(input)

	err := a.Input(i)
	if err != nil {
		fmt.Printf("failed to decode input: %v\\n", err)
		panic(err)
	}

	if i.ProxyConfigurationOptions != nil {
		err = a.CreateProxyConfiguration(i.ProxyConfigurationOptions)
		if err != nil {
			panic(err)
		}
	}

	scraper, err := newScraper(i, a)
	if err != nil {
		fmt.Printf("failed to create scrapper: %v\n", err)
		panic(err)
	}

	scraper.Run()
}
