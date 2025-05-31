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
}

func WithTypeCasts(casts map[string]string) Option {
	return func(e *Extension) {
		e.typeCasts = casts
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
				return strings.Title(s)
			},
			"until": func(n int) []int {
				a := make([]int, n)
				for i := range a {
					a[i] = i
				}
				return a
			},
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
	}
}

func (Extension) Options() []entc.Option {
	return nil
}
