package handler

import (
	"fmt"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := pathsToUrls[r.URL.Path]
		if url != "" {
			http.Redirect(w, r, url, http.StatusPermanentRedirect)
		} else {
			fallback.ServeHTTP(w, r)
		}
	})
}

func YAMLHandler(fileName string, fallback http.Handler) (http.HandlerFunc, error) {
	pathsToUrls := parseYamlFile(fileName)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := pathsToUrls[r.URL.Path]
		if url != "" {
			http.Redirect(w, r, url, http.StatusPermanentRedirect)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}), nil
}

func parseYamlFile(fileName string) map[string]string {

	var rawMappings []map[string]interface{}

	pathsToUrls := make(map[string]string)

	bytes, _ := os.ReadFile(fileName)

	yaml.Unmarshal(bytes, &rawMappings)

	for _, rawMapping := range rawMappings {
		path, pathOK := rawMapping["path"].(string)
		url, urlOK := rawMapping["url"].(string)
		if pathOK && urlOK {
			pathsToUrls[path] = url
		}
	}

	fmt.Print(pathsToUrls)

	return pathsToUrls
}
