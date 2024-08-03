package strcase

import "testing"

func TestToPascal(t *testing.T) {
	got := ToPascal("hello_world")
	want := "HelloWorld"

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}