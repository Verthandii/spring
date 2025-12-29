package main

import (
	"fmt"
	"log"
	"text/template"

	"github.com/Verthandii/spring/cmd/model-generate/writer"
)

func generateEntity(tables []*Table) {
	t := template.New("entity_tmpl")
	w := writer.NewTemplateWriter(
		t,
		&writer.Config{
			Terminal: Cfg.Terminal,
			Path:     "./common/domain/entity/",
		},
	)

	t, err := t.Parse(entityTmpl)
	if err != nil {
		log.Fatalln("template parse err:", err)
	}
	for _, table := range tables {
		if err = w.Write(table, fmt.Sprintf("%s.go", table.Name)); err != nil {
			log.Fatalln("write occurred error:", err)
		}
	}
}
