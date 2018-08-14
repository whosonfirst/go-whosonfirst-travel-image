package render

import (
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	wof_image "github.com/whosonfirst/go-whosonfirst-image"
	"github.com/whosonfirst/go-whosonfirst-travel-image/util"
	"golang.org/x/image/font"
	"golang.org/x/image/font/inconsolata"
	"golang.org/x/image/math/fixed"
	go_image "image"
	"image/color"
	"image/draw"
	"image/png"
	"path/filepath"
)

func RenderFeatureAsPNG(f geojson.Feature, opts *RenderOptions) (*Image, error) {

	fname := fmt.Sprintf("%s.png", f.Id())

	if opts.Prefix != "" {
		fname = fmt.Sprintf("%s-%s", opts.Prefix, fname)
	}

	path := filepath.Join(opts.Root, fname)

	im, err := RenderFeature(f, opts)

	if err != nil {
		return nil, err
	}

	fh, err := util.OpenFilehandle(path)

	if err != nil {
		return nil, err
	}

	err = png.Encode(fh, im)

	if err != nil {
		return nil, err
	}

	fh.Close()

	label := util.Label(f)

	i := Image{
		Filename: fname,
		Label:    label,
	}

	return &i, nil
}

func RenderFeature(f geojson.Feature, opts *RenderOptions) (go_image.Image, error) {

	img_opts := wof_image.NewDefaultOptions()
	im, err := wof_image.FeatureToImage(f, img_opts)

	if err != nil {
		return nil, err
	}

	final := im

	// TO DO : draw labels at the top of the image rather than bottom

	if opts.Labels {

		label := util.Label(f)

		bounds := im.Bounds()
		max := bounds.Max

		w := max.X
		h := max.Y + 52

		pt_x := 10
		pt_y := max.Y + 32

		im2 := go_image.NewRGBA(go_image.Rect(0, 0, w, h))

		draw.Draw(im2, bounds, im, go_image.ZP, draw.Src)

		col := color.RGBA{0, 0, 0, 255}

		point := fixed.Point26_6{
			fixed.Int26_6(pt_x * 64),
			fixed.Int26_6(pt_y * 64),
		}

		d := &font.Drawer{
			Dst:  im2,
			Src:  go_image.NewUniform(col),
			Face: inconsolata.Bold8x16,
			Dot:  point,
		}

		d.DrawString(label)

		final = im2
	}

	return final, nil
}
