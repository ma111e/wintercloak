package main

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	"path/filepath"
)

// FileItem represents a file or folder with its name, type, and nested items.
type FileItem struct {
	Name  string     `json:"name"`
	Type  string     `json:"type"`  // "file" or "folder"
	Items []FileItem `json:"items"` // Nested files and folders (only for folders)
}

// FolderSelectionResult holds the selected path and its file structure.
type FolderSelectionResult struct {
	Path  string     `json:"path"`
	Items []FileItem `json:"items"`
}

// getFolderContents recursively retrieves the contents of a directory.
func getFolderContents(dirPath string) ([]FileItem, error) {
	var items []FileItem

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		item := FileItem{
			Name: entry.Name(),
			Type: "file",
		}

		// If it's a folder, recursively fetch its contents
		if entry.IsDir() {
			item.Type = "folder"
			subItems, err := getFolderContents(filepath.Join(dirPath, entry.Name()))
			if err != nil {
				log.Println("Error reading folder:", err)
				continue
			}
			item.Items = subItems
		}

		items = append(items, item)
	}

	return items, nil
}

// SelectFolder opens a dialog to select a folder and retrieves its contents.
func (a *App) SelectFolder() (FolderSelectionResult, error) {
	dialogOptions := runtime.OpenDialogOptions{
		DefaultDirectory:     ".",
		Title:                "Select a folder",
		ShowHiddenFiles:      true,
		CanCreateDirectories: false,
	}

	// Open directory selection dialog
	selectedPath, err := runtime.OpenDirectoryDialog(a.ctx, dialogOptions)
	if err != nil {
		log.Println("Error opening folder dialog:", err)
		return FolderSelectionResult{}, err
	}

	// Retrieve full folder contents
	items, err := getFolderContents(selectedPath)
	if err != nil {
		log.Println("Error reading folder contents:", err)
		return FolderSelectionResult{Path: selectedPath}, err
	}

	return FolderSelectionResult{Path: selectedPath, Items: items}, nil
}

// SelectFile opens a dialog to select a folder and retrieves its contents.
func (a *App) SelectFile() (string, error) {
	dialogOptions := runtime.OpenDialogOptions{
		DefaultDirectory:     ".",
		Title:                "Select a file",
		ShowHiddenFiles:      true,
		CanCreateDirectories: true,
	}

	// Open directory selection dialog
	selectedPath, err := runtime.OpenFileDialog(a.ctx, dialogOptions)
	if err != nil {
		log.Println("Error opening folder dialog:", err)
		return "", err
	}

	return selectedPath, nil
}

// ReadFile reads the content of a file at the given path
func (a *App) ReadFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	if len(content) > 500000 {
		return "", errors.New("File size exceeds the 500KB limit")
	}
	return string(content), nil
}
