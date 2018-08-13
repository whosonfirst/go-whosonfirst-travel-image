package main

import (
	"flag"
	"github.com/whosonfirst/go-whosonfirst-cli/flags"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/whosonfirst"
	"github.com/whosonfirst/go-whosonfirst-travel/traveler"
	"log"
	"sync"
)

func main() {

	var belongs_to flags.MultiInt64
	flag.Var(&belongs_to, "belongs-to", "...")

	mode := flag.String("mode", "repo", "...")

	flag.Parse()

	lookup := make(map[int64]geojson.Feature)
	idx := make(map[int64][]int64)

	mu := new(sync.RWMutex)

	cb := func(f geojson.Feature, belongsto_id int64) error {

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

	for belongsto, descendants := range idx {
		log.Println(belongsto, descendants)
	}
}
