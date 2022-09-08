package provisioner

import (
	"context"
	"github.com/broswen/vex/internal/project"
	"github.com/broswen/vex/internal/token"
)

type Provisioner interface {
	ProvisionProject(ctx context.Context, p *project.Project) error
	DeprovisionProject(ctx context.Context, p *project.Project) error
	ProvisionToken(ctx context.Context, t *token.Token) error
	DeprovisionToken(ctx context.Context, t *token.Token) error
}
