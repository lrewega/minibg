// Command minibg generates a tiny image and aggressively sets it as the
// wallpaper for all displays on OS X. It can also set the wallpaper to
// a specific image file.

package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/fogleman/gg"
)

// randomImage produces a PNG image with random pixel colors and writes it
// to the given outputFilename.
func randomImage(outputFilename string) error {
	rng := rand.New(rand.NewSource(time.Now().Unix()))
	width, height := 2+rng.Intn(2), 2+rng.Intn(2)
	dc := gg.NewContext(width, height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dc.SetRGB(rng.Float64(), rng.Float64(), rng.Float64())
			dc.SetPixel(x, y)
		}
	}
	return dc.SavePNG(outputFilename)
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "usage: %s [-o <output_image> | <input_image>]\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	filename := os.ExpandEnv(fmt.Sprintf("$HOME/.%s.png", path.Base(os.Args[0])))
	flag.StringVar(&filename, "out", filename, "output filename")
	flag.Usage = usage
	flag.Parse()
	switch len(flag.Args()) {
	case 0:
		// no arguments
		if err := randomImage(filename); err != nil {
			fmt.Fprintf(os.Stderr, "error writing image to %q: %v\n", filename, err)
			os.Exit(2)
		}
	case 1:
		filename := flag.Arg(0)
		info, err := os.Stat(filename)
		if err != nil && !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "error accessing file %q: %v\n", filename, err)
			os.Exit(3)
		} else if !info.Mode().IsRegular() {
			fmt.Fprintf(os.Stderr, "file %q can't be used as it's not a regular file.\n", filename)
			os.Exit(4)
		}
	default:
		flag.Usage()
		os.Exit(1)
	}

	if err := SetWallpaper(filename); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(5)
	}
}
