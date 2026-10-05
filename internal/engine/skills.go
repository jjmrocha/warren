package engine

import (
	"errors"
	"slices"

	"github.com/jjmrocha/ai-toolkit/skills"
	"github.com/jjmrocha/go-algo/fn"
	"github.com/jjmrocha/warren/internal/config"
)

type coreSkill struct {
	name string
	help string
}

var coreSkills = []coreSkill{
	{name: "buffett-valuation", help: "Value a business as a whole"},
	{name: "company-research", help: "Research a company before forming a view"},
}

func coreSkillNames() []string {
	return fn.Map(coreSkills, func(skill coreSkill) string { return skill.name })
}

func newSkillCollection(cfg *config.Config) (*skills.Collection, error) {
	skillCollection := skills.NewCollection()

	problems := fn.Map(slices.Concat(coreSkillNames(), cfg.Skills), skillCollection.AddClaudeSkill)

	if err := errors.Join(problems...); err != nil {
		return nil, err
	}

	return skillCollection, nil
}
