package main

import (
	"fmt"
	"log"
	"text/template"

	"github.com/Verthandii/spring/cmd/model-generate/writer"
)

func generateModel(tables []*Table) {
	t := template.New("model_tmpl")
	t, err := t.Parse(modelTmpl)
	if err != nil {
		log.Fatalln("template parse err:", err)
	}

	for _, table := range tables {
		w := writer.NewTemplateWriter(
			t,
			&writer.Config{
				Terminal: Cfg.Terminal,
				Path:     "./common/infra/persistence/model/",
			},
		)
		if err = w.Write(table, fmt.Sprintf("%s%s.go", _prefix, table.Name)); err != nil {
			log.Fatalln("write occurred error:", err)
		}
	}
}
