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
	"sort"
	"strings"

	"github.com/r0bertz/ripple-go/csv"
	"github.com/r0bertz/ripple/data"
)

var (
	dir     = flag.String("dir", "", "directory contains tx data")
	account = flag.String("account", "", "ripple account")
	related = flag.String("related_accounts", "", "a comma-separated related ripple accounts")
	format  = flag.String("format", "", "csv file format")
	printTx = flag.Bool("print_tx", false, "whether to print an URL to the transaction in the last column")
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

	if _, ok := csv.Factory[*format]; !ok {
		log.Fatalf("unsupported format: %s", *format)
	}

	acct, err := data.NewAccountFromAddress(*account)
	if err != nil {
		log.Fatal(err)
	}
	var ra []data.Account
	if *related != "" {
		for _, a := range strings.Split(*related, ",") {
			r, err := data.NewAccountFromAddress(a)
			if err != nil {
				log.Fatal(err)
			}
			ra = append(ra, *r)
		}
	}

	files, err := ioutil.ReadDir(*dir)
	if err != nil {
		log.Fatal(err)
	}

	s := newFileSet("done")
	if err := s.load(); err != nil {
		log.Fatal(err)
	}

	errList := []string{}
	var rows csv.Slice
	for _, file := range files {
		if s.contains(file.Name()) {
			continue
		}
		content, err := ioutil.ReadFile(path.Join(*dir, file.Name()))
		if err != nil {
			log.Fatal(err)
		}
		row := csv.Factory[*format]()
		if err := row.New(string(content), *acct, ra); err != nil {
			if strings.HasPrefix(err.Error(), "not implemented") {
				errList = append(errList, err.Error())
				continue
			}
			log.Fatal(err)
		}
		s.add(file.Name())
		s.save()
		rows = append(rows, row)
	}
	sort.Sort(rows)
	for _, r := range rows {
		s := fmt.Sprintf("%s", r)
		if *printTx {
			s += fmt.Sprintf(",%s", r.TxURL())
		}
		fmt.Println(s)
	}
	sort.Strings(errList)
	for _, r := range errList {
		fmt.Fprintf(os.Stderr, "%s\n", r)
	}
}
