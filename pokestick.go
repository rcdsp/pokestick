package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/pelletier/go-toml"
)

type TomlMap map[string]map[string]any

type Flags struct {
	EnvPath string
	RequestPath string
}

type Req struct {
	Request Request
	Body map[string]any
	Headers map[string]any
	Save map[string]any
	Response Response
}

type Request =  struct {
	Url string
	Method Method
	ContentType string
	Log bool
}

type Response = struct {
	Log bool
}

type Method string

const (
	GET Method = "GET"
	POST Method = "POST"
	PUT Method = "PUT"
	PATCH Method = "PATCH"
	DELETE Method = "DELETE"
)

var env TomlMap

func main() {
	Init()
}

func Init() {
	flags := parseFlags()
	env = readTomlFile[TomlMap](flags.EnvPath)
	req := readTomlFile[Req](flags.RequestPath)

	executeRequest(req)
}

func parseFlags() Flags {
	var env string
	var req string

	flag.StringVar(
		&env, 
		"env", 
		"mocks/ps/dev.env.toml",
		"Path to the root of a pokestick project or a specific .env.toml file",
	)

	flag.StringVar(
		&req, 
		"path", 
		"mocks/ps/authenticate.toml",
		"Path to .toml file describing an api request", 
	)
	
	flag.Parse()

	var flags = Flags{
		EnvPath: env,
		RequestPath: req,
	}

	return flags
}

func readTomlFile[T TomlMap | Req](path string) T {
	if path == "" {
		panic("No file path provided")
	}
	
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var tomlMap T
	err = toml.Unmarshal(b, &tomlMap)
	if err != nil {
		panic(err)
	}

	return tomlMap
}

func executeRequest(config Req) {
	request := config.Request
	method := resolveExpression(string(request.Method))

	switch method {
	case "GET":
	case "POST":
		handleRequest(config)
	default:
		panic("Unsupported request method")
	}
}

func handleRequest(config Req) {
	url := resolveExpression(config.Request.Url)
	method := config.Request.Method
	contentType := resolveExpression(config.Request.ContentType)
	body, err := json.Marshal(config.Body)
	if err != nil {
		panic(err)
	}
	client := &http.Client{}

	req, err :=http.NewRequest(string(method), url, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", contentType)

	for k, v := range env["headers"] {
		if config.Headers[k] == false {
			continue
		}

		req.Header.Set(k, resolveExpression(v.(string)))
	}
	
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if config.Request.Log {
		renderTableHeader("Request")
		fmt.Printf("%s => %s\n", config.Request.Method, url)
		headers := req.Header

		for k:= range headers {
			fmt.Println(k, req.Header.Get(k))
		}
		fmt.Println("Body:", body)
	}
	
	if config.Response.Log {
		renderTableHeader("Response")
		fmt.Println(res.Status)
		fmt.Println(io.Reader(res.Body))	
	}
}


func resolveExpression(expression string) string {
	re := regexp.MustCompile(`\$\{([^}]*)\}`)
	match := re.FindStringSubmatch(expression)

	if len(match) == 0 {
		return expression
	}

	replace, key := match[0], match[1]

	value, ok := env["config"][key].(string)
	if !ok {
		fmt.Printf("Key %s not found in config %v", key, env["config"])
	}
	value = strings.Replace(expression, replace, value, -1)
	
	return value
}

func renderTableHeader(tableName string) {
	header := strings.Repeat(("="), 50)
	tableName = strings.Repeat(" ", 25 - (len(tableName) / 2)) + strings.ToTitle(tableName)
	fmt.Printf("\n%s\n%s\n%s\n", header, tableName, header)
}
