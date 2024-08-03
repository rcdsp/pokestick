package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"

	"internal/strcase"

	"github.com/pelletier/go-toml"
)

type Config struct {
	Env struct {
		Name  string `toml:"name"`
		BaseUrl string `toml:"base_url"`
		ApiKey string `toml:"api_key"`

		Headers struct {
			PsApiKey string `toml:"X-Ps-Api-Key"`
			PsAuthToken string `toml:"X-Ps-Auth-Token"`
		}
	}
}

func main() {
	Init()
}

func Init() {
	var path string

	flag.StringVar(
		&path, 
		"path", 
		"mocks/ps-gurl/",
		"Path to the root of a gurl project or a specific .toml file",
	)
	
	flag.Parse()

	file, err := os.Open(path + "8140.group.toml")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var config Config
	err = toml.Unmarshal(b, &config)
	if err != nil {
		panic(err)
	}

	fmt.Println(resolveExpression(config.Env.Headers.PsApiKey))
}

func resolveExpression(expression string) string {
	re := regexp.MustCompile(`\$\{([^}]*)\}`)
	match := re.FindStringSubmatch(expression)


	if len(match) == 0 {
		panic(fmt.Errorf("expression %s does not contain a valid key", expression))
	}

	key := match[1]

	return strcase.ToPascal(key)
}
	