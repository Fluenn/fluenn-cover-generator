package generate

import (
	"fmt"
	"image"

	_ "image/gif"
	"image/jpeg"

	_ "image/jpeg"

	_ "image/png"
	"os"

	"github.com/oliamb/cutter"
	"github.com/valyala/fastrand"
)

func getImageSize(imgPath string) *image.Config {

	fmt.Println(imgPath)
	file, err := os.Open(imgPath)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &image

}

func getCropCoordinate(originalWidth, originalHeight, cropWidth, cropHeight uint32) (uint, uint) {
	var minX uint32 = 0
	var minY uint32 = 0
	var maxX uint32 = originalWidth - cropWidth
	var maxY uint32 = originalHeight - cropHeight

	x := fastrand.Uint32n(maxX-minX) + minX
	y := fastrand.Uint32n(maxY-minY) + minY

	return uint(x), uint(y)
}

func cropImage(path string, width uint, height uint, x uint, y uint) (image.Image, error) {

	fmt.Println(path, width, height, x, y)

	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	croppedImg, err := cutter.Crop(img, cutter.Config{
		Width:  int(width),
		Height: int(height),
		Anchor: image.Point{int(x), int(y)},
		Mode:   cutter.TopLeft, // optional, default value
	})

	return croppedImg, err

}

func SaveJPEGImage(img image.Image, path string) {

	f, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	opt := jpeg.Options{
		Quality: 90,
	}

	err = jpeg.Encode(f, img, &opt)

	if err != nil {
		fmt.Println(err)
		return
	}
}

type CropConfig struct {
	Width  uint
	Height uint
}

func GenerateFrom(imgPath string, cropConfig *CropConfig) image.Image {

	imageConfig := getImageSize(imgPath)

	var cropWidth uint = cropConfig.Width
	var cropHeight uint = cropConfig.Height
	originalWidth := imageConfig.Width
	originalHeight := imageConfig.Height

	x, y := getCropCoordinate(
		uint32(originalWidth),
		uint32(originalHeight),
		uint32(cropWidth),
		uint32(cropHeight),
	)

	img, err := cropImage(
		imgPath,
		cropWidth,
		cropHeight,
		x,
		y,
	)

	if err != nil {
		panic(err)
	}

	return img
}
