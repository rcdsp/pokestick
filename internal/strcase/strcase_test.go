package strcase

import "testing"

var normalizedCase = "hello world"

func TestStandardizeCase(t *testing.T) {
	t.Run("Standardizes snake case", func(t *testing.T) {
		got := NormalizeCase("hello_world")
		want := normalizedCase

		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})
	t.Run("Standardizes camel case", func (t *testing.T) {
		got := NormalizeCase("helloWorld")
		want := normalizedCase

		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})
	t.Run("Standardizes pascal case", func (t *testing.T) {
		got := NormalizeCase("HelloWorld")
		want := normalizedCase

		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})
	t.Run("Standardizes kebab case", func (t *testing.T) {
		got := NormalizeCase("hello-world")
		want := normalizedCase

		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})
}

func TestToCamel(t *testing.T) {
	got := ToCamel(normalizedCase)
	want := "helloWorld"	

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestToPascal(t *testing.T) {
	got := ToPascal(normalizedCase)
	want := "HelloWorld"	

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestToSnake(t *testing.T) {
	got := ToSnake(normalizedCase)
	want := "hello_world"	

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestToKebab(t *testing.T) {
	got := ToKebab(normalizedCase)
	want := "hello-world"	

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}