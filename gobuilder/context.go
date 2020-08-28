package gobuilder

import "context"

type Env struct {
	TargetWasm string
}

var pipelineEnvKey = struct{}{}

func WithPipelineEnv(ctx context.Context, env *Env) context.Context {
	return context.WithValue(ctx, pipelineEnvKey, env)
}

func PipelineEnv(ctx context.Context) *Env {
	return ctx.Value(pipelineEnvKey).(*Env)
}
