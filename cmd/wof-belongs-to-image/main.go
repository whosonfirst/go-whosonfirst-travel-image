package main

// ./bin/wof-belongs-to -html -out ./belongs-to -source fs:///usr/local/data/sfomuseum-data-architecture/data -include-placetype concourse -include-placetype wing -belongs-to 1159157271 /usr/local/data/sfomuseum-data-architecture/

import (
	"context"
	"flag"
	"fmt"
	"github.com/sfomuseum/go-flags/multi"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/whosonfirst"
	_ "github.com/whosonfirst/go-whosonfirst-index/fs"
	"github.com/whosonfirst/go-whosonfirst-travel-image/render"
	"github.com/whosonfirst/go-whosonfirst-travel-image/util" // PLEASE RECONCILE ME
	"github.com/whosonfirst/go-whosonfirst-travel/traveler"
	"github.com/whosonfirst/go-whosonfirst-travel/utils" // PLEASE RECONCILE ME
	"log"
	"os"
	"path/filepath"
	"sync"
)

func main() {

	var belongs_to multi.MultiInt64
	flag.Var(&belongs_to, "belongs-to", "...")

	var include_placetype multi.MultiString
	flag.Var(&include_placetype, "include-placetype", "...")

	var exclude_placetype multi.MultiString
	flag.Var(&exclude_placetype, "exclude-placetype", "...")

	var sources multi.MultiString
	flag.Var(&sources, "source", "One or more filesystem based sources to use to read WOF ID data, which may or may not be part of the sources to graph. This is work in progress.")

	out := flag.String("out", "", "...")
	html := flag.Bool("html", false, "...")
	labels := flag.Bool("labels", false, "...")

	mode := flag.String("mode", "repo://", "...")

	flag.Parse()

	if *out == "" {

		cwd, err := os.Getwd()

		if err != nil {
			log.Fatal(err)
		}

		*out = cwd
	}

	ctx := context.Background()

	r, err := reader.NewMultiReaderFromURIs(ctx, sources...)

	if err != nil {
		log.Fatal(err)
	}

	lookup := make(map[int64]geojson.Feature)
	idx := make(map[int64][]int64)

	mu := new(sync.RWMutex)

	cb := func(f geojson.Feature, belongsto_id int64) error {

		pt := f.Placetype()
		// log.Println("placetype", pt)

		if len(include_placetype) > 0 {

			if !include_placetype.Contains(pt) {
				return nil
			}
		}

		if len(exclude_placetype) > 0 {

			if exclude_placetype.Contains(pt) {
				return nil
			}
		}

		mu.Lock()
		defer mu.Unlock()

		_, ok := lookup[belongsto_id]

		if ok {
			return nil
		}

		wof_id := whosonfirst.Id(f)

		lookup[wof_id] = f

		descendants, ok := idx[belongsto_id]

		if !ok {
			descendants = make([]int64, 0)
		}

		descendants = append(descendants, wof_id)
		idx[belongsto_id] = descendants

		return nil
	}

	t, err := traveler.NewDefaultBelongsToTraveler()
	t.Mode = *mode
	t.BelongsTo = belongs_to
	t.Callback = cb

	paths := flag.Args()
	err = t.Travel(paths...)

	if err != nil {
		log.Fatal(err)
	}

	//

	render_features := func(first geojson.Feature, features ...geojson.Feature) ([]*render.Image, error) {

		fname := util.Filename(first)
		root := filepath.Join(*out, fname)

		images := make([]*render.Image, 0)

		for i, f := range features {

			prefix := fmt.Sprintf("%03d", i+1)

			opts := &render.RenderOptions{
				Labels: *labels,
				Root:   root,
				Prefix: prefix,
			}

			im, err := render.RenderFeatureAsPNG(f, opts)

			if err != nil {
				return nil, err
			}

			images = append(images, im)
		}

		return images, nil
	}

	for belongsto_id, descendants := range idx {

		features := make([]geojson.Feature, len(descendants)+1)

		f, err := utils.LoadFeature(r, belongsto_id)

		if err != nil {
			log.Fatal(err)
		}

		features[0] = f

		for i, id := range descendants {

			f, ok := lookup[id]

			if !ok {
				log.Fatal("can't find ID")
			}

			features[i+1] = f
		}

		first := features[0]

		images, err := render_features(first, features...)

		if err != nil {
			log.Fatal(err)
		}

		if *html {

			root := filepath.Join(*out, first.Id())

			opts := &render.RenderOptions{
				Root: root,
			}

			err = render.RenderIndexForImages(images, opts)

			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
