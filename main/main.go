package main

import (
	"fmt"
	"openapi-explorer/openapi"
)

func isValidOpenAPIVersion(version string) bool {
	v := openapi.OpenAPIVersion(version)
	return v == openapi.OpenAPIV2 || v == openapi.OpenAPIV3
}

func main() {
	fmt.Println(isValidOpenAPIVersion("2"))
	fmt.Println(isValidOpenAPIVersion("3"))
	fmt.Println(isValidOpenAPIVersion("4"))
}
