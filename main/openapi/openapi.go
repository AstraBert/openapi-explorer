package openapi

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/pb33f/libopenapi"
	v2 "github.com/pb33f/libopenapi/datamodel/high/v2"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

type OpenAPIVersion string

const (
	OpenAPIV2 OpenAPIVersion = "2"
	OpenAPIV3 OpenAPIVersion = "3"
)

func processV3Model(model *libopenapi.DocumentModel[v3.Document]) (string, error) {
	paths := model.Model.Paths.PathItems
	stringPaths := ""
	for pair := paths.Oldest(); pair != nil; pair = pair.Next() {
		jsonBytes, err := json.MarshalIndent(pair.Value, "", "  ")
		if err != nil {
			fmt.Printf("Impossible to marshal JSON data for %s because of %s", pair.Key, err.Error())
		}
		s := fmt.Sprintf("Path: %s\nRequest schema: %s\n\n", pair.Key, string(jsonBytes))
		stringPaths += s
	}
	if stringPaths == "" {
		return stringPaths, errors.New("impossible to represent the OpenAPI schema in human readable language")
	}
	return stringPaths, nil
}

func processV2Model(model *libopenapi.DocumentModel[v2.Swagger]) (string, error) {
	paths := model.Model.Paths.PathItems
	stringPaths := ""
	for pair := paths.Oldest(); pair != nil; pair = pair.Next() {
		jsonBytes, err := json.MarshalIndent(pair.Value, "", "  ")
		if err != nil {
			fmt.Printf("Impossible to marshal JSON data for %s because of %s", pair.Key, err.Error())
		}
		s := fmt.Sprintf("Path: %s\nRequest schema: %s\n\n", pair.Key, string(jsonBytes))
		stringPaths += s
	}
	if stringPaths == "" {
		return stringPaths, errors.New("impossible to represent the OpenAPI schema in human readable language")
	}
	return stringPaths, nil
}

func OpenAPISpecToString(file []byte, version OpenAPIVersion) (string, error) {
	// create a new document from specification bytes
	document, err := libopenapi.NewDocument(file)

	if err != nil {
		return "", err
	}

	if version == OpenAPIV3 {
		v3Model, err := document.BuildV3Model()
		if err != nil {
			return "", nil
		}
		return processV3Model(v3Model)
	} else {
		v2Model, err := document.BuildV2Model()
		if err != nil {
			return "", nil
		}
		return processV2Model(v2Model)
	}
}
