package service

import (
	"image"
	"image/jpeg"
	"log"
	"os"

	"github.com/nfnt/resize"
)

type Service interface {
	ResizeImage(img *image.Image, cfg ResizeConfig) (*image.Image, error)
	ResizeImgFromFile(cfg ResizeInFileConfig) (*image.Image, error)
}

type Options struct {
}

func New(opts Options) Service {
	return &service{&opts}
}

type service struct {
	opts *Options
}

// InterpolationFunction constants
const (
	// Nearest-neighbor interpolation
	NearestNeighbor resize.InterpolationFunction = iota
	// Bilinear interpolation
	Bilinear
	// Bicubic interpolation (with cubic hermite spline)
	Bicubic
	// Mitchell-Netravali interpolation
	MitchellNetravali
	// Lanczos interpolation (a=2)
	Lanczos2
	// Lanczos interpolation (a=3)
	Lanczos3
)

type ResizeConfig struct {
	W, H   uint
	Interp resize.InterpolationFunction
}

func (svc *service) ResizeImage(img *image.Image, cfg ResizeConfig) (*image.Image, error) {
	resizedImg := resize.Resize(cfg.W, cfg.H, *img, cfg.Interp)

	log.Println("Image resized successfully")

	return &resizedImg, nil
}

// ---------------------------

type ResizeInFileConfig struct {
	In, Out string
	W, H    uint
	Interp  resize.InterpolationFunction
}

func (svc *service) ResizeImgFromFile(cfg ResizeInFileConfig) (*image.Image, error) {
	// read image
	file, err := os.Open(cfg.In)
	log.Println("Openning file ", cfg.In)
	if err != nil {
		log.Fatalln(err)
	}

	// decode file into image
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	// resize img
	resizedImg := resize.Resize(cfg.W, cfg.H, img, cfg.Interp)

	// create file
	newFile, err := os.Create(cfg.Out)
	if err != nil {
		log.Fatalln(err)
	}
	defer newFile.Close()

	// write new img to a file
	jpeg.Encode(newFile, resizedImg, &jpeg.Options{Quality: jpeg.DefaultQuality})

	log.Println("Image resized: ", cfg.Out)

	return &resizedImg, nil

}
