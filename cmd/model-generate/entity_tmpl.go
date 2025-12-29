package main

const entityTmpl = `
{{ $TableName := .BigCamelName }}
{{ $FirstLetter := .FirstLetter }}
package entity

// {{ $TableName }} {{ .Comment }}
type {{ $TableName }} struct {
{{ range $index, $field := .Fields }}{{ $field.BigCamelName }} {{ $field.DataType }}
{{ end }}
err error
}


{{ $SliceStruct := printf "%sSlice" $TableName }}

type {{ $SliceStruct }} []*{{ $TableName }}

func New{{ $TableName }}() *{{ $TableName }} {
    return &{{ $TableName }}{}
}

// region setter
{{ range $index, $field := .Fields }}
func ({{ $FirstLetter }} *{{ $TableName }}) Set{{ $field.BigCamelName }}(v {{ $field.DataType }}) *{{ $TableName }} {
	{{ $FirstLetter }}.{{ $field.BigCamelName }} = v
    return {{ $FirstLetter }}
}
{{ end }}
// endregion

func ({{ $FirstLetter }} *{{ $TableName }}) Valid() error {
    return {{ $FirstLetter }}.err
}

`
