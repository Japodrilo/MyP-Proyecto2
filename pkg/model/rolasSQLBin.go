package model

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _rolasSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x84\x54\xc1\x6e\xa3\x30\x10\xbd\xf3\x15\x73\x4b\x22\x11\xa9\xed\x71\x7b\x62\x2b\x37\x42\xdb\xa5\x95\x43\x57\xdb\x53\xe4\x84\xd9\x14\x95\x18\x64\x5c\xad\xfa\xf7\x95\x6d\x82\x81\x8c\x09\x37\xcf\xf3\xcc\xbc\xf7\xc6\xc3\x7a\x0d\x52\x9c\xf0\x07\x1c\x14\x0a\x8d\x6b\xfd\xd5\x60\xbb\xd6\x62\x5f\x61\xf4\xc0\x59\x92\x33\xc8\x93\x9f\x4f\x0c\x2c\x00\xcb\x08\x00\xa0\x2c\x76\xe6\x08\xee\x4b\xb3\x9c\x6d\x18\x87\x17\x9e\xfe\x4e\xf8\x1b\xfc\x62\x6f\xb1\xbd\x56\x60\x7b\x50\x65\xa3\xcb\x5a\x02\x40\xce\xfe\xe6\xd1\xea\x3e\x8a\xa8\x96\x37\x51\x9a\x6d\x19\xcf\x4d\xb1\xe7\xae\xd7\x9f\xe4\xe9\x95\x6d\x97\x37\xf1\xe2\x05\x55\x5b\xcb\xc5\xea\x9e\xcc\xbd\x0d\xe7\xde\xc6\x8b\x8d\xaa\x3f\x1b\x97\x7a\x91\x79\x17\xce\xbc\x8b\x17\xaf\xf2\x43\xd6\xff\x6d\xdb\x8b\xbe\x0d\xaa\x7f\xb5\x3a\xa1\x22\xbd\xf2\xa8\x37\xac\x8f\xcd\x18\x46\xfa\xea\x20\xd3\x1d\xfa\xcf\x78\xe9\xe2\x8f\xcf\x9c\xa5\x9b\xcc\xd4\x00\x80\x65\x57\x61\x05\x9c\x3d\x32\xce\xb2\x07\xb6\x75\xba\x7a\x84\x1c\x41\x63\x0d\x0e\x69\x31\xd0\x48\x48\x6b\x07\x3a\x23\xa4\xd5\xe2\x88\xbb\x33\x67\xcf\x56\xa1\xa8\xfa\xf0\x20\xbe\x2f\x95\x7e\xdf\x15\x42\x4f\xee\x17\x28\x26\x71\x92\xfe\xd1\xcc\x98\x64\xef\x10\x4f\xde\x9e\xaf\x3c\xdb\x90\xd5\xad\x16\x4a\x13\x24\x51\x16\x7d\x74\x86\xa4\xa8\xf6\x9f\x27\x92\xa4\x43\x3c\x49\x7b\xbe\x42\xb2\x11\xfa\x9d\x22\x19\x22\xff\x85\x42\xf9\x78\x57\x96\xe4\xa9\xea\x4a\x90\x34\x2d\xe0\x59\x9a\xe3\xb5\x3f\x00\xf9\xee\xe3\xb0\xce\x79\x6d\xba\xd4\x15\x52\x71\x25\x0e\x1f\x53\x6d\x33\xb2\x1d\x74\x44\xa9\xa8\x6a\x97\x1b\xd5\x4b\x18\xad\x95\x5f\xf2\xf1\x9d\x50\x15\xab\x76\x54\xc1\xcd\xdd\x63\xe4\x38\x4a\xe9\x1e\x2d\x35\x91\x33\x36\xb7\x9c\x71\xf8\xe9\x3b\x68\x30\x34\xaf\xb7\xad\x65\xdc\x27\x05\x25\xb9\x8b\x53\x57\xcc\xef\x62\x80\x86\x92\x5d\xe9\x61\x6e\xb7\xad\x1e\x34\x86\x7c\x07\x00\x00\xff\xff\x22\xcd\xd7\xda\xa1\x06\x00\x00")

func rolasSqlBytes() ([]byte, error) {
	return bindataRead(
		_rolasSql,
		"rolas.sql",
	)
}

func rolasSql() (*asset, error) {
	bytes, err := rolasSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "rolas.sql", size: 1697, mode: os.FileMode(420), modTime: time.Unix(1539836583, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"rolas.sql": rolasSql,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"rolas.sql": &bintree{rolasSql, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
