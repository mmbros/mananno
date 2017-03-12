// Code generated by go-bindata.
// sources:
// templates/tmpl/arenavision/schedule.tmpl
// templates/tmpl/ilcorsaronero/_base.tmpl
// templates/tmpl/ilcorsaronero/_browser-sync.tmpl
// templates/tmpl/ilcorsaronero/index.tmpl
// templates/tmpl/ilcorsaronero/index.tmpl.old
// templates/tmpl/partials/_footer.tmpl
// templates/tmpl/partials/_header.tmpl
// templates/tmpl/test/transmission.tmpl
// DO NOT EDIT!

package templates

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// bindataRead reads the given file from disk. It returns an error on failure.
func bindataRead(path, name string) ([]byte, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset %s at %s: %v", name, path, err)
	}
	return buf, err
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

// arenavisionScheduleTmpl reads file data from disk. It returns an error on failure.
func arenavisionScheduleTmpl() (*asset, error) {
	path := "/home/user/Code/go/src/github.com/mmbros/mananno/templates/tmpl/arenavision/schedule.tmpl"
	name := "arenavision/schedule.tmpl"
	bytes, err := bindataRead(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// ilcorsaronero_baseTmpl reads file data from disk. It returns an error on failure.
func ilcorsaronero_baseTmpl() (*asset, error) {
	path := "/home/user/Code/go/src/github.com/mmbros/mananno/templates/tmpl/ilcorsaronero/_base.tmpl"
	name := "ilcorsaronero/_base.tmpl"
	bytes, err := bindataRead(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// ilcorsaronero_browserSyncTmpl reads file data from disk. It returns an error on failure.
func ilcorsaronero_browserSyncTmpl() (*asset, error) {
	path := "/home/user/Code/go/src/github.com/mmbros/mananno/templates/tmpl/ilcorsaronero/_browser-sync.tmpl"
	name := "ilcorsaronero/_browser-sync.tmpl"
	bytes, err := bindataRead(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// ilcorsaroneroIndexTmpl reads file data from disk. It returns an error on failure.
func ilcorsaroneroIndexTmpl() (*asset, error) {
	path := "/home/user/Code/go/src/github.com/mmbros/mananno/templates/tmpl/ilcorsaronero/index.tmpl"
	name := "ilcorsaronero/index.tmpl"
	bytes, err := bindataRead(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// ilcorsaroneroIndexTmplOld reads file data from disk. It returns an error on failure.
func ilcorsaroneroIndexTmplOld() (*asset, error) {
	path := "/home/user/Code/go/src/github.com/mmbros/mananno/templates/tmpl/ilcorsaronero/index.tmpl.old"
	name := "ilcorsaronero/index.tmpl.old"
	bytes, err := bindataRead(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// partials_footerTmpl reads file data from disk. It returns an error on failure.
func partials_footerTmpl() (*asset, error) {
	path := "/home/user/Code/go/src/github.com/mmbros/mananno/templates/tmpl/partials/_footer.tmpl"
	name := "partials/_footer.tmpl"
	bytes, err := bindataRead(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// partials_headerTmpl reads file data from disk. It returns an error on failure.
func partials_headerTmpl() (*asset, error) {
	path := "/home/user/Code/go/src/github.com/mmbros/mananno/templates/tmpl/partials/_header.tmpl"
	name := "partials/_header.tmpl"
	bytes, err := bindataRead(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// testTransmissionTmpl reads file data from disk. It returns an error on failure.
func testTransmissionTmpl() (*asset, error) {
	path := "/home/user/Code/go/src/github.com/mmbros/mananno/templates/tmpl/test/transmission.tmpl"
	name := "test/transmission.tmpl"
	bytes, err := bindataRead(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
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
	"arenavision/schedule.tmpl": arenavisionScheduleTmpl,
	"ilcorsaronero/_base.tmpl": ilcorsaronero_baseTmpl,
	"ilcorsaronero/_browser-sync.tmpl": ilcorsaronero_browserSyncTmpl,
	"ilcorsaronero/index.tmpl": ilcorsaroneroIndexTmpl,
	"ilcorsaronero/index.tmpl.old": ilcorsaroneroIndexTmplOld,
	"partials/_footer.tmpl": partials_footerTmpl,
	"partials/_header.tmpl": partials_headerTmpl,
	"test/transmission.tmpl": testTransmissionTmpl,
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
	"arenavision": &bintree{nil, map[string]*bintree{
		"schedule.tmpl": &bintree{arenavisionScheduleTmpl, map[string]*bintree{}},
	}},
	"ilcorsaronero": &bintree{nil, map[string]*bintree{
		"_base.tmpl": &bintree{ilcorsaronero_baseTmpl, map[string]*bintree{}},
		"_browser-sync.tmpl": &bintree{ilcorsaronero_browserSyncTmpl, map[string]*bintree{}},
		"index.tmpl": &bintree{ilcorsaroneroIndexTmpl, map[string]*bintree{}},
		"index.tmpl.old": &bintree{ilcorsaroneroIndexTmplOld, map[string]*bintree{}},
	}},
	"partials": &bintree{nil, map[string]*bintree{
		"_footer.tmpl": &bintree{partials_footerTmpl, map[string]*bintree{}},
		"_header.tmpl": &bintree{partials_headerTmpl, map[string]*bintree{}},
	}},
	"test": &bintree{nil, map[string]*bintree{
		"transmission.tmpl": &bintree{testTransmissionTmpl, map[string]*bintree{}},
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

