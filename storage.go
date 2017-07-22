package main

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"
)

type MetaData struct {
	Filename  string
	Hash      string
	Created   time.Time
	Expire    time.Time
	OwnerMail string
	MailAddrs []string
}

func (m *MetaData) Notify() {
	// iterate over MailAddrs
	// use smtp to provide html mail with link to server
}

func (m *MetaData) DoExpire() {
	if time.Now().After(m.Expire) {
		// remove hashdir
	}
}

func (m *MetaData) MkHash() {
	data := []byte(fmt.Sprintf("%s %s %s", m.Filename, m.Created, m.Expire))
	sum := sha256.Sum256(data)
	m.Hash = fmt.Sprintf("%x", sum)
}

func (m *MetaData) Bytes() []byte {
	return []byte(fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n", m.Filename, m.Hash, m.Created, m.Expire, m.OwnerMail, m.MailAddrs))
}

type Storage struct {
	Root  string
	Files map[string]*MetaData
}

/*
.../storage/<hash of filename, time posted>/file
                                           /metadata
*/

func (s *Storage) Create(md *MetaData) error {
	md.MkHash()

	err := os.MkdirAll(path.Join(s.Root, md.Hash, md.Filename), 0755)
	if err != nil {
		return fmt.Errorf("mkdir: %v", err)
	}

	err = ioutil.WriteFile(path.Join(s.Root, md.Hash, "metadata"), md.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("create metadata: %v", err)
	}

	s.Files[md.Hash] = md

	return nil
}
