package swagger

import (
	"encoding/json"
	"fmt"

	"github.com/swaggo/swag"
)

// ServerInfo describes the combined API that will appear in the merged spec.
type ServerInfo struct {
	Title       string
	Description string
	Version     string
	// Host is used as the "host" field in swagger (e.g. "localhost:8080").
	// Leave empty to let swagger-ui use the current browser host.
	Host     string
	BasePath string
	Schemes  []string
}

// MergedSpecName is the swag instance name used for the combined merged spec.
// Callers must pass ginSwagger.InstanceName(swagger.MergedSpecName) when mounting
// the swagger-ui handler so it serves this spec.
const MergedSpecName = "merged"

// MergeAndRegister takes all doc JSONs from registry, merges paths, definitions,
// and securityDefinitions from every module, then registers the result as the
// "merged" swag spec so ginSwagger serves it at /swagger/index.html.
func MergeAndRegister(registry *Registry, info ServerInfo) error {
	return mergeAndRegister(registry.Docs(), info)
}

func mergeAndRegister(docs map[string]string, info ServerInfo) error {
	schemes := info.Schemes
	if len(schemes) == 0 {
		schemes = []string{"http", "https"}
	}

	paths := map[string]interface{}{}
	definitions := map[string]interface{}{}
	securityDefs := map[string]interface{}{
		// BearerAuth is used across all modules; seed it so it's always present.
		"BearerAuth": map[string]interface{}{
			"type":        "apiKey",
			"name":        "Authorization",
			"in":          "header",
			"description": "JWT Bearer token – format: 'Bearer <token>'",
		},
	}

	for moduleName, docJSON := range docs {
		var doc map[string]interface{}
		if err := json.Unmarshal([]byte(docJSON), &doc); err != nil {
			// Log and skip malformed modules rather than aborting the whole merge.
			fmt.Printf("[swagger/merger] skipping module %q: invalid JSON: %v\n", moduleName, err)
			continue
		}

		if p, ok := doc["paths"].(map[string]interface{}); ok {
			for k, v := range p {
				paths[k] = v
			}
		}
		if d, ok := doc["definitions"].(map[string]interface{}); ok {
			for k, v := range d {
				definitions[k] = v
			}
		}
		if s, ok := doc["securityDefinitions"].(map[string]interface{}); ok {
			for k, v := range s {
				securityDefs[k] = v
			}
		}
	}

	merged := map[string]interface{}{
		"swagger": "2.0",
		"info": map[string]interface{}{
			"title":       info.Title,
			"description": info.Description,
			"version":     info.Version,
			"contact":     map[string]interface{}{},
		},
		"host":                info.Host,
		"basePath":            info.BasePath,
		"schemes":             schemes,
		"securityDefinitions": securityDefs,
		"paths":               paths,
		"definitions":         definitions,
	}

	data, err := json.MarshalIndent(merged, "", "  ")
	if err != nil {
		return fmt.Errorf("swagger/merger: marshal: %w", err)
	}

	// Register as the "merged" swag instance (not the default "swagger") so it
	// does not clash with per-module init() registrations that use "swagger",
	// "calendar", "booking", etc. Callers mount ginSwagger with
	// ginSwagger.InstanceName(swagger.MergedSpecName).
	spec := &swag.Spec{
		Version:          info.Version,
		Host:             info.Host,
		BasePath:         info.BasePath,
		Schemes:          schemes,
		Title:            info.Title,
		Description:      info.Description,
		InfoInstanceName: MergedSpecName,
		// SwaggerTemplate is served verbatim – no template variables here because
		// we already resolved all values above.
		SwaggerTemplate: string(data),
	}
	swag.Register(MergedSpecName, spec)
	return nil
}
