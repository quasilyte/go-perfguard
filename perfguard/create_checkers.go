package perfguard

import (
	"github.com/quasilyte/go-perfguard/perfguard/checkers"
	"github.com/quasilyte/go-perfguard/perfguard/lint"

	_ "github.com/quasilyte/go-perfguard/perfguard/checkers/callcheckers" // for init()
)

type targetChecker struct {
	ctx  lint.Context
	impl checkers.PackageChecker
}

func (c *targetChecker) CheckTarget(target *lint.Target) error {
	c.ctx.Target = target

	return c.impl.CheckPackage(&c.ctx, target.Files)
}

func createCheckers(config *Config) []*targetChecker {
	packageCheckers := checkers.Create(func(doc checkers.Doc) bool {
		return true
	})

	targetCheckers := make([]*targetChecker, len(packageCheckers))
	for i := range packageCheckers {
		c := &targetChecker{
			impl: packageCheckers[i],
		}
		c.ctx.SetWarnFunc(config.Warn)
		targetCheckers[i] = c
	}

	return targetCheckers
}
