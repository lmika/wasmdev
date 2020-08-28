package gobuilder

import (
	"context"
	"os"
	"os/exec"
)

type Pipeline []PipelineStep

type PipelineStep interface {
	Exec(ctx context.Context) error
}



type GoBuildPipelineStep struct {
}

func (gb GoBuildPipelineStep) Exec(ctx context.Context) error {
	env := PipelineEnv(ctx)

	cmd := exec.CommandContext(ctx,"go", "build", "-o", env.TargetWasm, ".")
	cmd.Env = append(os.Environ(),
		"GOOS=js",
		"GOARCH=wasm",
	)

	return cmd.Run()
}