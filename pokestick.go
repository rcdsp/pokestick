package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"os"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
	"github.com/pelletier/go-toml"
)

type TomlMap map[string]map[string]any

type Flags struct {
	EnvPath string
	RequestPath string
}

type Req struct {
	Request Request
	Body map[string]string
	Headers map[string]any
	Save map[string]string
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
		"../poke-ps/dev.env.toml",
		"Path to the root of a pokestick project or a specific .env.toml file",
	)

	flag.StringVar(
		&req, 
		"path", 
		"../poke-ps/authenticate.toml",
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
	// SECTION build http request
	// build request reqUrl method and content type
	reqUrl := resolveExpression(config.Request.Url)
	method := config.Request.Method
	contentType := resolveExpression(config.Request.ContentType)

	
	// resolve and build body
	for k, v := range config.Body {
		config.Body[k] = resolveExpression(v)
	}
	payload := createKeyValuePairs(config.Body)

	// create http request
	req, err :=http.NewRequest(string(method), reqUrl, strings.NewReader(payload))
	if err != nil {
		panic(err)
	}

	// asdd headers to request
	req.Header.Add("Content-Type", contentType)
	for k, v := range env["headers"] {
		if config.Headers[k] == false {
			continue
		}

		req.Header.Set(k, resolveExpression(v.(string)))
	}

	// request logging
	if config.Request.Log {
		color.Green(
			fmt.Sprintf("REQ => %s %s\n", config.Request.Method, reqUrl), 
		)
		
		headers := req.Header

		for k:= range headers {
			fmt.Printf("%s: %s\n", k, req.Header.Get(k))
		}

		fmt.Println("payload:")
		for k, v := range config.Body {
			fmt.Printf( "%s: %s\n", k, v)
		}
	}
	

	// SECTION execute request and handle response
	// create http client and do request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	// process response
	var responseBody map[string]any
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	err =  oj.Unmarshal(bytes, &responseBody)
	if err != nil {
		panic(err)
	}

	saveResponseValues(env, config, responseBody) 
	
	if config.Response.Log {
		color.Green(fmt.Sprintf("\nRES => %s\n", res.Status)) 
		jsonString, err := json.MarshalIndent(responseBody, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(jsonString))
	}
}

// resolve inline interpolation expressions in the toml file
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

// small utility function to create a key value pair string from a map
func createKeyValuePairs(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
			fmt.Fprintf(b, "%s=%s&", key,url.QueryEscape(value))
	}
	return b.String()
}

// small utility function to save response values to a map
func saveResponseValues(env TomlMap, config Req, res map[string]any) {
	for k, v := range config.Save {
		exp, err := jp.ParseString(v)
		if err != nil {
			panic(err)
		}
		env["config"][k] = exp.Get(res)[0]
	}
}
