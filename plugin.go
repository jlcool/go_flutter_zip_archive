// File: main.go
package ziparchive

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/go-flutter-desktop/go-flutter"
	"github.com/go-flutter-desktop/go-flutter/plugin"
	"github.com/pkg/errors"
	"github.com/yeka/zip"
)

//  Make sure to use the same channel name as was used on the Flutter client side.
const channelName = "flutter_zip_archive"
const (
	PARAM_SRC      = "src"
	PARAM_ZIP      = "zip"
	PARAM_DEST     = "dest"
	PARAM_PASSWORD = "password"
)

type ZipArchivePlugin struct {
}

var _ flutter.Plugin = &ZipArchivePlugin{} // compile-time type check

func (p *ZipArchivePlugin) InitPlugin(messenger plugin.BinaryMessenger) error {
	channel := plugin.NewMethodChannel(messenger, channelName, plugin.StandardMethodCodec{})
	channel.HandleFunc("zip", handleZip)
	channel.HandleFunc("unzip", handleUnZip)
	return nil // no error
}

func handleZip(arguments interface{}) (reply interface{}, err error) {
	var ok bool
	var args map[interface{}]interface{}
	if args, ok = arguments.(map[interface{}]interface{}); !ok {
		return nil, errors.New("invalid arguments")
	}
	var src string
	var dest string
	var password string

	if src1, ok := args[PARAM_SRC]; ok {
		src = src1.(string)
	}
	if dest1, ok := args[PARAM_DEST]; ok {
		dest = dest1.(string)
	}
	if password1, ok := args[PARAM_PASSWORD]; ok {
		password = password1.(string)
	}

	fzip, err := os.Create(dest)
	if err != nil {
		log.Fatalln(err)
	}
	zipw := zip.NewWriter(fzip)
	defer zipw.Close()

	f, _ := os.Stat(src)
	if f.IsDir() {
		files, _ := ioutil.ReadDir(src)
		for _, file := range files {

			w, err := zipw.Encrypt(file.Name(), password, zip.StandardEncryption)
			b, err := os.OpenFile(path.Join(src, file.Name()), os.O_RDWR, 0666)
			_, err = io.Copy(w, b)
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		w, err := zipw.Encrypt(f.Name(), password, zip.StandardEncryption)
		b, err := os.OpenFile(src, os.O_RDWR, 0666)
		_, err = io.Copy(w, b)
		if err != nil {
			log.Fatal(err)
		}
	}

	zipw.Flush()

	return nil, nil
}
func handleUnZip(arguments interface{}) (reply interface{}, err error) {
	var ok bool
	var args map[interface{}]interface{}
	if args, ok = arguments.(map[interface{}]interface{}); !ok {
		return nil, errors.New("invalid arguments")
	}
	var zip_src string
	var dest string
	var password string

	if zip_path, ok := args[PARAM_ZIP]; ok {
		zip_src = zip_path.(string)
	}
	if destt, ok := args[PARAM_DEST]; ok {
		dest = destt.(string)
	}
	if passwordd, ok := args[PARAM_PASSWORD]; ok {
		password = passwordd.(string)
	}
	if zip_src == "" {
		log.Printf("[zip] %v\n", "invalid zip_src")
		return nil, errors.New("invalid zip_src")
	}

	r, err := zip.OpenReader(zip_src)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	os.MkdirAll(dest, os.ModePerm)
	// Iterate through the files in the archive,
	// printing some of their contents.
	for _, f := range r.File {
		if f.IsEncrypted() {
			f.SetPassword(password)
		}

		rc, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}
		buf, err := ioutil.ReadAll(rc)
		if err != nil {
			log.Fatal(err)
		}
		rc.Close()
		if f.FileInfo().IsDir() {
			os.MkdirAll(path.Join(dest, f.Name), os.ModePerm)
		} else {
			fmt.Printf("WriteFile %s:\n", dest+f.Name)
			ioutil.WriteFile(path.Join(dest, f.Name), buf, os.ModePerm)
		}
		fmt.Println()
	}
	return nil, nil
}
