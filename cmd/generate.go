package cmd

import (
	"fmt"
	"log"

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
			if manifest.GetKind() == "Secret" {

				existingAnnotations := manifest.GetAnnotations()
				if existingAnnotations == nil {
					existingAnnotations = map[string]string{
						"trustme": "imdevops",
					}
				} else {
					existingAnnotations["trustme"] = "imdevops"
				}

				manifest.SetAnnotations(existingAnnotations)
			}

			output, err := yaml.Marshal(manifest.Object)

			if err != nil {
				log.Fatal(err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "%s---\n", output)
		}
	},
}
