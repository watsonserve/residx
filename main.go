package main

// #cgo LDFLAGS: -lavformat -lavcodec -lavutil
//
// #include "stdafx.h"
import "C"

import (
	"encoding/json"
	"errors"
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

func search(root string) ([]*AudioMeta, error) {
	mem := C.malloc(C.size_t(BUFSIZ))

	files, err := WalkDir(root)
	if nil != err {
		return nil, err
	}

	audioList := make([]*AudioMeta, 0)
	for _, file := range files {
		meta, err := loadAudioMeta(file, mem)

		if nil != err {
			continue
		}

		audioList = append(audioList, meta)
	}

	C.free(mem)

	return audioList, nil
}

// root := os.Args[1]
// if len(os.Args) < 2 {
// 	fmt.Fprint(os.Stderr, os.Args[0]+" dirpath")
// 	return
// }
