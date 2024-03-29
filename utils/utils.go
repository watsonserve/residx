package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
	"path"

	"github.com/watsonserve/goutils"
	"go.mongodb.org/mongo-driver/bson"
)

func GetOption() (map[string][]string, error) {
	options := []goutils.Option{
		{
			Opt:       'h',
			Option:    "help",
			HasParams: false,
		},
		{
			Opt:       'c',
			Option:    "config",
			HasParams: true,
		},
	}
	argMap, params := goutils.GetOptions(options)
	cfg, err := goutils.GetConf(argMap["config"])
	if nil == err {
		cfg["listen"] = params
	}
	return cfg, err
}

func Sha1File(filePath string) (string, error) {
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

func FileBaseName(name string) string {
	name = path.Base(path.Clean(name))

	for i := len(name) - 1; i >= 0 && name[i] != '/'; i-- {
		if name[i] == '.' {
			name = name[:i]
			break
		}
	}

	return name
}

func ToMap(val interface{}) (bson.M, error) {
	bs, err := bson.Marshal(val)
	if nil != err {
		return nil, err
	}
	result := bson.M{}
	err = bson.Unmarshal(bs, result)
	return result, err
}

func MapToKvList(m bson.M) bson.D {
	result := make([]bson.E, 0)
	for k, v := range m {
		result = append(result, bson.E{Key: k, Value: v})
	}
	return result
}

func MapFilter(m bson.M, ks []string) bson.M {
	result := bson.M{}
	for _, k := range ks {
		v, has := m[k]
		if has {
			result[k] = v
		}
	}
	return result
}
