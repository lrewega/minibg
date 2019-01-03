package main

import (
	"os/exec"

	"go.uber.org/multierr"
)

func SetWallpaper(filename string) error {
	return multierr.Combine(
		exec.Command("termux-wallpaper", "-f", filename).Run(),
		exec.Command("termux-wallpaper", "-l", "-f", filename).Run(),
	)
}
