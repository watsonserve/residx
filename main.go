package main

// #cgo LDFLAGS: -lavformat -lavcodec -lavutil
//
// #include "stdafx.h"
import "C"

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"unsafe"
)

const BUFSIZ = 2048

const UNKNOW = 0
const PICTURE = 1
const AUDIO = 2
const VIEDO = 3

type EnMediaType int

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
	if nil == err {
		if "" == meta.Title {
			meta.Title = fileBaseName(file)
		}
		meta.Url = file
		meta.Hash, err = sha1File(file)
	}

	return meta, err
}

type FileError struct {
	filename string
	err      error
}

func search(root string) ([]*AudioMeta, []*FileError, error) {
	mem := C.malloc(C.size_t(BUFSIZ))
	if nil == mem {
		return nil, nil, errors.New("no memary")
	}
	defer C.free(mem)

	audioList := make([]*AudioMeta, 0)
	errList := make([]*FileError, 0)

	err := filepath.WalkDir(root, func(filename string, info fs.DirEntry, err error) error {
		if nil != err {
			errList = append(errList, &FileError{filename, err})
			return filepath.SkipDir
		}

		if info.IsDir() || AUDIO != mediaType(filename) {
			return nil
		}

		meta, err := loadAudioMeta(filename, mem)
		if nil == err {
			audioList = append(audioList, meta)
		} else {
			errList = append(errList, &FileError{filename, err})
		}

		return nil
	})

	if 0 == len(errList) {
		errList = nil
	}

	return audioList, errList, err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, os.Args[0]+" dirpath")
		return
	}

	root := os.Args[1]
	audioList, errList, err := search(root)
	if nil != err {
	}
	if nil != errList {

	}

	var db *sql.DB
	result, err := db.Exec("INSERR INTO ? VALUES (), ()?", id)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	if rows != 1 {
		log.Fatalf("expected to affect 1 row, affected %d", rows)
	}
}
