// Code generated by go-bindata.
// sources:
// static/js/main.js
// DO NOT EDIT!

package static

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

var _staticJsMainJs = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x56\x4b\x6b\x2c\x45\x14\xde\x0f\xcc\x7f\x38\x96\xc3\xd0\x4d\x92\x6e\x35\xe8\x22\xf3\x00\x65\x22\xba\x70\x93\x19\xc8\x42\x5c\x74\xba\x4f\xba\x9b\x74\x57\x35\x5d\x95\xe9\xc4\x30\x10\x90\xe8\xc2\x60\x76\x2a\xba\x92\xf8\x42\xc1\xc7\x42\x82\xe2\x25\xbf\xa6\x73\x27\xab\xfb\x17\x2e\xd5\xef\x9a\xcc\x84\xb9\xc9\x26\xd5\x55\xe7\x7c\xdf\x77\x9e\xcc\x58\xc4\x3e\x75\x8d\x28\x66\x82\x89\xd3\x08\x8d\x18\xa3\xc0\xb2\xf1\xdd\x20\x80\x01\x1c\x1e\x53\x5b\xf8\x8c\x82\xc6\xdf\xdc\x04\xfe\x96\x0e\x67\xed\x16\x00\x40\x8c\xe2\x38\xa6\x20\x3c\x9f\x97\x1e\x1a\xc5\x04\xf6\xd0\xdd\x3d\x89\x32\x6b\xe2\x86\x44\xcf\x9c\x7a\xed\xd6\xac\xdd\x6a\xb7\x3a\x5a\x8d\x57\x21\x4d\xad\x18\xdc\x18\xa3\x8f\x90\x73\xcb\x45\x18\x00\x21\xbd\xfa\xa9\x13\xe6\xf7\x1c\x06\xd0\xd1\xc8\xeb\xe5\x27\xd1\x1b\x46\x56\x14\x21\x75\x6a\x84\x9a\x26\xe4\x6e\xc5\x54\x5a\x0b\x0f\x47\xfe\x54\xe2\x95\x60\x1f\xbf\xf1\x49\x4f\xb5\x71\xd8\xd8\x8e\x59\x96\x83\xdc\xdc\xe0\xd9\xf7\x84\x45\x30\x58\xb8\xfb\x00\x7d\xd7\x13\xb0\x55\xde\xda\x81\x8f\x54\xe4\xb7\x0d\xd8\x4e\x26\xc6\xc8\xb5\x4e\x98\x56\xb1\xeb\x0d\x23\xff\x10\xb4\x92\x5b\x11\x2e\xff\x1e\x2a\x79\x92\x90\x59\x7e\x9c\x35\x12\x98\xc0\x00\x64\xf9\xf6\xf1\x60\xcc\xec\x23\x14\x1a\x49\xf8\x8e\x69\x12\xd8\x80\xc4\xa7\x0e\x4b\x8c\x80\xd9\x96\xcc\xa9\xe1\x31\x2e\x60\x03\x88\x99\x70\xf9\xbc\x87\x21\x13\x08\x1b\xf0\xfe\x87\xa3\x32\x92\xc4\x60\x94\x45\x48\x95\x52\x28\xe1\x28\x15\xd3\x88\xf0\x10\x02\xe6\x72\x90\x5e\xaf\x55\xb5\x95\x12\x6b\x44\x3b\x60\x1c\xd7\x86\xec\x3b\xfe\x74\xd8\xb7\x91\x0a\x8c\x87\x7d\x6f\x7b\x38\xf2\xb9\xcd\x28\x45\x5b\xa0\xd3\x37\xbd\xed\x61\xdf\x2c\x5f\x4d\x69\xab\xb0\x96\x9c\xe1\x1a\x3d\x25\x6b\xa6\x74\xb0\x6c\xe1\x07\xc5\x5b\x26\x4f\xe6\x2f\xe4\xae\xe1\x58\xc2\x92\x29\x5d\x10\x92\x89\x01\x0c\x38\xbe\x2a\x58\x63\x8a\x9b\xd2\x36\x81\xf4\x63\x74\x32\xd3\xa6\xe2\x8c\x3a\x7b\xd0\x57\xc9\x58\x2c\x88\x9c\x45\x09\xf1\x9e\xa0\x44\x97\x9d\x66\x1f\x2d\x9b\xee\xb2\xc1\x8a\x3c\x4e\x4e\x44\x31\xc7\xd2\x77\x72\x22\x14\x12\x75\x0d\xd4\x2e\xc6\xd4\x0a\x34\xdd\x10\x2c\x5f\x56\x5a\xd3\x27\x31\x38\x52\xa7\x19\x64\x55\xc6\xf2\x20\xe9\xec\x00\x2d\xba\x86\xd6\x6a\x28\x0d\x4f\x84\x81\x46\x88\x82\xd6\x18\x18\xb9\x8f\xf2\xc9\xe0\x63\xff\x53\xac\x34\x35\x00\xb2\xb9\xd3\x92\xf2\x04\x5b\xf0\xce\xdb\xd9\x2e\xcc\xa1\x2a\x05\x0a\x4c\x73\x2d\x26\x3e\xdd\xf7\x1d\xe1\x6d\xca\x93\x32\xc6\xa6\x09\xf3\xaf\x6e\xd2\xab\xaf\x9f\xff\xf6\x4d\x7a\x75\x9d\xfe\xf1\x2c\xfd\xef\xe7\xfc\x49\x76\x63\x31\xb2\x3e\xa5\x18\x67\x08\x4a\x8c\x25\x6c\x15\x41\xc3\xb0\x8c\x36\xff\x97\xb5\x9e\x04\xd4\x1c\x66\x1f\x87\x48\x85\x71\xc0\x9c\x53\x1d\xba\x5d\x50\xaf\x8a\x5d\x93\x93\xad\x62\x5b\xe9\x51\xb3\x2e\x0d\xef\xfe\xf7\x6f\x57\x86\x97\xa7\x65\x91\xb1\xd8\x84\x6a\x80\x4a\x06\x9f\x1a\x61\xc1\xb7\x92\x70\xb5\x8f\x4a\x6c\x9a\xf7\xe7\xdf\xcd\x6f\xbf\xb8\xbb\xf9\x3b\xbd\xf8\x09\x46\x85\x1b\xa4\x9f\x5f\xdc\x7f\xf6\x6b\xfa\xe7\xbf\x20\x01\x60\x7e\xfb\xfd\xfc\x87\xcb\xbb\xeb\xf3\xbb\x7f\xbe\x7c\xf1\xff\xa5\x52\xf4\x1f\x7f\x49\xff\xba\xaa\xb3\x52\x31\x97\x87\xdd\x00\x33\xc8\x6e\x17\x56\xbd\x29\x02\xd7\x30\x5c\xda\x4c\x0f\x63\x7f\x8c\xa5\xf7\x78\x6b\x3c\xc2\xab\x26\xb0\xf8\x05\x72\x46\x12\xf9\x46\x76\x1a\xd3\x42\xf2\x81\xcb\xef\x72\x56\xb9\xb8\x66\x2f\x03\x00\x00\xff\xff\x8c\xe4\xcb\x0e\xec\x08\x00\x00")

func staticJsMainJsBytes() ([]byte, error) {
	return bindataRead(
		_staticJsMainJs,
		"static/js/main.js",
	)
}

func staticJsMainJs() (*asset, error) {
	bytes, err := staticJsMainJsBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "static/js/main.js", size: 2284, mode: os.FileMode(420), modTime: time.Unix(1512097443, 0)}
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
	"static/js/main.js": staticJsMainJs,
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
	"static": &bintree{nil, map[string]*bintree{
		"js": &bintree{nil, map[string]*bintree{
			"main.js": &bintree{staticJsMainJs, map[string]*bintree{}},
		}},
	}},
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

