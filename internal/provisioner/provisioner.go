package provisioner

import (
	"context"
	"github.com/broswen/vex/internal/flag"
	"github.com/broswen/vex/internal/project"
	"github.com/broswen/vex/internal/token"
)

type Provisioner interface {
	ProvisionFlag(ctx context.Context, f *flag.Flag) error
	DeprovisionFlag(ctx context.Context, f *flag.Flag) error
	ProvisionProject(ctx context.Context, p *project.Project) error
	DeprovisionProject(ctx context.Context, p *project.Project) error
	ProvisionToken(ctx context.Context, t *token.Token) error
	DeprovisionToken(ctx context.Context, t *token.Token) error
}
