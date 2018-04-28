package main

import (
	"bytes"
	"container/heap"
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/golang/glog"
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

	ff, ok := csv.FormatterFactory[*format]
	if !ok {
		glog.Fatalf("unsupported format: %s", *format)
	}
	formatter := ff()

	acct, err := data.NewAccountFromAddress(*account)
	if err != nil {
		glog.Fatal(err)
	}
	var ra []data.Account
	if *related != "" {
		for _, a := range strings.Split(*related, ",") {
			r, err := data.NewAccountFromAddress(a)
			if err != nil {
				glog.Fatal(err)
			}
			ra = append(ra, *r)
		}
	}

	files, err := ioutil.ReadDir(*dir)
	if err != nil {
		glog.Fatal(err)
	}

	s := newFileSet("done")
	if err := s.load(); err != nil {
		glog.Fatal(err)
	}

	errList := []string{}
	c := csv.New(*acct, ra)
	for _, file := range files {
		if s.contains(file.Name()) {
			continue
		}
		content, err := ioutil.ReadFile(path.Join(*dir, file.Name()))
		if err != nil {
			glog.Fatal(err)
		}
		var resp csv.TxResponse
		dec := json.NewDecoder(strings.NewReader(string(content)))
		if err := dec.Decode(&resp); err != nil {
			glog.Fatalf("error decoding transaction: %v", err)
		}
		t := resp.Result
		if err := c.Add(t.TransactionWithMetaData); err != nil {
			if strings.HasPrefix(err.Error(), "not implemented") {
				errList = append(errList, err.Error())
				goto out
			}
			glog.Fatal(err)
		}
	out:
		s.add(file.Name())
		s.save()
	}
	fmt.Println(formatter.Header())
	for c.Rows.Len() > 0 {
		r := heap.Pop(&c.Rows).(csv.Row)
		s, err := formatter.Format(r)
		if err != nil {
			if strings.HasPrefix(err.Error(), "not implemented") {
				errList = append(errList, err.Error())
			}
			continue
		}
		if *printTx {
			s += fmt.Sprintf(",%s", r.URL())
		}
		fmt.Println(s)
	}
	sort.Strings(errList)
	for _, r := range errList {
		fmt.Fprintf(os.Stderr, "%s\n", r)
	}
}
