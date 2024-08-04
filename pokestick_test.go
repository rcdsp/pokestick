package main

import "testing"

func TestPokeStick(t *testing.T) {
	var config Config

	config.Env.ApiKey = "eONADJwhZbU"
	t.Run("Parse expression", func(t *testing.T) {
		got := resolveExpression("${api_key}", config)
		want :=  config.Env.ApiKey

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
