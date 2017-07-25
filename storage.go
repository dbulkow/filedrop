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

type MetaData struct {
	Type     StorageType `json:"type"`
	From     string      `json:"from"` // IP address of uploader
	Filename string      `json:"filename"`
	Hash     string      `json:"hash"`
	Created  time.Time   `json:"created"`
	Expire   time.Time   `json:"expire"`
}

func (m *MetaData) MkHash() {
	data := []byte(fmt.Sprintf("%s %s %s", m.Filename, m.Created, m.Expire))
	sum := sha256.Sum256(data)
	m.Hash = fmt.Sprintf("%x", sum)
}

func (m *MetaData) Bytes() []byte {
	b, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		log.Fatalf("marshal: %v", err)
	}

	return b
}

type Storage struct {
	Root  string
	Files map[string]*MetaData
	sync.Mutex
}

func NewStorage(root string) *Storage {
	s := &Storage{Root: root}

	s.Files = make(map[string]*MetaData)

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
		s.Files[hash] = md
		activeFiles.Inc()

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

func (s *Storage) Create(md *MetaData) (io.WriteCloser, error) {
	md.MkHash()

	err := os.MkdirAll(path.Join(s.Root, md.Hash), 0755)
	if err != nil {
		return nil, fmt.Errorf("mkdir: %v", err)
	}

	err = ioutil.WriteFile(path.Join(s.Root, md.Hash, "metadata"), md.Bytes(), 0644)
	if err != nil {
		return nil, fmt.Errorf("create metadata: %v", err)
	}

	s.Files[md.Hash] = md
	activeFiles.Inc()

	return os.Create(path.Join(s.Root, md.Hash, md.Filename))
}

func (s *Storage) Expire() {
	for {
		time.Sleep(time.Minute)

		s.Lock()

		now := time.Now()
		for _, md := range s.Files {
			if now.After(md.Expire) {
				log.Printf("expire %s\n", md.Hash)
				if err := os.RemoveAll(path.Join(s.Root, md.Hash)); err != nil {
					log.Printf("remove %s: %v", md.Hash, err)
				}
				delete(s.Files, md.Hash)
				activeFiles.Dec()
			}
		}

		s.Unlock()
	}
}
