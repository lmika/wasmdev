package config

import (
	"go.starlark.net/starlark"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func FromString(s string) (*Config, error) {
	c := &Config{}
	if err := c.readFromReader("str", strings.NewReader(s)); err != nil {
		return nil, err
	}
	return c, nil
}

// FromWasmDevFile returns the config from a WASM dev file.  If there's no
// WASM dev file, an empty config is retuned.
func FromWasmDevFile() (*Config, error) {
	c := &Config{}
	if err := c.readFromFile("wasmdev.star"); err != nil {
		return nil, err
	}
	return c, nil
}

type Config struct {
	globals		starlark.StringDict
}


func (c *Config) readFromFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			c.globals = starlark.StringDict{}
			return nil
		} else {
			return err
		}
	}
	defer f.Close()

	return c.readFromReader(filename, f)
}

func (c *Config) readFromReader(name string, r io.Reader) error {
	bts, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	thread := &starlark.Thread{Name: name}
	globals, err := starlark.ExecFile(thread, name, string(bts), nil)
	if err != nil {
		return err
	}

	c.globals = globals
	return nil
}

// GetString returns the value of a string path.
func (c *Config) GetString(name string, def string) string {
	val, hasVal := c.resolveValueOfPath(name)
	if !hasVal {
		return def
	}

	strVal, wasStringable := starlark.AsString(val)
	if !wasStringable {
		return def
	}

	return strVal
}

func (c *Config) GetBool(name string, def bool) bool {
	val, hasVal := c.resolveValueOfPath(name)
	if !hasVal {
		return def
	}

	return bool(val.Truth())
}

// valueOfPath looks up a value based on the path.  Paths a logical and are separated by dots.
func (c *Config) resolveValueOfPath(name string) (starlark.Value, bool) {
	realName := strings.Replace(name, ".", "_", -1)
	val, hasVal := c.globals[realName]
	return val, hasVal
}