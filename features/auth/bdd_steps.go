package auth

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
	ctx.Step(`^user has already created an account$`, s.createUser)
	ctx.Step(`^user login with "([^"]*)" by "([^"]*)"$`, s.login)
	ctx.Step(`^user login successfully by "([^"]*)"$`, s.assertToken)
}
