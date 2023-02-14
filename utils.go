package main

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

func mediaType(filename string) EnMediaType {
	switch path.Ext(filename) {
	case ".mp3":
		fallthrough
	case ".wav":
		fallthrough
	case ".flac":
		fallthrough
	case ".ape":
		fallthrough
	case ".wma":
		fallthrough
	case ".aac":
		fallthrough
	case ".aiff":
		return AUDIO
	default:
	}
	return UNKNOW
}

func WalkDir(root string) ([]string, error) {
	files := make([]string, 0)
	fn := func(filename string, info fs.DirEntry, err error) error {
		if nil != err {
			return err
		}

		// if ".git" == info.Name() {
		// 	return filepath.SkipDir
		// }

		if !info.IsDir() && AUDIO == mediaType(filename) {
			files = append(files, filename)
		}

		return nil
	}

	err := filepath.WalkDir(root, fn)

	return files, err
}

func sha1File(filePath string) (string, error) {
	file, err := os.Open(filePath)
	hash := sha1.New()
	hashVal := ""

	if nil == err {
		defer file.Close()
		_, err = io.Copy(hash, file)
	}

	if nil == err {
		hashVal = hex.EncodeToString(hash.Sum(nil))
	}

	return hashVal, err
}

func fileBaseName(name string) string {
	name = path.Base(path.Clean(name))

	for i := len(name) - 1; i >= 0 && name[i] != '/'; i-- {
		if name[i] == '.' {
			name = name[:i]
			break
		}
	}

	return name
}
