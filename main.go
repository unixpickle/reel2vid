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
	var frameRepeat int
	var loops int
	flag.IntVar(&width, "width", -1, "width of each frame")
	flag.IntVar(&height, "height", -1, "height of each frame")
	flag.Float64Var(&fps, "fps", 12.0, "frame rate of exported video")
	flag.IntVar(&frameRepeat, "frame-repeat", 1, "number of times to repeat each frame, for lower FPS")
	flag.IntVar(&loops, "loops", 1, "number of times to repeat the whole video")
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
	} else if width == -1 && height == -1 {
		essentials.Die("Provide -width or -height. See -help.")
	}
	inputPath := flag.Args()[0]
	outputPath := flag.Args()[1]

	f, err := os.Open(inputPath)
	essentials.Must(err)
	img, _, err := image.Decode(f)
	f.Close()
	essentials.Must(err)

	if height == -1 {
		height = img.Bounds().Dy()
	}
	if width == -1 {
		width = img.Bounds().Dx()
	}

	if img.Bounds().Dy()%height != 0 {
		essentials.Die("height does not divide file heighr", img.Bounds().Dy())
	} else if img.Bounds().Dx()%width != 0 {
		essentials.Die("width does not divide file width", img.Bounds().Dx())
	}

	writer, err := ffmpego.NewVideoWriter(outputPath, width, height, fps)
	essentials.Must(err)
	defer func() {
		essentials.Must(writer.Close())
	}()
	for i := 0; i < loops; i++ {
		for y := 0; y < img.Bounds().Dy(); y += height {
			for x := 0; x < img.Bounds().Dx(); x += width {
				crop := cropImage(img, x, y, width, height)
				for j := 0; j < frameRepeat; j++ {
					essentials.Must(writer.WriteFrame(crop))
				}
			}
		}
	}
}

func cropImage(img image.Image, sx, sy, width, height int) *image.RGBA {
	res := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			res.Set(x, y, img.At(x+sx+img.Bounds().Min.X, y+sy+img.Bounds().Min.Y))
		}
	}
	return res
}
