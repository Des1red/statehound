package collector

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"statehound/internal/model"
)

func collectWatchedFiles() (map[string]FileWatch, error) {
	files := make(map[string]FileWatch)

	paths := watchedFilePaths()

	for _, path := range paths {
		item := collectFile(path)
		files[path] = item
	}

	return files, nil
}

func watchedFilePaths() []string {
	var paths []string

	paths = append(paths, model.SystemCronFile)

	paths = appendGlob(paths, model.SystemCronDir+"/*")

	paths = appendGlob(paths, model.SystemdSystemDir+"/*.service")
	paths = appendGlob(paths, model.SystemdSystemDir+"/*.timer")
	paths = appendGlob(paths, model.SystemdSystemDir+"/*.socket")
	paths = appendGlob(paths, model.SystemdSystemDir+"/*.d/*.conf")

	paths = appendGlob(paths, model.SystemdUserGlobalDir+"/*.service")
	paths = appendGlob(paths, model.HomeUserSystemdGlob)

	paths = appendGlob(paths, model.XDGAutostartDir+"/*.desktop")
	paths = appendGlob(paths, model.HomeAutostartGlob)

	paths = appendGlob(paths, model.RootAuthorizedKeysGlob)
	paths = appendGlob(paths, model.HomeAuthorizedKeysGlob)

	paths = append(paths, model.RootBashrc, model.RootProfile)
	paths = appendGlob(paths, model.HomeBashrcGlob)
	paths = appendGlob(paths, model.HomeProfileGlob)
	paths = appendGlob(paths, model.HomeZshrcGlob)
	return paths
}

func appendGlob(paths []string, pattern string) []string {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return paths
	}

	return append(paths, matches...)
}

func collectFile(path string) FileWatch {
	info, err := os.Stat(path)
	if err != nil {
		return FileWatch{
			Path:   path,
			Exists: false,
		}
	}

	item := FileWatch{
		Path:    path,
		Exists:  true,
		Size:    info.Size(),
		Mode:    info.Mode().String(),
		ModTime: info.ModTime(),
	}

	if !info.IsDir() {
		item.Hash = hashFile(path)
	}

	return item
}

func hashFile(path string) string {
	file, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer file.Close()

	h := sha256.New()

	if _, err := io.Copy(h, file); err != nil {
		return ""
	}

	return hex.EncodeToString(h.Sum(nil))
}
