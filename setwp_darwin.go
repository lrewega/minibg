package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"

	_ "github.com/mattn/go-sqlite3"
)

func SetWallpaper(filename string) error {
	if err := updateWallpaper(filename); err != nil {
		return fmt.Errorf("failed to set wallpaper: %v", err)
	}

	if err := restartDock(); err != nil {
		return fmt.Errorf("failed to restart dock: %v", err)
	}

	return nil
}

func updateWallpaper(filename string) error {
	const dbPath = "$HOME/Library/Application Support/Dock/desktoppicture.db"
	db, err := sql.Open("sqlite3", os.ExpandEnv(dbPath))
	if err != nil {
		return err
	}
	defer db.Close()

	const statement = `
		-- remove all image:space,display relations.
		delete from preferences;
		-- remove all registered images;
		delete from data;
		-- register the image;
		insert into data values (?);
		-- set the image for every space and display combination.
		insert into preferences
			select
				1, -- magical key to indicate this is a Wallpaper preference.
				data.ROWID, -- the row containing the image in data.
				pictures.ROWID -- the row referencing the space+display combo.
			from pictures
			inner join data
				on data.value = ?;
	`
	_, err = db.Exec(statement, filename, filename)
	return err
}

// restartDock kills the Dock process, which restarts immediately, and updates
// the wallpaper at start.
func restartDock() error {
	return exec.Command("killall", "Dock").Run()
}
