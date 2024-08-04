package strcase

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// Find string's original case and normalize it: "HelloWorld" || "hello_world" -> "hello world"
func NormalizeCase(s string) string {
	camelCaseRegex := regexp.MustCompile(`^[a-z]+(?:[A-Z][a-z]+)*$`)
	pascalCaseRegex := regexp.MustCompile(`^[A-Z][a-z]+(?:[A-Z][a-z]+)*$`)
	snakeCaseRegex := regexp.MustCompile(`^[a-z]+(?:_[a-z]+)*$`)
	kebabCaseRegex := regexp.MustCompile(`^[a-z]+(?:-[a-z]+)*$`)

	switch {
	case camelCaseRegex.MatchString(s):
	fmt.Println("camel")
		return FromCamel(s)
	case pascalCaseRegex.MatchString(s):
		fmt.Println("pascal")
		return FromPascal(s)
	case snakeCaseRegex.MatchString(s):
		return FromSnake(s)
	case kebabCaseRegex.MatchString(s):
		return FromKebab(s)
	default:
		return s
	}
}

// Normalize from camel case
func FromCamel(s string) string {
	var indexes []int

	for i, r := range s {
		if unicode.IsUpper(r) {
			indexes = append(indexes, i)
		}
	}

	if len(indexes) == 0 {
		return s
	}

	var standardString string

	for i := 0; i < len(indexes); i++ {
		camelSplit := string(s[indexes[i]])
		standardJoint := string(' ') +  strings.ToLower(camelSplit)
		standardString = strings.Replace(s, camelSplit, standardJoint, 1)
	}

	return standardString
}

// Normalize from pascal case
func FromPascal(s string) string {
	s = FromCamel(s)
	s = strings.ToLower(s)
	return s
}

// Normalize from snake case
func FromSnake(s string) string {
	s = strings.ReplaceAll(s, "_", " ")
	return s
}

// Normalize from kebab case
func FromKebab(s string) string {
	s = strings.ReplaceAll(s, "-", " ")
	return s
}

// Normalize any string and convert it to camel case
func ToCamel(s string) string {	
	s = NormalizeCase(s)

	for strings.ContainsAny(s, " ") {
		i := strings.Index(s, " ")
		s = strings.Replace(s, " ", "", 1)
		s = strings.Replace(s, string(s[i]), strings.ToUpper(string(s[i])), 1)
	}

	return s
}

// Normalize any string and convert it to pascal case
func ToPascal(s string) string {
	s = NormalizeCase(s)
	s = ToCamel(s)
	s = strings.ToUpper(s[0:1]) + s[1:]
	return s
}

// Normalize any string and convert it to snake case
func ToSnake(s string) string {
	s = NormalizeCase(s)
	s = strings.ReplaceAll(s, " ", "_")
	return s
}

// Normalize any string and convert it to kebab case
func ToKebab(s string) string {
	s = NormalizeCase(s)
	s = strings.ReplaceAll(s, " ", "-")
	return s
}