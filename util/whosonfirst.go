package util

import (
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/whosonfirst"
	"regexp"
	"strings"
)

var re_name *regexp.Regexp

func init() {
	re_name = regexp.MustCompile(`[^a-z0-9\-]+`)
}

func Filename(f geojson.Feature) string {

	id := f.Id()
	name := f.Name()

	name = strings.ToLower(name)
	name = strings.Replace(name, " ", "-", -1)
	name = re_name.ReplaceAllString(name, "")

	ds := whosonfirst.DateSpan(f)

	// always put the dates at the beginning for sorting
	// always put the ID at the end so it's easy to tease apart

	return fmt.Sprintf("%s-%s-%s", ds, name, id)
}

func FilenameWithExtension(f geojson.Feature, ext string) string {

	if strings.HasPrefix(ext, ".") {
		ext = strings.TrimLeft(ext, ".")
	}

	fname := Filename(f)

	return fmt.Sprintf("%s.%s", fname, ext)
}

func Label(f geojson.Feature) string {

	id := f.Id()
	pt := f.Placetype()

	label := whosonfirst.LabelOrDerived(f)

	return fmt.Sprintf("%s %s %s", id, label, pt)
}
