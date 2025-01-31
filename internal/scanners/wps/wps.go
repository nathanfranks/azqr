// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package wps

import (
	"github.com/rs/zerolog/log"

	"github.com/Azure/azqr/internal/scanners"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/webpubsub/armwebpubsub"
)

// WebPubSubScanner - Scanner for WebPubSub
type WebPubSubScanner struct {
	config *scanners.ScannerConfig
	client *armwebpubsub.Client
}

// Init - Initializes the WebPubSubScanner
func (c *WebPubSubScanner) Init(config *scanners.ScannerConfig) error {
	c.config = config
	var err error
	c.client, err = armwebpubsub.NewClient(config.SubscriptionID, config.Cred, config.ClientOptions)
	return err
}

// Scan - Scans all WebPubSub in a Resource Group
func (c *WebPubSubScanner) Scan(resourceGroupName string, scanContext *scanners.ScanContext) ([]scanners.AzureServiceResult, error) {
	log.Info().Msgf("Scanning WebPubSub in Resource Group %s", resourceGroupName)

	WebPubSub, err := c.listWebPubSub(resourceGroupName)
	if err != nil {
		return nil, err
	}
	engine := scanners.RuleEngine{}
	rules := c.GetRules()
	results := []scanners.AzureServiceResult{}

	for _, w := range WebPubSub {
		rr := engine.EvaluateRules(rules, w, scanContext)

		results = append(results, scanners.AzureServiceResult{
			SubscriptionID: c.config.SubscriptionID,
			ResourceGroup:  resourceGroupName,
			ServiceName:    *w.Name,
			Type:           *w.Type,
			Location:       *w.Location,
			Rules:          rr,
		})
	}
	return results, nil
}

func (c *WebPubSubScanner) listWebPubSub(resourceGroupName string) ([]*armwebpubsub.ResourceInfo, error) {
	pager := c.client.NewListByResourceGroupPager(resourceGroupName, nil)

	WebPubSubs := make([]*armwebpubsub.ResourceInfo, 0)
	for pager.More() {
		resp, err := pager.NextPage(c.config.Ctx)
		if err != nil {
			return nil, err
		}
		WebPubSubs = append(WebPubSubs, resp.Value...)
	}
	return WebPubSubs, nil
}
