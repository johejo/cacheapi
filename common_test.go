package cacheapi

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	const (
		envName = "TEST_ENV"
		envVar  = "VEVTVF9FTlYK"
	)

	if err := os.Setenv(envName, envVar); err != nil {
		t.Errorf("failed to set env: %s", err.Error())
	}

	var (
		expected, actual string
	)

	expected = envVar
	actual = GetEnv(envName, "")
	if actual != expected {
		t.Errorf("failed to get env: envName=%s, envVar=%s, expected=%s, actual=%s", envName, envVar, expected, actual)
	}

	if err := os.Unsetenv(envName); err != nil {
		t.Errorf("failed to unset env: %s", err.Error())
	}

	expected = "DEFAULT"
	actual = GetEnv(envName, "DEFAULT")
	if actual != expected {
		t.Errorf("failed to get env: envName=%s, envVar=%s, expected=%s, actual=%s", envName, envVar, expected, actual)
	}
}
