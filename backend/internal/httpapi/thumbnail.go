package httpapi

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"io"

	_ "image/gif"
	_ "image/png"

	xdraw "golang.org/x/image/draw"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/webp"
)

const thumbMaxSide = 480

func makeThumbnail(r io.Reader) ([]byte, error) {
	src, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()
	if w <= thumbMaxSide && h <= thumbMaxSide {
		return encodeJpeg(src)
	}
	var tw, th int
	if w >= h {
		tw = thumbMaxSide
		th = h * thumbMaxSide / w
	} else {
		th = thumbMaxSide
		tw = w * thumbMaxSide / h
	}
	if tw < 1 {
		tw = 1
	}
	if th < 1 {
		th = 1
	}
	dst := image.NewRGBA(image.Rect(0, 0, tw, th))
	xdraw.Draw(dst, dst.Bounds(), image.NewUniform(color.White), image.Point{}, xdraw.Src)
	xdraw.CatmullRom.Scale(dst, dst.Bounds(), src, b, xdraw.Over, nil)
	return encodeJpeg(dst)
}

func encodeJpeg(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	bg := image.NewRGBA(img.Bounds())
	xdraw.Draw(bg, bg.Bounds(), image.NewUniform(color.White), image.Point{}, xdraw.Src)
	xdraw.Draw(bg, bg.Bounds(), img, img.Bounds().Min, xdraw.Over)
	if err := jpeg.Encode(&buf, bg, &jpeg.Options{Quality: 82}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func thumbKey(fileKey string) string {
	return fileKey + ".thumb.jpg"
}

func isRasterImage(mime string) bool {
	switch mime {
	case "image/png", "image/jpeg", "image/gif", "image/webp", "image/bmp":
		return true
	}
	return false
}
