package main

// #cgo LDFLAGS: -lavformat -lavcodec -lavutil
//
// #include "stdafx.h"
import "C"

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"unsafe"
)

type EnMediaType int

const UNKNOW = 0
const PICTURE = 1
const AUDIO = 2
const VIEDO = 3

const BUFSIZ = 2048

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

type AudioMeta struct {
	Url        string `json:"url"`
	Hash       string `json:"hash"`
	Title      string `json:"title"`
	Album      string `json:"album"`
	Artist     string `json:"artist"`
	SampleRate int64  `json:"sample_rate"`
	BitRate    int64  `json:"bit_rate"`
	Channels   int64  `json:"channels"`
	Duration   int64  `json:"duration"`
}

func loadAudioMeta(file string, mem unsafe.Pointer) (*AudioMeta, error) {
	c_filename := C.CString(file)
	cret := C.load_audio(c_filename, mem, C.size_t(BUFSIZ))
	C.free(unsafe.Pointer(c_filename))

	if int(cret) < 0 {
		return nil, errors.New("load audio failed")
	}

	meta := &AudioMeta{}
	err := json.Unmarshal(C.GoBytes(mem, cret), meta)
	if nil != err {
		return nil, err
	}

	return meta, nil
}

func SHA1(src string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(src)))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, os.Args[0]+" dirpath")
		return
	}

	root := os.Args[1]
	mem := C.malloc(C.size_t(BUFSIZ))

	files, err := WalkDir(root)
	if nil != err {
		fmt.Fprint(os.Stderr, err.Error())
		return
	}

	for _, file := range files {
		meta, err := loadAudioMeta(file, mem)

		if nil != err {
			continue
		}

		if "" == meta.Title {
			meta.Title = file
		}

		fmt.Println(meta)
	}

	C.free(mem)
}
