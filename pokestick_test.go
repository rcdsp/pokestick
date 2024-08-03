package main

import "testing"

func TestPokeStick(t *testing.T) {
	t.Run("Parse expression", func(t *testing.T) {
		got := resolveExpression("${api_key}")
		want := "ApiKey"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
