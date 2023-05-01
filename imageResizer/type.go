package imageResizer

type ResizeImageParams struct {
	Image         []byte
	RequiredSizes []ImageSize
}

type ImageSize struct {
	Key    string
	Width  int
	Height int
}

type ResizeImageResult struct {
	Images []*ResizedImage
}

type ResizedImage struct {
	Key   string
	Bytes []byte
}
