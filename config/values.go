package config

import "go.starlark.net/starlark"

type ResolvedValue struct {
	value starlark.Value
	isSet bool
}

type ValueOption func(resolvedValue ResolvedValue) ResolvedValue

// IdentityValueOption is a value option the returns the resolved value
func IdentityValueOption(resolvedValue ResolvedValue) ResolvedValue {
	return resolvedValue
}

func WithStringOverride(override string) ValueOption {
	if override == "" {
		return IdentityValueOption
	}

	return func(resolvedValue ResolvedValue) ResolvedValue {
		return ResolvedValue{value: starlark.String(override), isSet: true}
	}
}
