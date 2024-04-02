package provider

import "context"

type AddonProvider interface {
	ProvisionAddon(ctx context.Context, req ProvisionReq) (*ProvisionRes, error)
	DeProvisionAddon(ctx context.Context, req DeProvisionReq) error
	PlanChange(ctx context.Context, req PlanChangeReq) (*PlanChangeRes, error)
	SSO(ctx context.Context, req SSORequest) (*SSOResponse, error)
}
