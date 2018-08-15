package main

import (
	"flag"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"github.com/whosonfirst/go-whosonfirst-image"
	"github.com/whosonfirst/go-whosonfirst-svg"
	"github.com/whosonfirst/warning"
	"log"
	"os"
)

func main() {

	var width = flag.Int("width", 1024, "...")
	var height = flag.Int("height", 1024, "...")

	var style = flag.String("style", "", "...")
	var fill = flag.String("fill", "", "...")

	flag.Parse()

	opts := image.NewDefaultOptions()

	opts.Writer = os.Stdout // this is redundant but whatever
	opts.Width = *width
	opts.Height = *height

	switch *style {
	case "":
		// pass
	case "dopplr":
		opts.StyleFunction = svg.NewDopplrStyleFunction()
	case "fill":
		if *fill == "" {
			log.Fatal("Missing -fill colour")
		}
		opts.StyleFunction = svg.NewFillStyleFunction(*fill)
	default:
		log.Fatal("Invalid style")
	}

	for _, path := range flag.Args() {

		f, err := feature.LoadFeatureFromFile(path)

		if err != nil && !warning.IsWarning(err) {
			log.Fatal(err)
		}

		err = image.FeatureToPNG(f, opts)

		if err != nil {
			log.Fatal(err)
		}
	}
}
