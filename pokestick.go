package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"
	"regexp"

	"internal/strcase"

	"github.com/pelletier/go-toml"
)

type Env struct {
	Config struct {
		Name  string `toml:"name"`
		BaseUrl string `toml:"base_url"`
		ApiKey string `toml:"api_key"`

	}
	Headers struct {
		PsApiKey string `toml:"X-Ps-Api-Key"`
		PsAuthToken string `toml:"X-Ps-Auth-Token"`
	}
}

func main() {
	Init()
}

func Init() {
	var env string
	var req string

	flag.StringVar(
		&env, 
		"env", 
		"mocks/ps/",
		"Path to the root of a gurl project or a specific .env.toml file",
	)

	flag.StringVar(
		&req, 
		"path", 
		"mocks/ps/",
		"Path to .toml file describing an api request", 
	)
	
	flag.Parse()

	file, err := os.Open(env + "8140.group.toml")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var environment Env
	err = toml.Unmarshal(b, &environment)
	if err != nil {
		panic(err)
	}

	resolveEnvironment(reflect.ValueOf(environment), environment)

	fmt.Println(environment)
}

func resolveEnvironment(values reflect.Value, env Env) {
	for i := 0; i < values.NumField(); i++ {
		fieldName := values.Type().Field(i).Name
		fieldValue := values.Field(i)
		
		if fieldValue.Kind() == reflect.Struct {
			fmt.Println(fieldName)
			resolveEnvironment(fieldValue, env)
		}

		if fieldValue.Kind() == reflect.String {
			fieldValue = reflect.ValueOf(resolveExpression(fieldValue.String(), env))
			fmt.Println(values.Type().Field(i).Name, fieldValue)
		}
	}
}

func resolveExpression(expression string, env Env) string {
	re := regexp.MustCompile(`\$\{([^}]*)\}`)
	match := re.FindStringSubmatch(expression)


	if len(match) == 0 {
		return expression
	}

	key := strcase.ToPascal(match[1])
	value := reflect.ValueOf(env.Config).FieldByName(key)

	if value.Kind() == reflect.Invalid {
		return expression
	}	

	return value.String()
}
