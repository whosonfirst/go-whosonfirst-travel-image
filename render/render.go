package render

type Image struct {
	Filename string
	Label    string
}

type RenderOptions struct {
	Labels bool
	Root   string
	Prefix string
}
