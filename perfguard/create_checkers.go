package perfguard

import (
	"github.com/quasilyte/go-perfguard/perfguard/checkers"
	"github.com/quasilyte/go-perfguard/perfguard/lint"

	_ "github.com/quasilyte/go-perfguard/perfguard/checkers/callcheckers" // for init()
	_ "github.com/quasilyte/go-perfguard/perfguard/checkers/funccheckers" // for init()
)

type targetChecker struct {
	ctx  lint.SharedContext
	impl checkers.PackageChecker
}

func (c *targetChecker) CheckTarget(target *lint.Target) error {
	c.ctx.Reset(target)
	return c.impl.CheckPackage(&c.ctx, target.Files)
}

func createCheckers(config *Config) []*targetChecker {
	packageCheckers := checkers.Create(func(doc checkers.Doc) bool {
		if doc.NeedsProfile && config.Heatmap == nil {
			return false
		}
		return true
	})

	targetCheckers := make([]*targetChecker, len(packageCheckers))
	for i := range packageCheckers {
		c := &targetChecker{
			impl: packageCheckers[i],
		}
		c.ctx.Heatmap = config.Heatmap
		c.ctx.Warn = config.Warn
		targetCheckers[i] = c
	}

	return targetCheckers
}
