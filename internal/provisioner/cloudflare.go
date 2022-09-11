package provisioner

import (
	"context"
	"github.com/broswen/vex/internal/flag"
	"github.com/broswen/vex/internal/project"
	"github.com/broswen/vex/internal/token"
	"github.com/cloudflare/cloudflare-go"
	"github.com/rs/zerolog/log"
)

type CloudflareProvisioner struct {
	api                  *cloudflare.API
	projectKVNamespaceID string
	tokenKVNamespaceID   string
	projectStore         project.ProjectStore
	flagStore            flag.FlagStore
	tokenStore           token.TokenStore
}

func NewCloudflareProvisioner(apiToken, accountID, projectVNamespaceID, tokenKVNamespaceID string, projectStore project.ProjectStore, flagStore flag.FlagStore, tokenStore token.TokenStore) (*CloudflareProvisioner, error) {
	api, err := cloudflare.NewWithAPIToken(apiToken)
	api.AccountID = accountID
	if err != nil {
		return nil, err
	}
	return &CloudflareProvisioner{
		api:                  api,
		projectKVNamespaceID: projectVNamespaceID,
		tokenKVNamespaceID:   tokenKVNamespaceID,
		projectStore:         projectStore,
		flagStore:            flagStore,
		tokenStore:           tokenStore,
	}, nil
}

func (p *CloudflareProvisioner) ProvisionProject(ctx context.Context, pr *project.Project) error {
	project, err := p.projectStore.Get(ctx, pr.ID)
	if err != nil {
		return err
	}
	flags, err := p.flagStore.List(ctx, pr.ID, 1000, 0)
	if err != nil {
		return err
	}
	rendered, err := flag.RenderConfig(flags)
	if err != nil {
		return err
	}
	resp, err := p.api.WriteWorkersKVBulk(ctx, p.projectKVNamespaceID, cloudflare.WorkersKVBulkWriteRequest{
		{
			Key:      project.ID,
			Value:    string(rendered),
			Metadata: project.AccountID,
		},
	})

	if !resp.Success {
		log.Warn().Msgf("errors: %s", resp.Errors)
		log.Warn().Msgf("messages: %s", resp.Messages)
	}
	return err
}

func (p *CloudflareProvisioner) DeprovisionProject(ctx context.Context, pr *project.Project) error {
	resp, err := p.api.DeleteWorkersKV(ctx, p.projectKVNamespaceID, pr.ID)
	if !resp.Success {
		log.Warn().Msgf("errors: %s", resp.Errors)
		log.Warn().Msgf("messages: %s", resp.Messages)
	}
	return err
}

func (p *CloudflareProvisioner) ProvisionToken(ctx context.Context, t *token.Token) error {
	tok, err := p.tokenStore.Get(ctx, t.ID)
	if err != nil {
		return err
	}
	resp, err := p.api.WriteWorkersKVBulk(ctx, p.tokenKVNamespaceID, cloudflare.WorkersKVBulkWriteRequest{
		{
			Key:   tok.Token,
			Value: tok.AccountID,
		},
	})

	if !resp.Success {
		log.Warn().Msgf("errors: %s", resp.Errors)
		log.Warn().Msgf("messages: %s", resp.Messages)
	}
	return err
}

func (p *CloudflareProvisioner) DeprovisionToken(ctx context.Context, t *token.Token) error {
	resp, err := p.api.DeleteWorkersKV(ctx, p.tokenKVNamespaceID, t.Token)
	if !resp.Success {
		log.Warn().Msgf("errors: %s", resp.Errors)
		log.Warn().Msgf("messages: %s", resp.Messages)
	}
	return err
}
