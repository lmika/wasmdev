package gobuilder

import (
	"context"
)

type GoBuilder struct {
	TargetWasm string

	Pipeline Pipeline
	Hook BuildHook
}

func (gb *GoBuilder) Build() error {
	ctx := WithPipelineEnv(context.Background(), &Env{
		TargetWasm: gb.TargetWasm,
	})

	return gb.withHooks(ctx, func(ctx context.Context) error {
		for _, pipeline := range gb.Pipeline {
			if err := pipeline.Exec(ctx); err != nil {
				return err
			}
		}
		return nil
	})
}

func (gb *GoBuilder) withHooks(ctx context.Context, buildFn func(ctx context.Context) error) error {
	if gb.Hook != nil {
		gb.Hook.OnBuildTriggered()
	}

	err := buildFn(ctx)

	if err != nil {
		if gb.Hook != nil {
			gb.Hook.OnBuildFailed()
		}
		return err
	}

	if gb.Hook != nil {
		gb.Hook.OnBuildSuccess()
	}
	return nil
}
