package main

// this is one of those scripts that does too much stuff and should
// almost certainly be chunked out in to package/library code but
// today is is not... (20180807/thisisaaronland)

import (
	"context"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-bindata-html-template"
	"github.com/whosonfirst/go-whosonfirst-cli/flags"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/whosonfirst"
	"github.com/whosonfirst/go-whosonfirst-image"
	"github.com/whosonfirst/go-whosonfirst-readwrite/reader"
	"github.com/whosonfirst/go-whosonfirst-travel"
	"github.com/whosonfirst/go-whosonfirst-travel-image/assets/html"
	"github.com/whosonfirst/go-whosonfirst-travel/utils"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"sync"
)

func OpenFilehandle(path string) (*os.File, error) {

	abs_path, err := filepath.Abs(path)

	if err != nil {
		return nil, err
	}

	root := filepath.Dir(abs_path)

	_, err = os.Stat(root)

	if os.IsNotExist(err) {

		err := os.MkdirAll(root, 0755)

		if err != nil {
			return nil, err
		}
	}

	return os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
}

func main() {

	var sources flags.MultiString
	flag.Var(&sources, "source", "One or more filesystem based sources to use to read WOF ID data, which may or may not be part of the sources to graph. This is work in progress.")

	var follow flags.MultiString
	flag.Var(&follow, "follow", "...")

	out := flag.String("out", "", "...")

	parent_id := flag.Bool("parent", false, "...")
	supersedes := flag.Bool("supersedes", false, "...")
	superseded_by := flag.Bool("superseded-by", false, "...")
	hierarchies := flag.Bool("hierarchies", false, "...")
	singleton := flag.Bool("singleton", true, "...")
	timings := flag.Bool("timings", false, "...")

	index_page := flag.Bool("html", false, "...")

	flag.Parse()

	if *out == "" {

		cwd, err := os.Getwd()

		if err != nil {
			log.Fatal(err)
		}

		*out = cwd
	}

	abs_root, err := filepath.Abs(*out)

	if err != nil {
		log.Fatal(err)
	}

	r, err := reader.NewMultiReaderFromStrings(sources...)

	if err != nil {
		log.Fatal(err)
	}

	for _, str_id := range flag.Args() {

		// see this? we're making a separate traveler and callback thingy
		// for each ID because, for example, we may want to include the same
		// hierarchy or parents or whatever only once for each ID but at
		// least once for each ID - that sort of thing (20180807/thisisaaronland)

		f, err := utils.LoadFeatureFromString(r, str_id)

		if err != nil {
			log.Fatal(err)
		}

		images := make([][]string, 0)
		mu := new(sync.RWMutex)

		cb := func(f geojson.Feature) error {

			id := f.Id()

			fname := fmt.Sprintf("%s.png", id)

			root := filepath.Join(abs_root, str_id)
			path := filepath.Join(root, fname)

			fh, err := OpenFilehandle(path)

			if err != nil {
				return err
			}

			opts := image.NewDefaultOptions()
			im, err := image.FeatureToImage(f, opts)

			if err != nil {
				return err
			}

			err = png.Encode(fh, im)

			if err != nil {
				return err
			}

			fh.Close()

			label := whosonfirst.LabelOrDerived(f)
			label = fmt.Sprintf("%s %s", id, label)

			mu.Lock()
			defer mu.Unlock()

			images = append(images, []string{fname, label})
			return nil
		}

		opts, err := travel.DefaultTravelOptions()

		if err != nil {
			log.Fatal(err)
		}

		opts.Reader = r

		opts.Callback = cb
		opts.Singleton = *singleton
		opts.ParentID = *parent_id
		opts.Hierarchy = *hierarchies
		opts.Supersedes = *supersedes
		opts.SupersededBy = *superseded_by
		opts.Timings = *timings

		tr, err := travel.NewTraveler(opts)

		if err != nil {
			log.Fatal(err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		err = tr.TravelFeature(ctx, f)

		if err != nil {
			log.Fatal(err)
		}

		// speaking of things that really need to go in a separate package...
		// (20180807/thisisaaronland)
		
		if *index_page {

			fname := "index.html"

			root := filepath.Join(abs_root, str_id)
			path := filepath.Join(root, fname)

			fh, err := OpenFilehandle(path)

			if err != nil {
				log.Fatal(err)
			}

			defer fh.Close()

			type HTMLVars struct {
				ID     string
				Images [][]string
			}

			tpl := template.New("index", html.Asset)

			tpl, err = tpl.ParseFiles("templates/html/index.html")

			if err != nil {
				log.Fatal(err)
			}

			vars := HTMLVars{
				ID:     str_id,
				Images: images,
			}

			err = tpl.Execute(fh, vars)

			if err != nil {
				log.Fatal(err)
			}
		}
	}

}
