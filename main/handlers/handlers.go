package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"openapi-explorer/ai"
	"openapi-explorer/models"
	"openapi-explorer/openapi"
	"openapi-explorer/templates"
	"os"
	"strings"

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
	c.Set("Content-Type", "text/html")
	if !ok {
		return templates.ErrorBanner(errors.New("invalid OpenAPI version provided (should be 2 or 3)")).Render(c.Context(), c.Response().BodyWriter())
	}
	if err != nil {
		return templates.ErrorBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}
	src, err := fl.Open()
	if err != nil {
		return templates.ErrorBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}
	bts, err := io.ReadAll(src)
	if err != nil {
		return templates.ErrorBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}
	openApiSchema, err := openapi.OpenAPISpecToString(bts, *typedV)
	if err != nil {
		return templates.ErrorBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}
	response, err := ai.StructuredChat[models.ApiRequestPythonCode](
		fmt.Sprintf("Based on this OpenAPI schema (version %s):\n\n```text\n%s\n```\n\n, create a python code (and list the needed dependencies) using `requests` and all other needed packages to send an API request to the API server based on this user request: %s", v, openApiSchema, msg),
		"You are a senior python engineer with great expertise in API requests using the `requests` package. You produce code along with the needed dependencies to run it. Be accurate, but not too verbose in the produced code.",
		"ApiRequestPythonCode",
		"Python code to send an API request to an API server, as well as the needed dependencies to run the code",
	)
	if err != nil {
		return templates.ErrorBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}
	typedResponse, ok := response.(models.ApiRequestPythonCode)
	if !ok {
		return templates.ErrorBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}
	deps := make([]string, len(typedResponse.Dependencies))
	for _, dep := range typedResponse.Dependencies {
		deps = append(deps, dep.Name+dep.VersionConstraint)
	}
	return templates.GeneratedCode(typedResponse.Code, deps).Render(c.Context(), c.Response().BodyWriter())
}

func HomeRoute(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html")
	return templates.Home().Render(c.Context(), c.Response().BodyWriter())
}

func parsePythonDependency(dep string) (name string, constraints string) {
	dep = strings.TrimSpace(dep)

	// Find first occurrence of version constraint characters
	for i, ch := range dep {
		if ch == '=' || ch == '<' || ch == '>' || ch == '!' || ch == '~' {
			return strings.TrimSpace(dep[:i]), strings.TrimSpace(dep[i:])
		}
	}

	// No version constraint found
	return dep, ""
}

func HandleCodeRun(c *fiber.Ctx) error {
	code := c.FormValue("codeToRun")
	deps := c.FormValue("dependencies")
	dependencies := strings.Split(deps, "\n")
	depsToSubmit := make([]models.Dependency, len(dependencies))
	for i, dep := range dependencies {
		name, vers := parsePythonDependency(dep)
		depsToSubmit[i] = models.Dependency{Name: name, VersionConstraint: vers}
	}
	requestBodyRaw := models.ApiRequestPythonCode{Code: code, Dependencies: depsToSubmit}
	apiKey := os.Getenv("SANDBOX_API_KEY")
	apiEndpoint := os.Getenv("SANDBOX_API_ENDPOINT")
	jsonData, err := json.Marshal(requestBodyRaw)

	c.Set("Content-Type", "text/html")
	if err != nil {
		return templates.ErrorBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return templates.ErrorBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return templates.ErrorBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return templates.ErrorBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return templates.ErrorBanner(fmt.Errorf("response has status %d: %s", resp.StatusCode, string(body))).Render(c.Context(), c.Response().BodyWriter())
	}

	var response models.ApiResponseCodeRun

	err = json.Unmarshal(body, &response)

	if err != nil {
		return templates.ErrorBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}

	return templates.CodeRunResult(response.Output, response.Error).Render(c.Context(), c.Response().BodyWriter())
}
