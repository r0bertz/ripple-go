package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/r0bertz/ripple-go/csv"
)

var (
	dir     = flag.String("dir", "", "directory contains tx data")
	account = flag.String("account", "", "ripple account")
)

type fileSet struct {
	m    map[string]struct{}
	path string
}

func newFileSet(path string) *fileSet {
	return &fileSet{m: map[string]struct{}{}, path: path}
}

func (f fileSet) add(file string) {
	f.m[file] = struct{}{}
}

func (f fileSet) contains(file string) bool {
	_, ok := f.m[file]
	return ok
}

func (f fileSet) load() error {
	_, err := os.Stat(f.path)
	if os.IsNotExist(err) {
		return ioutil.WriteFile(f.path, []byte{}, 0644)
	}
	c, err := ioutil.ReadFile(f.path)
	if err != nil {
		return err
	}
	d := gob.NewDecoder(bytes.NewReader(c))
	return d.Decode(&f.m)
}

func (f fileSet) save() error {
	b := new(bytes.Buffer)
	e := gob.NewEncoder(b)

	// Encoding the map
	err := e.Encode(f.m)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(f.path, b.Bytes(), 0644)
}

func main() {
	flag.Parse()

	files, err := ioutil.ReadDir(*dir)
	if err != nil {
		log.Fatal(err)
	}

	s := newFileSet("done")
	if err := s.load(); err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if s.contains(file.Name()) {
			continue
		}
		content, err := ioutil.ReadFile(path.Join(*dir, file.Name()))
		if err != nil {
			log.Fatal(err)
		}
		row, err := csv.NewRow(string(content), *account)
		if err != nil {
			if err.Error() == "not implemented" {
				continue
			}
			if strings.HasPrefix(err.Error(), "hash") {
				fmt.Println(err)
				continue
			}
			log.Fatal(err)
		}
		s.add(file.Name())
		s.save()
		fmt.Println(row.String())
	}
}
