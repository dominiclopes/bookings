package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/dominiclopes/bookings/pkg/config"
	"github.com/dominiclopes/bookings/pkg/models"
)

var app *config.AppConfig

func NewAppConfig(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	// get the template cache
	var tc map[string]*template.Template
	var err error
	if app.InProduction {
		tc = app.TemplateCache
	} else {
		tc, err = CreateTemplateCache()
		if err != nil {
			log.Fatalf("error creating the template cache, err: %v", err)
		}
	}

	// get the requested template from the cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatalf("error requested template %s not present in the template cache", tmpl)
	}

	// render the template
	buf := new(bytes.Buffer)
	td = AddDefaultData(td)
	err = t.Execute(buf, td)
	if err != nil {
		log.Printf("error executing the parsed template %s, err: %v", tmpl, err)
		return
	}
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Printf("error rendering the template %s, err: %v", tmpl, err)
		return
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	log.Println("Creating the template cache")
	myCache := map[string]*template.Template{}

	// get all the files ending with *.page.tmpl
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		log.Printf("error getting filepaths of all the page templates, err: %v", err)
		return myCache, err
	}

	// range through the matched files and parse them
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			log.Printf("error parsing the template %s, err: %v", name, err)
			return myCache, err
		}

		// get all the files ending with *.layout.tmpl
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			log.Printf("error getting filepaths of all the layout templates, err: %v", err)
			return myCache, err
		}

		// parse all the layout templates
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				log.Printf("error parsing the layout templates, err: %v", err)
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
