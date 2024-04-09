package util

import (
	"archive/tar"
	"bytes"
	"os"

	"github.com/twpayne/go-vfs/v4"
)

type Tar struct {
	files []file

	fs vfs.FS
}

type file struct {
	path string
	data []byte
	inFS bool
}

func NewTar(fs vfs.FS) *Tar {
	return &Tar{
		fs: fs,
	}
}

func (t *Tar) AddPath(path string) *Tar {
	t.files = append(t.files, file{path: path, inFS: true})
	return t
}

func (t *Tar) AddBytes(path string, data []byte) *Tar {
	t.files = append(t.files, file{path: path, data: data, inFS: false})
	return t
}

func (t *Tar) Build() ([]byte, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	for _, f := range t.files {
		if f.inFS {
			if err := t.buildFromFS(tw, f.path); err != nil {
				return nil, err
			}
		} else {
			if err := t.buildFromBytes(tw, f.path, f.data); err != nil {
				return nil, err
			}
		}
	}

	return buf.Bytes(), nil
}

func (t *Tar) buildFromFS(tw *tar.Writer, p string) error {
	return vfs.Walk(t.fs, p, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fi.IsDir() {
			return nil
		}

		f, err := t.fs.ReadFile(path)
		if err != nil {
			return err
		}

		tarHeader := &tar.Header{
			Name: path,
			Size: int64(len(f)),
		}

		if err := tw.WriteHeader(tarHeader); err != nil {
			return err
		}

		if _, err := tw.Write(f); err != nil {
			return err
		}

		return nil
	})
}

func (t *Tar) buildFromBytes(tw *tar.Writer, p string, data []byte) error {
	tarHeader := &tar.Header{
		Name: p,
		Size: int64(len(data)),
	}

	if err := tw.WriteHeader(tarHeader); err != nil {
		return err
	}

	if _, err := tw.Write(data); err != nil {
		return err
	}

	return nil
}
