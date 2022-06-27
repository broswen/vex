package provisioner

import (
	"context"
	"github.com/broswen/vex/internal/project"
)

type Provisioner interface {
	ProvisionProject(ctx context.Context, p *project.Project) error
	DeprovisionProject(ctx context.Context, p *project.Project) error
}
