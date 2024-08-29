package cmd

import (
	"encoding/base64"
	"fmt"
	"log"
	"regexp"

	"github.com/doezaza12/custom-argocd-plugin/pkg/gitlab"
	"github.com/doezaza12/custom-argocd-plugin/pkg/types"
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
				if val, ok := existingAnnotations[types.AnnotationIndicator]; ok {
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
							decodedVal, err := base64.StdEncoding.DecodeString(val.(string))
							if err != nil {
								log.Fatal(err)
							}

							decodedStringVal := string(decodedVal)

							matched, err := regexp.MatchString("<[a-zA-Z0-9_]*>", decodedStringVal)
							if err != nil {
								log.Fatal(err)
							}

							if matched {
								transformedVal := decodedStringVal[1 : len(decodedStringVal)-1]
								obj[key] = base64.StdEncoding.EncodeToString([]byte(secretMap[transformedVal]))
							}
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
