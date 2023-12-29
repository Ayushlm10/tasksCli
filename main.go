package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	gap "github.com/muesli/go-app-paths"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Setup directory to store the tasks
func setupXDGPath() string {
	scope := gap.NewScope(gap.User, "tasks")
	dirs, err := scope.DataDirs()
	if err != nil {
		log.Fatal(err)
	}
	var tasksDir string
	if len(dirs) > 0 {
		tasksDir = dirs[0]
	} else {
		tasksDir, err = os.UserHomeDir()
		if err != nil {
			log.Fatalf("Couldn't get user home directory %s", err)
		}
	}

	if err := initTaskDir(tasksDir); err != nil {
		log.Fatal(err)
	}
	return tasksDir
}

// initialize the directory if it does not exist
func initTaskDir(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return os.Mkdir(path, 0o770)
		}
		return err
	}
	return nil
}

func openDb(path string) (*taskDb, error) {
	db, err := sql.Open("sqlite3", filepath.Join(path, "tasks.db"))
	if err != nil {
		log.Fatal(err)
	}
	t := taskDb{db, path}
	if !t.isTableExists("tasks") {
		err := t.createTable()
		if err != nil {
			return nil, err
		}
	}
	return &t, nil
}
