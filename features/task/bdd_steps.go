package task

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
	ctx.Step(`^user creates task with "([^"]*)"$`, s.createTask)
	ctx.Step(`^the task is created successfully$`, s.assertCreatedTask)
}
