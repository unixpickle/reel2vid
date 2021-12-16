package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/unixpickle/essentials"
	"github.com/unixpickle/ffmpego"
)

func main() {
	var width int
	var height int
	var fps float64
	flag.IntVar(&width, "width", -1, "width of each frame")
	flag.IntVar(&width, "height", -1, "height of each frame")
	flag.Float64Var(&fps, "fps", 12.0, "frame rate of exported video")
	flag.Usage = func() {
		fmt.Fprintln(
			os.Stderr,
			"Usage: reel2vid [flags] -width X -height Y <input_image> <output_video>",
		)
		fmt.Fprintln(os.Stderr)
		flag.PrintDefaults()
		os.Exit(1)
	}
	flag.Parse()

	if len(flag.Args()) != 2 {
		flag.Usage()
	} else if width == -1 || height == -1 {
		essentials.Die("Provide -width and -height. See -help.")
	}
	inputPath := flag.Args()[0]
	outputPath := flag.Args()[1]

	f, err := os.Open(inputPath)
	essentials.Must(err)
	img, _, err := image.Decode(f)
	f.Close()
	essentials.Must(err)

	if img.Bounds().Dy() != height {
		essentials.Die("height does not match file, expected", img.Bounds().Dy())
	} else if img.Bounds().Dx()%width != 0 {
		essentials.Die("width does not divide file width", img.Bounds().Dx())
	}

	writer, err := ffmpego.NewVideoWriter(outputPath, width, height, fps)
	essentials.Must(err)
	defer func() {
		essentials.Must(writer.Close())
	}()
	for x := 0; x < img.Bounds().Dx(); x += width {
		crop := cropImage(img, x, 0, width, height)
		essentials.Must(writer.WriteFrame(crop))
	}
}

func cropImage(img image.Image, x, y, width, height int) *image.RGBA {
	res := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			res.Set(x, y, img.At(x+width+img.Bounds().Min.X, y+height+img.Bounds().Min.Y))
		}
	}
	return res
}
