package provisioner

import (
	"context"
	"github.com/broswen/vex/internal/flag"
	"github.com/broswen/vex/internal/project"
	"github.com/broswen/vex/internal/token"
	"github.com/cloudflare/cloudflare-go"
	"log"
)

// Ideally this should be replaced with a client that submits change events to a kafka topic
// then a group of workers can process the flag changes in order and scale properly
type CloudflareProvisioner struct {
	api           *cloudflare.API
	kvNamespaceID string
	store         flag.FlagStore
}

func NewCloudflareProvisioner(apiToken, accountID, kvNamespaceID string, store flag.FlagStore) (*CloudflareProvisioner, error) {
	api, err := cloudflare.NewWithAPIToken(apiToken)
	api.AccountID = accountID
	if err != nil {
		return nil, err
	}
	return &CloudflareProvisioner{
		api:           api,
		kvNamespaceID: kvNamespaceID,
		store:         store,
	}, nil
}

func (p *CloudflareProvisioner) ProvisionProject(ctx context.Context, pr *project.Project) error {
	flags, err := p.store.List(ctx, pr.ID, 1000, 0)
	if err != nil {
		return err
	}
	rendered, err := flag.RenderConfig(flags)
	if err != nil {
		return err
	}
	resp, err := p.api.WriteWorkersKV(ctx, p.kvNamespaceID, pr.ID, rendered)

	if !resp.Success {
		log.Println(resp.Messages, resp.Errors)
	}
	return err
}

func (p *CloudflareProvisioner) DeprovisionProject(ctx context.Context, pr *project.Project) error {
	resp, err := p.api.DeleteWorkersKV(ctx, p.kvNamespaceID, pr.ID)
	if !resp.Success {
		log.Println(resp.Messages, resp.Errors)
	}
	return err
}

func (p *CloudflareProvisioner) ProvisionToken(ctx context.Context, t *token.Token) error {
	//TODO implement me
	panic("implement me")
}

func (p *CloudflareProvisioner) DeprovisionToken(ctx context.Context, t *token.Token) error {
	//TODO implement me
	panic("implement me")
}
