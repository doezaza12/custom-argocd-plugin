package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
)

func readFilesAsManifests(paths []string) (result []unstructured.Unstructured, errs []error) {
	for _, path := range paths {
		rawdata, err := os.ReadFile(path)
		if err != nil {
			errs = append(errs, fmt.Errorf("could not read file: %s from disk: %s", path, err))
		}
		manifest, err := readManifestData(bytes.NewReader(rawdata))
		if err != nil {
			errs = append(errs, fmt.Errorf("could not read file: %s from disk: %s", path, err))
		}
		result = append(result, manifest...)
	}

	return result, errs
}

func readManifestData(yamlData io.Reader) ([]unstructured.Unstructured, error) {
	decoder := k8syaml.NewYAMLOrJSONDecoder(yamlData, 1)

	var manifests []unstructured.Unstructured
	for {
		nxtManifest := unstructured.Unstructured{}
		err := decoder.Decode(&nxtManifest)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		// Skip empty manifests
		if len(nxtManifest.Object) > 0 {
			manifests = append(manifests, nxtManifest)
		}
	}

	return manifests, nil
}
