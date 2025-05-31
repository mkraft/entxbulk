package entxbulk

import (
	"log"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

type Option func(*Extension)

type Extension struct {
	typeCasts map[string]string
	goTypes   map[string]string
}

func WithTypeCasts(casts map[string]string) Option {
	return func(e *Extension) {
		e.typeCasts = casts
	}
}

func WithGoTypes(goTypes map[string]string) Option {
	return func(e *Extension) {
		e.goTypes = goTypes
	}
}

func NewExtension(opts ...Option) Extension {
	e := Extension{
		typeCasts: make(map[string]string),
	}
	for _, opt := range opts {
		opt(&e)
	}
	return e
}

type TypeCastAnnotation struct {
	TypeCasts map[string]string
}

func (TypeCastAnnotation) Name() string { return "TypeCast" }

type GoTypeCastAnnotation struct {
	GoTypes map[string]string
}

func (GoTypeCastAnnotation) Name() string { return "GoTypeCast" }

func (e Extension) Templates() []*gen.Template {
	_, currentFile, _, _ := runtime.Caller(0)
	dir := filepath.Dir(currentFile)
	templatePath := filepath.Join(dir, "template", "bulk_update.tmpl")

	log.Printf("Loading template from: %s", templatePath)

	t, err := gen.NewTemplate("bulk_update").
		Funcs(template.FuncMap{
			"ToExported": func(s string) string {
				if s == "" {
					return ""
				}
				// Convert snake_case to CamelCase
				parts := strings.Split(s, "_")
				for i, part := range parts {
					if len(part) > 0 {
						parts[i] = strings.ToUpper(part[:1]) + part[1:]
					}
				}
				return strings.Join(parts, "")
			},
			"until": func(n int) []int {
				a := make([]int, n)
				for i := range a {
					a[i] = i
				}
				return a
			},
			"hasPrefix": strings.HasPrefix,
		}).
		ParseFiles(templatePath)
	if err != nil {
		log.Printf("Failed to parse template: %v", err)
		panic(err)
	}

	log.Printf("Successfully loaded template")
	return []*gen.Template{t}
}

func (e Extension) Hooks() []gen.Hook {
	return nil
}

func (e Extension) Annotations() []entc.Annotation {
	return []entc.Annotation{
		TypeCastAnnotation{TypeCasts: e.typeCasts},
		GoTypeCastAnnotation{GoTypes: e.goTypes},
	}
}

func (Extension) Options() []entc.Option {
	return nil
}
