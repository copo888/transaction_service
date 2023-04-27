package {{.PkgName}}

import (
	"net/http"
	{{if .HasRequest}}"{{.CommonPath}}/common/vaildx"
    {{end}}"{{.CommonPath}}/common/responsex"
    {{if .HasRequest}}"encoding/json"
	{{end}}{{if .After1_1_10}}{{if .HasRequest}}"github.com/zeromicro/go-zero/rest/httpx"{{end}}{{end}}
    {{if .HasRequest}}"go.opentelemetry.io/otel/attribute"
    {{end}}"go.opentelemetry.io/otel/trace"
	{{.ImportPackages}}
)

func {{.HandlerName}}(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		span := trace.SpanFromContext(r.Context())
        defer span.End()

		{{if .HasRequest}}var req types.{{.RequestType}}

        if err := httpx.ParseJsonBody(r, &req); err != nil {
            responsex.Json(w, r, responsex.FAIL, nil, err)
            return
        }

		if err := vaildx.Validator.Struct(req); err != nil {
			responsex.Json(w, r, responsex.INVALID_PARAMETER, nil, err)
			return
		}

		if requestBytes, err := json.Marshal(req); err == nil {
            span.SetAttributes(attribute.KeyValue{
                Key:   "request",
                Value: attribute.StringValue(string(requestBytes)),
            })
        }

		{{end}}l := {{.LogicName}}.New{{.LogicType}}(r.Context(), ctx)
		{{if .HasResp}}resp, {{end}}err := l.{{.Call}}({{if .HasRequest}}&req{{end}})
		if err != nil {
			responsex.Json(w, r, err.Error(), nil, err)
		} else {
			{{if .HasResp}}responsex.Json(w, r, responsex.SUCCESS, resp, err){{else}}responsex.Json(w, r, responsex.SUCCESS, nil, err){{end}}
		}
	}
}
