//go:generate go run ./gen
package buildfs_test

import (
	"embed"
	"fmt"
	"io/fs"
	"log"

	"github.com/l4go/buildfs"
)

func MustSub(fsys fs.FS, subpath string) fs.FS {
	subfs, err := fs.Sub(fsys, subpath)
	if err != nil {
		panic("Fail fs.Sub: " + err.Error())
	}

	return subfs
}

//go:embed testfs
var raw_testFS embed.FS
var testFS = MustSub(buildfs.BuildInFS(raw_testFS, BuildTime), "testfs")

func Example() {
	f, err := testFS.Open(".")
	if err != nil {
		log.Fatal(err)
	}

	df, ok := f.(fs.ReadDirFile)
	if !ok {
		log.Fatal("Not fs.ReadDirFile")
	}

	dir, err := df.ReadDir(-1)
	if err != nil {
		log.Fatal(err)
	}

	for _, de := range dir {
		fi, err := de.Info()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(de.Name(), fi.ModTime().IsZero())
	}

	// Unordered output:
	// foo.txt false
	// bar.txt false
}

func Example_file() {
	f, err := testFS.Open("foo.txt")
	if err != nil {
		log.Fatal(err)
	}

	fi, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fi.Name(), fi.ModTime().IsZero())

	// Output:
	// foo.txt false
}

func Example_dir() {
	f, err := testFS.Open(".")
	if err != nil {
		log.Fatal(err)
	}

	fi, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fi.IsDir(), fi.ModTime().IsZero())

	// Output:
	// true false
}
