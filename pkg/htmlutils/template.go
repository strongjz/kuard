/*
Copyright 2017 The KUAR Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package htmlutils

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/jbeda/kuard/pkg/config"
	"github.com/jbeda/kuard/pkg/debugsitedata"
	"github.com/jbeda/kuard/pkg/sitedata"
)

type TemplateGroup struct {
	t *template.Template
}

func (g *TemplateGroup) Render(w http.ResponseWriter, name string, context interface{}) {
	t := g.GetTemplate(name)
	buf := &bytes.Buffer{}
	err := t.Execute(buf, context)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	w.WriteHeader(http.StatusOK)
	buf.WriteTo(w)
}

func (g *TemplateGroup) GetTemplate(name string) *template.Template {
	if g.t == nil || *config.Debug {
		g.t = g.LoadTemplates()
	}
	t := g.t.Lookup(name)
	if t == nil {
		panic(fmt.Sprintf("Could not load template %v", name))
	}
	return t
}

// LoadTemplates loads the templates for our toy server
func (g *TemplateGroup) LoadTemplates() *template.Template {
	assetDir := sitedata.AssetDir
	asset := sitedata.Asset
	if *config.Debug {
		assetDir = debugsitedata.AssetDir
		asset = debugsitedata.Asset
	}

	tFiles, err := assetDir("templates")
	if err != nil {
		panic(err)
	}

	t := template.New("").Funcs(FuncMap())

	for _, tFile := range tFiles {
		fullName := path.Join("templates", tFile)
		data, err := asset(fullName)
		if err != nil {
			continue
		}
		log.Printf("Loading template for %v", tFile)
		_, err = t.New(tFile).Parse(string(data))
		if err != nil {
			log.Printf("ERROR: Could parse template %v: %v", tFile, err)
		}
	}
	return t
}
