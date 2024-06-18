package main

import (
	"log"
	"time"

	"github.com/gzuidhof/tygo/tygo"
)

func main() {
	var defaultTypeMappings = map[string]string{
		"time.Time": "string /* RFC3339 */",
	}

	var baseTypesOutputDir = "frontend/src/types"

	cfg := &tygo.Config{
		Packages: []*tygo.PackageConfig{
			{
				Path:         "github.com/cory-evans/record-rummage/internal/models",
				OutputPath:   baseTypesOutputDir + "/models.ts",
				TypeMappings: defaultTypeMappings,
			},
			{
				Path:       "github.com/zmb3/spotify/v2",
				OutputPath: baseTypesOutputDir + "/spotify/gen.ts",
				TypeMappings: map[string]string{
					"time.Time": "string /* RFC3339 */",
					"[]Image":   "Image[] | null",
				},
				ExcludeFiles: []string{
					"request_options.go",
					"spotify.go",
				},
				Frontmatter: `
import {
	ID,
	Numeric,
	URI,
	Image,
	Followers,
	Error,
} from "./base"`,
			},
		},
	}

	gen := tygo.New(cfg)
	gen.Generate()
	log.Println(time.Now().Format(time.RFC3339), "Done")
}
