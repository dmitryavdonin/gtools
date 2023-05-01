package imageResizer

import (
	"bytes"
	"image"
	"image/jpeg"

	"github.com/disintegration/imaging"
)

type imageResizer struct{}

func NewImageResizer() *imageResizer {
	return &imageResizer{}
}

func (i *imageResizer) ResizeImages(params *ResizeImageParams) (result *ResizeImageResult, err error) {
	images := make([]*ResizedImage, len(params.RequiredSizes))
	for index, size := range params.RequiredSizes {
		bytes, err := i.ResizeImage(params.Image, size.Width, size.Height)
		if err != nil {
			return nil, err
		}
		images[index] = &ResizedImage{
			Key:   size.Key,
			Bytes: bytes,
		}
	}
	return &ResizeImageResult{Images: images}, nil
}

func (i *imageResizer) ResizeImage(imgBytes []byte, width int, height int) (bytesR []byte, err error) {
	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return
	}
	var converted *image.NRGBA
	if height != 0 {
		converted = imaging.Fill(img, width, height, imaging.Center, imaging.Lanczos)
	} else {
		converted = imaging.Resize(img, width, 0, imaging.Lanczos)
	}

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, converted, nil)
	if err != nil {
		return nil, err
	}
	bytesR = buf.Bytes()
	return
}
