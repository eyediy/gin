// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"html/template"
	"net/http"
)

type (
	HTMLRender interface {
		Instance(string, interface{}) Render
		Funcs(fm template.FuncMap)
	}

	HTMLProduction struct {
		Template *template.Template
	}

	HTMLDebug struct {
		Files []string
		Glob  string
		FuncMaps template.FuncMap
	}

	HTML struct {
		Template *template.Template
		Name     string
		Data     interface{}
	}
)

var htmlContentType = []string{"text/html; charset=utf-8"}

func (r HTMLProduction) Funcs(fm template.FuncMap) {
	r.Template.Funcs(fm)
}
func (r HTMLDebug) Funcs(fm template.FuncMap) {
	for name, fn := range fm {
		r.FuncMaps[name] = fn
	}
}

func (r HTMLProduction) Instance(name string, data interface{}) Render {
	return HTML{
		Template: r.Template,
		Name:     name,
		Data:     data,
	}
}

func (r HTMLDebug) Instance(name string, data interface{}) Render {
	html := HTML{
		Template: r.LoadTemplate(),
		Name:     name,
		Data:     data,
	}
	if r.FuncMaps != nil {
		html.Template.Funcs(r.FuncMaps)
	}
	return html
}
func (r HTMLDebug) LoadTemplate() *template.Template {
	if len(r.Files) > 0 {
		return template.Must(template.ParseFiles(r.Files...))
	}
	if len(r.Glob) > 0 {
		return template.Must(template.ParseGlob(r.Glob))
	}
	panic("the HTML debug render was created without files or glob pattern")
}

func (r HTML) Render(w http.ResponseWriter) error {
	writeContentType(w, htmlContentType)
	if len(r.Name) == 0 {
		return r.Template.Execute(w, r.Data)
	}
	return r.Template.ExecuteTemplate(w, r.Name, r.Data)
}
