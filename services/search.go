package services

// #cgo CFLAGS: -I /usr/include/x86_64-linux-gnu
// #cgo LDFLAGS: -L /usr/lib/x86_64-linux-gnu -lavformat -lavcodec -lavutil
//
// #include "stdafx.h"
import "C"

import (
	"encoding/json"
	"errors"
	"io/fs"
	"path"
	"path/filepath"
	"unsafe"

	"github.com/watsonserve/scaner/entities"
	"github.com/watsonserve/scaner/utils"
)

const BUFSIZ = 2048

const UNKNOW = 0
const PICTURE = 1
const AUDIO = 2
const VIEDO = 3

type EnMediaType int

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

func loadAudioMeta(file string, mem unsafe.Pointer) (*entities.AudioMeta, error) {
	c_filename := C.CString(file)
	cret := C.load_audio(c_filename, mem, C.size_t(BUFSIZ))
	C.free(unsafe.Pointer(c_filename))

	if int(cret) < 0 {
		return nil, errors.New("load audio failed")
	}

	meta := &entities.AudioMeta{}
	err := json.Unmarshal(C.GoBytes(mem, cret), meta)
	if nil == err {
		if "" == meta.Title {
			meta.Title = utils.FileBaseName(file)
		}
		meta.Url = file
		meta.Hash, err = utils.Sha1File(file)
	}

	return meta, err
}

type FileError struct {
	filename string
	err      error
}

func search(root string) ([]*entities.AudioMeta, []*FileError, error) {
	mem := C.malloc(C.size_t(BUFSIZ))
	if nil == mem {
		return nil, nil, errors.New("no memary")
	}
	defer C.free(mem)

	audioList := make([]*entities.AudioMeta, 0)
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
