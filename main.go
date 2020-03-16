package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Info struct {
	Title       string `yaml:"title"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
}
type Feature struct {
	FeatureName string              `yaml:"feature"`
	ID          string              `yaml:"id"`
	Description string              `yaml:"description"`
	Endpoints   []map[string]string `yaml:"endpoints"`
}
type Method struct {
	Description string                 `yaml:"description"`
	OperationID string                 `yaml:"operationId"`
	Responses   map[string]interface{} `yaml:"responses"`
	RBACFeature []string               `yaml:"x-rbac-feature"`
}
type Server struct {
	URL string "yaml:`url`"
}
type OpenApiSpec struct {
	Info    Info                         `yaml:"info"`
	Servers []Server                     `yaml:"servers"`
	Paths   map[string]map[string]Method `yaml:"paths"`
}
type Result struct {
	ServiceName string    `yaml:"serviceName"`
	Features    []Feature `yaml:"features"`
}

func main() {
	// Define flags
	inputFilePath := flag.String("input", "./Rest-api.yaml", "Input YAML file path")
	outputFilePath := flag.String("output", "./result.yaml", "Output YAML file path")
	flag.Parse()

	//Read from file
	inputSpec, err := ioutil.ReadFile(*inputFilePath)
	if err != nil {
		panic(err.Error())
	}

	// Define variables
	var openApiSpec OpenApiSpec
	var result Result

	yaml.Unmarshal(inputSpec, &openApiSpec)

	// Parse file
	result.ServiceName = openApiSpec.Info.Title

	for pathKey, path := range openApiSpec.Paths {
		for methodKey, method := range path {
			for _, feature := range method.RBACFeature {

				var endpoints []map[string]string
				endpoint := make(map[string]string)
				endpoint[pathKey] = methodKey
				endpoints = append(endpoints, endpoint)

				result.Features = append(result.Features, Feature{
					FeatureName: strings.ReplaceAll(feature, "_", " "),
					ID:          feature,
					Description: method.Description,
					Endpoints:   endpoints,
				})
			}
		}
	}

	// Create open file
	outputFile, err := os.Create(*outputFilePath)
	if err != nil {
		panic(err.Error())
	}

	// Write data into outputFile
	yaml.NewEncoder(outputFile).Encode(result)

	fmt.Println(result)
	// Close file
	err = outputFile.Close()
	if err != nil {
		panic(err.Error())
	}
}
