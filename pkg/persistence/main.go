package persistence

import (
	"github.com/williamflynt/topolith/pkg/world"
	"os"
	"path/filepath"
	"runtime"
)

// Persistence defines the interface for saving, loading, and managing worlds.
type Persistence interface {
	Save(world world.World) error
	Load(name string) (world.World, error)
	ListWorlds() ([]string, error)
	SetSourcePath(pathOrUrl string)
}

// filePersistence is the unexported struct that implements the Persistence interface.
type filePersistence struct {
	directory string
}

// NewFilePersistence creates a new instance of filePersistence with the appropriate directory based on the OS.
func NewFilePersistence() Persistence {
	dir := getDefaultDirectory()
	return &filePersistence{directory: dir}
}

// getDefaultDirectory returns the default directory based on the OS.
func getDefaultDirectory() string {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("APPDATA"), "topolith")
	case "darwin":
		return filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "topolith")
	default: // Unix-like systems
		return filepath.Join(os.Getenv("HOME"), ".local", "share", "topolith")
	}
}

// Save saves a world to a file.
func (fp *filePersistence) Save(world world.World) error {
	if err := os.MkdirAll(fp.directory, 0755); err != nil {
		return err
	}

	filePath := filepath.Join(fp.directory, world.Name()+".world")
	data := world.String()
	return os.WriteFile(filePath, []byte(data), 0644)
}

// Load loads a world from a file.
func (fp *filePersistence) Load(name string) (world.World, error) {
	filePath := name
	if filepath.Ext(name) != ".world" {
		filePath = filepath.Join(fp.directory, name+".world")
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	println(string(data))
	return world.FromString(string(data))
}

// ListWorlds scans the directory for world files and returns their names.
func (fp *filePersistence) ListWorlds() ([]string, error) {
	files, err := os.ReadDir(fp.directory)
	if err != nil {
		return nil, err
	}

	var worlds []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".world" {
			worlds = append(worlds, file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))])
		}
	}
	return worlds, nil
}

// SetSourcePath allows setting a custom persistence layer location at runtime.
func (fp *filePersistence) SetSourcePath(dir string) {
	fp.directory = dir
}
