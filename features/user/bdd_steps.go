package user

import (
	"github.com/cucumber/godog"
	"github.com/huynhtruongson/simple-todo/features/common"
)

type suite struct {
	*common.Suite
	// define more here
}

func NewUserSuite(s *common.Suite) *suite {
	return &suite{
		Suite: s,
	}
}

func (s *suite) InitStep(ctx *godog.ScenarioContext) {
	ctx.Step(`^user creates account with "([^"]*)"$`, s.createUser)
	ctx.Step(`^the account is created successfully$`, s.assertCreatedUser)
	ctx.Step(`^user can login with username and password$`, s.login)
}
