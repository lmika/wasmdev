package gobuilder

import (
	"log"
	"os"
	"os/exec"
)

type GoBuilder struct {
	TargetWasm string

	Hook BuildHook
}

func (gb *GoBuilder) Build() {
	// TODO: cwd
	cmd := exec.Command("go", "build", "-o", gb.TargetWasm, ".")
	cmd.Env = append(os.Environ(),
		"GOOS=js",
		"GOARCH=wasm",
	)

	if gb.Hook != nil {
		gb.Hook.OnBuildTriggered()
	}

	outErr, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("go build error:\n%v", string(outErr))
		if gb.Hook != nil {
			gb.Hook.OnBuildFailed()
		}
	} else {
		log.Println("rebuilt")
		if gb.Hook != nil {
			gb.Hook.OnBuildSuccess()
		}
	}
}
