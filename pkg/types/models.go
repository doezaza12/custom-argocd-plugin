package types

import "io"

// Common Models
type ApiRequest struct {
	Url     string
	Method  string
	Headers map[string][]string
	Body    io.Reader
}

type ApiResponse struct {
	Data       []byte
	StatusCode int
}

// GitLab Models
type GroupVariable struct {
	Key              string `json:"key"`
	VariableType     string `json:"variable_type"`
	Value            string `json:"value"`
	Protected        bool   `json:"protected"`
	Masked           bool   `json:"masked"`
	Raw              bool   `json:"raw"`
	EnvironmentScope string `json:"environment_scope"`
	Description      string `json:"description"`
}
