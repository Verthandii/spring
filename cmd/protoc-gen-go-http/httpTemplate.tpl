{{$svrType := .ServiceType}}
{{$svrName := .ServiceName}}

type {{.ServiceType}}HTTPServer interface {
{{- range .MethodSets}}
	{{- if ne .Comment ""}}
	{{.Comment}}
	{{- end}}
	{{.Name}}(*app.Context, *{{.Request}}) (*{{.Reply}}, error)
{{- end}}
}

func Register{{.ServiceType}}HTTPServer(srv *ihttp.Server, s {{.ServiceType}}HTTPServer) {
    g := srv.Group("")
	{{- range .Methods}}
	g.{{.Method}}("{{.Path}}", {{.Name}}{{.Num}}(s))
	{{- end}}
}

{{range .Methods}}
func {{.Name}}{{.Num}}(s {{$svrType}}HTTPServer) ihttp.Handler {
    return func(c *ihttp.Context) {
        req := &{{.Request}}{}
    	if err := c.ShouldBind(req); err != nil {
    		c.Failure(err)
    		return
    	}

    	if err := req.ValidateAll(); err != nil {
    		c.Failure(err)
    		return
    	}

    	resp, err := s.{{.Name}}(app.GetContext(c), req)
    	if err != nil {
    		c.Failure(err)
    		return
    	}
    	c.Success(resp)
    }
}
{{end}}

type {{.ServiceType}}HTTPClient interface {
{{- range .MethodSets}}
    {{- if ne .Comment ""}}
	{{.Comment}}
	{{- end}}
	{{.Name}}(ctx context.Context, req *{{.Request}}) (rsp *{{.Reply}}, err error)
{{- end}}
}

type {{.ServiceType}}HTTPClientImpl struct{
	cc *ihttp.Client
}

func New{{.ServiceType}}HTTPClient (client *ihttp.Client) {{.ServiceType}}HTTPClient {
	return &{{.ServiceType}}HTTPClientImpl{client}
}

{{range .MethodSets}}
func (c *{{$svrType}}HTTPClientImpl) {{.Name}}(ctx context.Context, in *{{.Request}}) (*{{.Reply}}, error) {
	var out {{.Reply}}
	pattern := "{{.Path}}"
	path := ihttp.EncodeURL(pattern, in, {{not .HasBody}})
	{{if .HasBody -}}
	err := c.cc.Invoke(ctx, "{{.Method}}", path, "application/json", in{{.Body}}, &out{{.ResponseBody}})
	{{else -}}
	err := c.cc.Invoke(ctx, "{{.Method}}", path, "application/json", nil, &out{{.ResponseBody}})
	{{end -}}
	if err != nil {
		return nil, err
	}
	return &out, err
}
{{end}}
