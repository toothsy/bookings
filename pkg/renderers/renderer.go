package renderers

import (
	"bytes"
	"github/toothsy/bookings/pkg/config"
	"github/toothsy/bookings/pkg/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var appConf *config.AppConfig

func AddDefaultTemplateData(td *models.TemplateData) *models.TemplateData {
	return td
}
func RenderTemplate(w http.ResponseWriter, tmpl string, templateData *models.TemplateData) {
	// create cache
	var templateCache map[string]*template.Template
	if appConf.UseCache {
		templateCache = appConf.TemplateCache
	} else {
		templateCache, _ = CreateTemplateCache()
	}
	// get the template from cache
	template, okay := templateCache[filepath.Base(tmpl)]
	if !okay {
		log.Fatal("couldnt get the cache")
	}

	buff := new(bytes.Buffer)
	templateData = AddDefaultTemplateData(templateData)
	_ = template.Execute(buff, templateData)

	//render the template
	_, err := buff.WriteTo(w)
	if err != nil {
		log.Fatal("Coudlnt render", err)

	}

}

// allows outside function t o initialise appConfig
func SetConfig(a *config.AppConfig) {
	appConf = a
}
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}
	//getting all tmpl files
	pages, err := filepath.Glob("../../templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		fileName := filepath.Base(page)
		templateSet, err := template.New(fileName).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		layouts, err := filepath.Glob("../../templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}
		if len(layouts) > 0 {
			templateSet, err = templateSet.ParseGlob("../../templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}
		myCache[fileName] = templateSet

	}
	return myCache, nil

}
