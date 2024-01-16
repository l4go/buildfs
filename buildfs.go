package buildfs

import (
	"embed"
	"errors"
	"io"
	"io/fs"
	"time"
)

type buildFS struct {
	fsys       embed.FS
	build_time time.Time
}

func BuildInFS(fsys embed.FS, build_time time.Time) fs.FS {
	return &buildFS{fsys: fsys, build_time: build_time}
}

func (bfs buildFS) Open(name string) (fs.File, error) {
	f, err := bfs.fsys.Open(name)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if !s.IsDir() {
		return &buildFile{f: f, mod_time: bfs.build_time}, nil
	}

	return &buildDir{f: f, mod_time: bfs.build_time}, nil
}

func (bf buildFS) ReadDir(name string) ([]fs.DirEntry, error) {
	f, err := bf.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dir, ok := f.(*buildDir)
	if !ok {
		return nil, &fs.PathError{Op: "read", Path: name, Err: errors.New("not a directory")}
	}

	return dir.ReadDir(-1)
}

func (bf buildFS) ReadFile(name string) ([]byte, error) {
	f, err := bf.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return io.ReadAll(f)
}

type buildFileInfo struct {
	fi fs.FileInfo
	mt time.Time
}

func (bfi *buildFileInfo) Name() string {
	return bfi.fi.Name()
}

func (bfi *buildFileInfo) Size() int64 {
	return bfi.fi.Size()
}

func (bfi *buildFileInfo) Mode() fs.FileMode {
	return bfi.fi.Mode()
}

func (bfi *buildFileInfo) ModTime() time.Time {
	return bfi.mt
}

func (bfi *buildFileInfo) IsDir() bool {
	return bfi.fi.IsDir()
}

func (bfi *buildFileInfo) Sys() any {
	return bfi.fi.Sys()
}

type buildDirEntry buildFileInfo

func newBuildDirEntry(bfi *buildFileInfo) *buildDirEntry {
	return (*buildDirEntry)(bfi)
}

func (bde *buildDirEntry) Name() string {
	return (*buildFileInfo)(bde).Name()
}

func (bde *buildDirEntry) IsDir() bool {
	return (*buildFileInfo)(bde).IsDir()
}

func (bde *buildDirEntry) Type() fs.FileMode {
	return (*buildFileInfo)(bde).Mode()
}

func (bde *buildDirEntry) Info() (fs.FileInfo, error) {
	return (*buildFileInfo)(bde), nil
}

type buildFile struct {
	f        fs.File
	mod_time time.Time
}

func (bf *buildFile) Close() error {
	return bf.f.Close()
}

func (bf *buildFile) Stat() (fs.FileInfo, error) {
	fi, err := bf.f.Stat()
	if err != nil {
		return nil, err
	}

	return &buildFileInfo{fi: fi, mt: bf.mod_time}, nil
}

func (bf *buildFile) Read(b []byte) (int, error) {
	return bf.f.Read(b)
}

func (bf *buildFile) Seek(offset int64, whence int) (int64, error) {
	s, ok := bf.f.(io.Seeker)
	if !ok {
		panic("Not embed.FS .")
	}

	return s.Seek(offset, whence)
}

func (bf *buildFile) ReadAt(b []byte, offset int64) (int, error) {
	ra, ok := bf.f.(io.ReaderAt)
	if !ok {
		panic("Not embed.FS .")
	}

	return ra.ReadAt(b, offset)
}

type buildDir struct {
	f        fs.File
	mod_time time.Time
}

func (bd *buildDir) Close() error {
	return bd.f.Close()
}

func (bd *buildDir) Stat() (fs.FileInfo, error) {
	fi, err := bd.f.Stat()
	if err != nil {
		return nil, err
	}

	return &buildFileInfo{fi: fi, mt: bd.mod_time}, nil
}

func (bd *buildDir) Read(b []byte) (int, error) {
	return bd.f.Read(b)
}

func (bd *buildDir) ReadDir(count int) ([]fs.DirEntry, error) {
	rd, ok := bd.f.(fs.ReadDirFile)
	if !ok {
		panic("Not embed.FS .")
	}

	ents, err := rd.ReadDir(count)
	if err != nil {
		return nil, err
	}

	new_ents := make([]fs.DirEntry, len(ents))
	for i, ent := range ents {
		fi, ierr := ent.Info()
		if ierr != nil {
			return nil, ierr
		}
		new_ents[i] = newBuildDirEntry(&buildFileInfo{fi: fi, mt: bd.mod_time})
	}

	return new_ents, nil
}

