package gitlab

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/doezaza12/custom-argocd-plugin/pkg/types"
	"github.com/doezaza12/custom-argocd-plugin/pkg/utils"
)

func GetGitlabDefaultHeaders(contentType string) map[string][]string {
	// set default
	if contentType == "" {
		contentType = "application/json"
	}
	headers := make(map[string][]string)
	headers["PRIVATE-TOKEN"] = []string{os.Getenv(types.EnvGitLabPassword)}
	headers["Content-Type"] = []string{contentType}

	return headers
}

func ListGitLabVariables(groupPath string) ([]types.GroupVariable, error) {
	requestUrl := fmt.Sprintf("%s/groups/%s/variables",
		os.Getenv(types.EnvGitLabAPIV4),
		url.PathEscape(groupPath))
	headers := GetGitlabDefaultHeaders("")

	req := types.ApiRequest{
		Url:     requestUrl,
		Method:  http.MethodGet,
		Headers: headers,
	}
	res, err := utils.Request(req)
	if err != nil {
		return []types.GroupVariable{}, err
	}

	resModel := &[]types.GroupVariable{}
	err = json.Unmarshal(res.Data, resModel)
	if err != nil {
		return []types.GroupVariable{}, err
	}

	return *resModel, nil
}
