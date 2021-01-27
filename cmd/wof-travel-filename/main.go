package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/sfomuseum/go-flags/multi"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-travel-image/util" // PLEASE RECONCILE ME
	"github.com/whosonfirst/go-whosonfirst-travel/utils"      // PLEASE RECONCILE ME
	"log"
)

func main() {

	var sources multi.MultiString
	flag.Var(&sources, "source", "One or more filesystem based sources to use to read WOF ID data, which may or may not be part of the sources to graph. This is work in progress.")

	flag.Parse()

	ctx := context.Background()

	r, err := reader.NewMultiReaderFromURIs(ctx, sources...)

	if err != nil {
		log.Fatal(err)
	}

	for _, str_id := range flag.Args() {

		f, err := utils.LoadFeatureFromString(r, str_id)

		if err != nil {
			log.Fatal(err)
		}

		fname := util.Filename(f)
		fmt.Println(fname)
	}
}
