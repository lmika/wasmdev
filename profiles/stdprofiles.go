package profiles

import "github.com/lmika/wasmdev/gobuilder"

var StandardProfiles = map[string]Profile {
	"go": {
		Pipeline: []gobuilder.PipelineStep{
			gobuilder.GoBuildPipelineStep{},
		},
	},
}
