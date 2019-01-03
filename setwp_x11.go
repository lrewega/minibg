// +build !android,linux

package main

import "os/exec"

func SetWallpaper(filename string) error {
	return exec.Command("feh", "--bg-scale", filename).Run()
}
