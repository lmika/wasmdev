package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig_ShouldReadValues(t *testing.T) {
	conf, err := FromString(`
devserver_listen = ":1234"
devserver_enabled = False
build_target = "somewhere_else.wasm"
	`)

	// Then
	assert.NoError(t, err)
	assert.False(t, conf.GetBool("devserver.enabled", true))
	assert.Equal(t, ":1234", conf.GetString("devserver.listen", ":1234"))
	assert.Equal(t, "somewhere_else.wasm", conf.GetString("build.target", "main.wasm"))
}

func TestConfig_ShouldReadValueWithOverrides(t *testing.T) {
	conf, err := FromString(`
devserver_listen = ":1234"
devserver_enabled = False
build_target = "somewhere_else.wasm"
	`)

	// Then
	assert.NoError(t, err)
	assert.False(t, conf.GetBool("devserver.enabled", true))
	assert.Equal(t, ":9999", conf.GetString("devserver.listen", ":8080", WithStringOverride(":9999")))
	assert.Equal(t, "somewhere_else.wasm", conf.GetString("build.target", "main.wasm", WithStringOverride("")))
}