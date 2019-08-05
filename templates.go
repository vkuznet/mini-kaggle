package main

// templates module
//
// Copyright (c) 2019 - Valentin Kuznetsov <vkuznet@gmail.com>
//

import (
	"bytes"
	"html/template"
	"path/filepath"
)

// consume list of templates and release their full path counterparts
func fileNames(tdir string, filenames ...string) []string {
	flist := []string{}
	for _, fname := range filenames {
		flist = append(flist, filepath.Join(tdir, fname))
	}
	return flist
}

// parse template with given data
func parseTmpl(tdir, tmpl string, data interface{}) string {
	buf := new(bytes.Buffer)
	filenames := fileNames(tdir, tmpl)
	funcMap := template.FuncMap{
		// The name "oddFunc" is what the function will be called in the template text.
		"oddFunc": func(i int) bool {
			if i%2 == 0 {
				return true
			}
			return false
		},
		"incFunc": func(i int) int {
			return i + 1
		},
	}
	t := template.Must(template.New(tmpl).Funcs(funcMap).ParseFiles(filenames...))
	err := t.Execute(buf, data)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

// ServerTemplates structure
type ServerTemplates struct {
	top, bottom, searchForm, cards, dasError, dasKeys, dasZero string
}

// Top method for ServerTemplates structure
func (q ServerTemplates) Top(tdir string, tmplData map[string]interface{}) string {
	if q.top != "" {
		return q.top
	}
	q.top = parseTmpl(Config.Templates, "top.tmpl", tmplData)
	return q.top
}

// Bottom method for ServerTemplates structure
func (q ServerTemplates) Bottom(tdir string, tmplData map[string]interface{}) string {
	if q.bottom != "" {
		return q.bottom
	}
	q.bottom = parseTmpl(Config.Templates, "bottom.tmpl", tmplData)
	return q.bottom
}

// Home method for ServerTemplates structure
func (q ServerTemplates) Home(tdir string, tmplData map[string]interface{}) string {
	if q.searchForm != "" {
		return q.searchForm
	}
	q.searchForm = parseTmpl(Config.Templates, "home.tmpl", tmplData)
	return q.searchForm
}

// Confirm method for ServerTemplates structure
func (q ServerTemplates) Confirm(tdir string, tmplData map[string]interface{}) string {
	if q.top != "" {
		return q.top
	}
	q.top = parseTmpl(Config.Templates, "confirm.tmpl", tmplData)
	return q.top
}

// Dashboard method for ServerTemplates structure
func (q ServerTemplates) Dashboard(tdir string, tmplData map[string]interface{}) string {
	if q.dasError != "" {
		return q.dasError
	}
	q.dasError = parseTmpl(Config.Templates, "dashboard.tmpl", tmplData)
	return q.dasError
}

// PrivateLogin method for ServerTemplates structure
func (q ServerTemplates) PrivateLogin(tdir string, tmplData map[string]interface{}) string {
	if q.dasError != "" {
		return q.dasError
	}
	q.dasError = parseTmpl(Config.Templates, "login.tmpl", tmplData)
	return q.dasError
}
