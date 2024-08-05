package main

import "testing"

func TestPokeStick(t *testing.T) {
	var env Env

	env.Config.ApiKey = "eONADJwhZbU"
	t.Run("Parse expression", func(t *testing.T) {
		got := resolveExpression("${api_key}", env)
		want :=  env.Config.ApiKey

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
