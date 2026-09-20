package engine

import (
	"errors"
	"slices"

	"github.com/jjmrocha/ai-toolkit/skills"
	"github.com/jjmrocha/go-algo/fn"
	"github.com/jjmrocha/warren/internal/config"
)

var coreSkills = []string{
	"buffett-valuation",
	"company-research",
}

func newSkillCollection(cfg *config.Config) (*skills.Collection, error) {
	skillCollection := skills.NewCollection()

	problems := fn.Map(slices.Concat(coreSkills, cfg.Skills), skillCollection.AddClaudeSkill)

	if err := errors.Join(problems...); err != nil {
		return nil, err
	}

	return skillCollection, nil
}
