package httpdoc

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/mercari/go-httpdoc/static"
)

// defaultTmpl is default template file to use.
var defaultTmpl = "tmpl/doc.md.tmpl"

// Generate writes documentation into the given file. Generation is skipped
// if EnvHTTPDoc is empty. If directory does not exist or any, it returns error.
func (d *Document) Generate(path string) error {

	// Only generate documentation when EnvHttpDoc has non-empty value
	if os.Getenv(EnvHTTPDoc) == "" {
		return nil
	}

	path, _ = filepath.Abs(path)
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	return d.generate(f)
}

func (d *Document) generate(w io.Writer) error {
	if d.tmpl == "" {
		d.tmpl = defaultTmpl
	}

	buf, err := static.Asset(d.tmpl)
	if err != nil {
		return err
	}

	return d.tmplExecute(w, string(buf))
}

func (d *Document) tmplExecute(w io.Writer, text string) error {
	tmpl, err := template.New("httpdoc").Funcs(funcMap()).Parse(text)
	if err != nil {
		return err
	}

	if err := tmpl.Execute(w, d); err != nil {
		return err
	}
	return nil
}

func funcMap() template.FuncMap {
	return template.FuncMap{
		"lower": strings.ToLower,
		"stripslash": func(s string) string {
			return strings.Replace(s, "/", "", -1)
		},
	}
}
