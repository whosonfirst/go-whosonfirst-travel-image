package util

import (
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/whosonfirst"
)

func Label(f geojson.Feature) string {

	id := f.Id()
	pt := f.Placetype()

	label := whosonfirst.LabelOrDerived(f)

	return fmt.Sprintf("%s %s %s", id, label, pt)
}
