package render

import (
	"github.com/whosonfirst/go-whosonfirst-travel-image/templates/html"
	"html/template"
	"log"
	"path/filepath"
)

func RenderIndexForImages(images []*Image, opts *RenderOptions) error {

	fname := "index.html"

	path := filepath.Join(opts.Root, fname)

	fh, err := util.OpenFilehandle(path)

	if err != nil {
		return err
	}

	defer fh.Close()

	type IndexVars struct {
		Title  string
		Images []*Image
	}

	tpl := template.New("images")

	tpl, err = tpl.ParseFS(html.FS, "*.html")

	if err != nil {
		log.Fatal(err)
	}

	vars := IndexVars{
		Title:  "",
		Images: images,
	}

	err = tpl.Execute(fh, vars)

	if err != nil {
		return err
	}

	return nil
}
