package utils

import (
	"crypto/md5"
	"encoding/hex"
	"path/filepath"
	"runtime"
)

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func GetProfileImageFromUsername(username string) string {
	// Get the path to the current file
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "default_profile.png"
	}

	// Get the path to the images folder
	dir := filepath.Join(filepath.Dir(currentFile), "..", "..", "..", "images")

	// Check if there is a matching file
	matches, err := filepath.Glob(filepath.Join(dir, username+".*"))
	if err != nil || len(matches) < 1 {
		return "default_profile.png"
	}
	return filepath.Base(matches[0])
}
