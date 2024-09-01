package cmd

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

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
		var swapGeneratedName map[string]string

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

							matched, err := regexp.MatchString("<[a-zA-Z0-9_]+>", decodedStringVal)
							if err != nil {
								log.Fatal(err)
							}

							if matched {
								transformedVal := decodedStringVal[1 : len(decodedStringVal)-1]
								obj[key] = base64.StdEncoding.EncodeToString([]byte(secretMap[transformedVal]))
							}
						}

						// Calculate new hash based-on current secrets for secret name suffix
						data, err := json.Marshal(manifest.Object["data"])
						if err != nil {
							log.Fatal(err)
						}

						manifestName := manifest.GetName()
						hashSuffix := fmt.Sprintf("%x", sha256.Sum256([]byte(data)))
						enc := []rune(hashSuffix[:10])
						if _, ok := swapGeneratedName[manifestName]; ok {
							swapGeneratedName[manifest.GetName()] = string(enc)
						} else {
							swapGeneratedName = map[string]string{
								manifestName: fmt.Sprintf("%s%s", manifestName[:len(manifestName)-10], string(enc)),
							}
						}
					}
				}
			}

			output, err := yaml.Marshal(manifest.Object)
			if err != nil {
				log.Fatal(err)
			}

			var finalOutput string
			for key, val := range swapGeneratedName {
				finalOutput = strings.ReplaceAll(string(output), key, val)
			}

			if finalOutput != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "%s---\n", finalOutput)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "%s---\n", output)
			}
		}
	},
}
