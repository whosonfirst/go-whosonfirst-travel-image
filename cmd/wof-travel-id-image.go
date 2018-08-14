package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-cli/flags"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-readwrite/reader"
	"github.com/whosonfirst/go-whosonfirst-travel"
	"github.com/whosonfirst/go-whosonfirst-travel-image/render"
	"github.com/whosonfirst/go-whosonfirst-travel/utils"
	"log"
	"os"
	"path/filepath"
	"sync"
)

func main() {

	var sources flags.MultiString
	flag.Var(&sources, "source", "One or more filesystem based sources to use to read WOF ID data, which may or may not be part of the sources to graph. This is work in progress.")

	var follow flags.MultiString
	flag.Var(&follow, "follow", "...")

	out := flag.String("out", "", "...")

	labels := flag.Bool("labels", false, "...")
	parent_id := flag.Bool("parent", false, "...")
	supersedes := flag.Bool("supersedes", false, "...")
	superseded_by := flag.Bool("superseded-by", false, "...")
	hierarchies := flag.Bool("hierarchies", false, "...")
	singleton := flag.Bool("singleton", true, "...")
	timings := flag.Bool("timings", false, "...")

	html := flag.Bool("html", false, "...")

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

		images := make([]*render.Image, 0)
		mu := new(sync.RWMutex)

		cb := func(f geojson.Feature, step int64) error {

			prefix := fmt.Sprintf("%03d", step)

			root := filepath.Join(abs_root, str_id)

			opts := &render.RenderOptions{
				Labels: *labels,
				Root:   root,
				Prefix: prefix,
			}

			im, err := render.RenderFeatureAsPNG(f, opts)

			if err != nil {
				return err
			}

			mu.Lock()
			defer mu.Unlock()

			images = append(images, im)
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

		if *html {

			root := filepath.Join(abs_root, str_id)

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
