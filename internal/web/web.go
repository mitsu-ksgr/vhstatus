package web

import (
	"html/template"
	"net/http"

	"github.com/mitsu-ksgr/vhstatus/internal/vhstatus"
)

//-----------------------------------------------------------------------------
// templateDirPath ... path to the template to use when rendering html.
var templateDirPath string

func getTemplateDirPath() string {
	return templateDirPath
}

func SetTemplateDirPath(p string) {
	if len(p) > 0 && p[len(p)-1:] != "/" {
		templateDirPath = p + "/"
	} else {
		templateDirPath = p
	}
}

//-----------------------------------------------------------------------------
// funcFetchVHStatus
var funcFetchVHStatus func() vhstatus.Params

func getVHStatusParams() vhstatus.Params {
	if funcFetchVHStatus == nil {
		return vhstatus.Params{
			Status: "WARN: VHStatus - funcFetchVHStatus is nil",
		}
	}
	return funcFetchVHStatus()
}

func SetFechVHStatusParamsFunc(f func() vhstatus.Params) {
	funcFetchVHStatus = f
}

//-----------------------------------------------------------------------------
// Rendering helper
func render404(w http.ResponseWriter) {
	if temp, err := template.ParseFiles(getTemplateDirPath() + "404.html"); err == nil {
		if err := temp.Execute(w, nil); err == nil {
			w.WriteHeader(404)
			return
		}
	}
	http.Error(w, "404 Not Found", 404)
}

func render500(w http.ResponseWriter) {
	if temp, err := template.ParseFiles(getTemplateDirPath() + "500.html"); err == nil {
		if err := temp.Execute(w, nil); err == nil {
			w.WriteHeader(500)
			return
		}
	}
	http.Error(w, "500 Internal Server Error", 500)
}
