package models

type Dependency struct {
	Name              string `json:"name" jsonschema_description:"Name of the dependency"`
	VersionConstraint string `json:"version_constraints" jsonschema_description:"Version constraints for the dependency, e.g. >1, <2, >=0.2.3..."`
}

type ApiRequestPythonCode struct {
	Code         string       `json:"code" jsonschema_description:"Generated python code to run the API request. It should use the 'requests' package to send requests and receive responses."`
	Dependencies []Dependency `json:"dependencies" jsonschema_description:"Dependencies needed for the code to run."`
}

type ApiResponseCodeRun struct {
	Output string `json:"output"`
	Error  string `json:"error"`
}
