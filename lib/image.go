package lib

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"os"

	"golang.org/x/image/draw"
)

//ImageSizer make image thumbnail
func ImageSizer(size int, src io.Reader) (io.Reader, error) {
	var w = new(bytes.Buffer)

	imgSrc, _, err := image.Decode(src)
	if err != nil {
		return w, err
	}

	rctSrc := imgSrc.Bounds()

	var imgDst *image.RGBA
	switch {
	case rctSrc.Dx() <= rctSrc.Dy():
		imgDst = image.NewRGBA(image.Rect(0, 0, int(float64(size)*float64(rctSrc.Dx())/float64(rctSrc.Dy())), size))
	case rctSrc.Dx() > rctSrc.Dy():
		imgDst = image.NewRGBA(image.Rect(0, 0, size, int(float64(size)*float64(rctSrc.Dy())/float64(rctSrc.Dx()))))
	}

	draw.CatmullRom.Scale(imgDst, imgDst.Bounds(), imgSrc, imgSrc.Bounds(), draw.Over, nil)

	if err := png.Encode(w, imgDst); err != nil {
		return w, err
	}

	return w, nil
}

//ImageConverter converts image to png and export it to the file path
func ImageConverter(src io.Reader, filepath string) error {
	imgSrc, _, err := image.Decode(src)
	if err != nil {
		return err
	}

	rctSrc := imgSrc.Bounds()

	var imgDst *image.RGBA
	imgDst = image.NewRGBA(image.Rect(0, 0, rctSrc.Dx(), rctSrc.Dy()))

	draw.CatmullRom.Scale(imgDst, imgDst.Bounds(), imgSrc, imgSrc.Bounds(), draw.Over, nil)

	dst, err := os.Create(filepath)
	if err != nil {
		return err
	}

	if err := png.Encode(dst, imgDst); err != nil {
		dst.Close()
		return err
	}

	dst.Close()

	if err = os.Chmod(filepath, 0777); err != nil {
		return err
	}

	return nil
}
