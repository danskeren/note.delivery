package templates

import (
	"html/template"
	"io"
	"io/ioutil"
	"log"

	"github.com/packago/config"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"note.delivery/templates/tmpl"
)

var tmpls map[string]*template.Template = make(map[string]*template.Template)

func init() {
	if config.File().GetString("environment") == "development" {
		minifyCSS()
	}
	layout := minifyHTML(tmpl.LayoutHTML)

	tmpls["index.html"] = template.Must(template.New("layout").Parse(layout))
	template.Must(tmpls["index.html"].New("index").Parse(minifyHTML(tmpl.IndexHTML)))
	tmpls["note.html"] = template.Must(template.New("layout").Parse(layout))
	template.Must(tmpls["note.html"].New("note").Parse(minifyHTML(tmpl.NoteHTML)))

	tmpls["privacy-policy.html"] = template.Must(template.New("layout").Parse(layout))
	template.Must(tmpls["privacy-policy.html"].New("privacy-policy").Parse(minifyHTML(tmpl.PrivacyPolicyHTML)))

	tmpls["rate-limit.html"] = template.Must(template.New("layout").Parse(layout))
	template.Must(tmpls["rate-limit.html"].New("rate-limit").Parse(minifyHTML(tmpl.RateLimitHTML)))
	tmpls["not-found.html"] = template.Must(template.New("layout").Parse(layout))
	template.Must(tmpls["not-found.html"].New("not-found").Parse(minifyHTML(tmpl.NotFoundHTML)))
}

func Render(wr io.Writer, template string, data interface{}) {
	err := tmpls[template].Execute(wr, data)
	if err != nil {
		log.Panicf("[ERROR] Error rendering %s: %s\n", template, err)
	}
}

func minifyHTML(input string) string {
	min := minify.New()
	min.Add("text/html", &html.Minifier{})
	str, err := min.String("text/html", input)
	if err != nil {
		panic(err)
	}
	return str
}

func minifyCSS() {
	filepaths := []string{"./static/css/main.css"}

	min := minify.New()
	min.AddFunc("text/css", css.Minify)
	var minifiedCSS string
	for _, filepath := range filepaths {
		unminifiedCSS, err := ioutil.ReadFile(filepath)
		minCSS, err := min.String("text/css", string(unminifiedCSS))
		if err != nil {
			panic(err)
		}
		minifiedCSS += minCSS
	}
	err := ioutil.WriteFile("./static/main.min.css", []byte(minifiedCSS), 0644)
	if err != nil {
		panic(err)
	}
}
