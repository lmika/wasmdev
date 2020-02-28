package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig_ShouldReadValues(t *testing.T) {
	conf, err := FromString(`
devserver_listen = ":1234"
devserver_enabled = False
	`)

	// Then
	assert.NoError(t, err)
	assert.False(t, conf.GetBool("devserver.enabled", true))
	assert.Equal(t, ":1234", conf.GetString("devserver.listen", ":1234"))
}
