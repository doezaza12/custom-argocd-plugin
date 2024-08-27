package cmd

import (
	"fmt"
	"log"
	"regexp"

	"github.com/doezaza12/custom-argocd-plugin/pkg/gitlab"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "",
	Long:  ``,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var manifests []unstructured.Unstructured
		var err error

		// read input from stdin
		manifests, err = readManifestData(cmd.InOrStdin())
		if err != nil {
			log.Fatal(err)
		}

		for _, manifest := range manifests {
			existingAnnotations := manifest.GetAnnotations()

			if manifest.GetKind() == "Secret" {
				id := uuid.New()

				if existingAnnotations == nil {
					existingAnnotations = map[string]string{
						"checksum": id.String(),
					}
				} else {
					existingAnnotations["checksum"] = id.String()
				}

				manifest.SetAnnotations(existingAnnotations)
			}

			if val, ok := existingAnnotations["abc.argocd.io/path"]; ok {
				listGroupVariables, err := gitlab.ListGitLabVariables(val)
				if err != nil {
					log.Fatal(err)
				}

				// Construct map to easier retrieve
				secretMap := make(map[string]string)
				for _, groupVariable := range listGroupVariables {
					secretMap[groupVariable.Key] = groupVariable.Value
				}

				if obj, ok := manifest.Object["data"].(map[string]interface{}); ok {
					for key, val := range obj {
						matched, err := regexp.MatchString("<[a-zA-Z0-9_]*>", val.(string))
						if err != nil {
							log.Fatal(err)
						}

						if matched {
							transformedVal := (val.(string))[1 : len(val.(string))-1]
							obj[key] = secretMap[transformedVal]
						}
					}
				}
			}

			output, err := yaml.Marshal(manifest.Object)

			if err != nil {
				log.Fatal(err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "%s---\n", output)
		}
	},
}
