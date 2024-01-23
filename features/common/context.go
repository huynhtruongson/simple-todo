package common

import "context"

var (
	StepStateKey = struct{}{}
)

type StepState struct {
	m map[string]interface{}
}

func NewStepState() StepState {
	return StepState{
		m: map[string]interface{}{},
	}
}

func (s StepState) Get(key string) interface{} {
	return s.m[key]
}
func (s StepState) Set(key string, value interface{}) {
	s.m[key] = value
}

func StepStateToContext(ctx context.Context, state StepState) context.Context {
	return context.WithValue(ctx, StepStateKey, state)
}

func StepStateFromContext(ctx context.Context) StepState {
	return ctx.Value(StepStateKey).(StepState)
}
