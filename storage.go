package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"
)

type StorageType string

const (
	StorageFile     StorageType = "file"
	StorageText     StorageType = "text"
	StorageGraphics StorageType = "image"
)

type File struct {
	Name       string `json:"filename"`
	Type       string `json:"ctype"`
	Size       int64  `json:"size"`
	downloaded bool
}

type MetaData struct {
	Type    StorageType `json:"type"`
	From    string      `json:"from"` // IP address of uploader
	Files   []File      `json:"files"`
	Hash    string      `json:"hash"`
	Created time.Time   `json:"created"`
	Expire  time.Time   `json:"expire"`
	OnRead  bool        `json:"onread"` // delete after first read/download
	hashdir string
}

func (m *MetaData) MkHash() {
	names := make([]string, 0)
	for _, f := range m.Files {
		names = append(names, f.Name)
	}

	data := []byte(fmt.Sprintf("%s %q %s", names, m.Created, m.Expire))
	sum := sha256.Sum256(data)
	m.Hash = fmt.Sprintf("%x", sum)
}

func (m *MetaData) Marshal() []byte {
	b, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		log.Fatalf("marshal: %v", err)
	}

	return b
}

type Storage struct {
	Root string
	Dirs map[string]*MetaData
	sync.Mutex
}

func NewStorage(root string) *Storage {
	s := &Storage{Root: path.Join(root, "downloads")}

	s.Dirs = make(map[string]*MetaData)

	visit := func(pathname string, fi os.FileInfo, err error) error {
		if path.Base(pathname) != "metadata" {
			return nil
		}

		log.Printf("Loading %s\n", pathname)

		f, err := os.Open(pathname)
		if err != nil {
			return fmt.Errorf("metadata open: %v", err)
		}
		defer f.Close()

		md := &MetaData{}

		if err := json.NewDecoder(f).Decode(md); err != nil {
			return fmt.Errorf("metadata decode: %v", err)
		}

		hash := path.Base(path.Dir(pathname))
		s.Dirs[hash] = md
		activeDirs.Inc()

		return nil
	}

	if err := filepath.Walk(s.Root, visit); err != nil {
		log.Fatalf("walk: %v", err)
	}

	go s.Expire()

	return s
}

/*
.../storage.Root/<hash of filename, time posted>/file
                                                /metadata
*/

func (s *Storage) Mkdir(md *MetaData) error {
	md.MkHash()
	md.hashdir = path.Join(s.Root, md.Hash)

	err := os.MkdirAll(md.hashdir, 0755)
	if err != nil {
		return fmt.Errorf("mkdir: %v", err)
	}

	return nil
}

func (s *Storage) Create(md *MetaData, filename string) (io.WriteCloser, error) {
	hashdir := path.Join(s.Root, md.Hash)
	return os.Create(path.Join(hashdir, filename))
}

func (s *Storage) WriteMeta(md *MetaData) error {
	mdfile := path.Join(md.hashdir, "metadata")

	err := ioutil.WriteFile(mdfile, md.Marshal(), 0644)
	if err != nil {
		return fmt.Errorf("create metadata %s: %v", md.hashdir, err)
	}

	s.Dirs[md.Hash] = md
	activeDirs.Inc()
	activeFiles.Add(float64(len(md.Files)))

	return nil
}

func (s *Storage) expire(md *MetaData) {
	log.Printf("expire %s\n", md.Hash)
	hashdir := path.Join(s.Root, md.Hash)
	if err := os.RemoveAll(hashdir); err != nil {
		log.Printf("remove %s: %v", hashdir, err)
	}
	delete(s.Dirs, md.Hash)
	activeDirs.Dec()
	activeFiles.Sub(float64(len(md.Files)))
}

func (s *Storage) Expire() {
	for {
		time.Sleep(time.Minute)

		s.Lock()

		now := time.Now()
		for _, md := range s.Dirs {
			if now.After(md.Expire) {
				s.expire(md)
			}

			if md.OnRead {
				count := 0
				for _, f := range md.Files {
					if f.downloaded {
						count++
					}
				}
				if len(md.Files) == count {
					s.expire(md)
				}
			}
		}

		s.Unlock()
	}
}
