// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package cae

import (
	"strings"

	"github.com/Azure/azqr/internal/scanners"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appcontainers/armappcontainers"
)

// GetRules - Returns the rules for the ContainerAppsScanner
func (a *ContainerAppsScanner) GetRules() map[string]scanners.AzureRule {
	return map[string]scanners.AzureRule{
		"cae-001": {
			Id:          "cae-001",
			Category:    scanners.RulesCategoryReliability,
			Subcategory: scanners.RulesSubcategoryReliabilityDiagnosticLogs,
			Description: "ContainerApp should have diagnostic settings enabled",
			Severity:    scanners.SeverityMedium,
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				service := target.(*armappcontainers.ManagedEnvironment)
				_, ok := scanContext.DiagnosticsSettings[strings.ToLower(*service.ID)]
				return !ok, ""
			},
			Url:   "https://learn.microsoft.com/en-us/azure/container-apps/log-options#diagnostic-settings",
			Field: scanners.OverviewFieldDiagnostics,
		},
		"cae-002": {
			Id:          "cae-002",
			Category:    scanners.RulesCategoryReliability,
			Subcategory: scanners.RulesSubcategoryReliabilityAvailabilityZones,
			Description: "ContainerApp should have availability zones enabled",
			Severity:    scanners.SeverityHigh,
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				app := target.(*armappcontainers.ManagedEnvironment)
				zones := *app.Properties.ZoneRedundant
				return !zones, ""
			},
			Url:   "https://learn.microsoft.com/en-us/azure/container-apps/disaster-recovery?tabs=bash#set-up-zone-redundancy-in-your-container-apps-environment",
			Field: scanners.OverviewFieldAZ,
		},
		"cae-003": {
			Id:          "cae-003",
			Category:    scanners.RulesCategoryReliability,
			Subcategory: scanners.RulesSubcategoryReliabilitySLA,
			Description: "ContainerApp should have a SLA",
			Severity:    scanners.SeverityHigh,
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				return false, "99.95%"
			},
			Url:   "https://azure.microsoft.com/en-us/support/legal/sla/container-apps/v1_0/",
			Field: scanners.OverviewFieldSLA,
		},
		"cae-004": {
			Id:          "cae-004",
			Category:    scanners.RulesCategorySecurity,
			Subcategory: scanners.RulesSubcategorySecurityPrivateEndpoint,
			Description: "ContainerApp should have private endpoints enabled",
			Severity:    scanners.SeverityHigh,
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				app := target.(*armappcontainers.ManagedEnvironment)
				pe := app.Properties.VnetConfiguration != nil && *app.Properties.VnetConfiguration.Internal
				return !pe, ""
			},
			Url:   "https://learn.microsoft.com/en-us/azure/container-apps/vnet-custom-internal?tabs=bash&pivots=azure-portal",
			Field: scanners.OverviewFieldPrivate,
		},
		"cae-006": {
			Id:          "cae-006",
			Category:    scanners.RulesCategoryOperationalExcellence,
			Subcategory: scanners.RulesSubcategoryOperationalExcellenceCAF,
			Description: "ContainerApp Name should comply with naming conventions",
			Severity:    scanners.SeverityLow,
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				c := target.(*armappcontainers.ManagedEnvironment)
				caf := strings.HasPrefix(*c.Name, "cae")
				return !caf, ""
			},
			Url:   "https://learn.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/resource-abbreviations",
			Field: scanners.OverviewFieldCAF,
		},
		"cae-007": {
			Id:          "cae-007",
			Category:    scanners.RulesCategoryOperationalExcellence,
			Subcategory: scanners.RulesSubcategoryOperationalExcellenceTags,
			Description: "ContainerApp should have tags",
			Severity:    scanners.SeverityLow,
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				c := target.(*armappcontainers.ManagedEnvironment)
				return len(c.Tags) == 0, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/azure-resource-manager/management/tag-resources?tabs=json",
		},
	}
}
