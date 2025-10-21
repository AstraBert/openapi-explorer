package handlers

import (
	"fmt"
	"io"
	"openapi-explorer/ai"
	"openapi-explorer/models"
	"openapi-explorer/openapi"

	"github.com/gofiber/fiber/v2"
)

func validateOpenAPIVersion(version string) (*openapi.OpenAPIVersion, bool) {
	v := openapi.OpenAPIVersion(version)
	if v == openapi.OpenAPIV2 || v == openapi.OpenAPIV3 {
		return &v, true
	}
	return nil, false
}

func HandleCodeGeneration(c *fiber.Ctx) error {
	fl, err := c.FormFile("openapiSpec")
	v := c.FormValue("openapiVersion")
	msg := c.FormValue("inputMessage")
	typedV, ok := validateOpenAPIVersion(v)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "invalid OpenAPI version provided (should be 2 or 3)"})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error while parsing the input file: " + err.Error()})
	}
	src, err := fl.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error while opening the input file: " + err.Error()})
	}
	bts, err := io.ReadAll(src)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error while reading the input file: " + err.Error()})
	}
	openApiSchema, err := openapi.OpenAPISpecToString(bts, *typedV)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error while converting the JSON schema to human-readable format: " + err.Error()})
	}
	response := ai.StructuredChat[models.ApiRequestPythonCode](
		fmt.Sprintf("Based on this OpenAPI schema (version %s):\n\n```text\n%s\n```\n\n, create a python code (and list the needed dependencies) using `requests` and all other needed packages to send an API request to the API server based on this user request: %s", v, openApiSchema, msg),
		"You are a senior python engineer with great expertise in API requests using the `requests` package. You produce code along with the needed dependencies to run it. Be accurate, but not too verbose in the produced code.",
		"ApiRequestPythonCode",
		"Python code to send an API request to an API server, as well as the needed dependencies to run the code",
	)
	typedResponse, ok := response.(models.ApiRequestPythonCode)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error while generating the code for the API request"})
	}
	deps := make([]string, len(typedResponse.Dependencies))
	for _, dep := range typedResponse.Dependencies {
		deps = append(deps, dep.Name+dep.VersionConstraint)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"code": typedResponse.Code, "dependencies": deps})
}
