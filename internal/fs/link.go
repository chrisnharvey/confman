package fs

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Link struct {
	Source      string
	Destination string
}

func NewLink(source, destination string) *Link {
	return &Link{
		Source:      source,
		Destination: destination,
	}
}

func NewLinks(links map[string]string) []*Link {
	var linkSlice []*Link

	for source, dest := range links {
		linkSlice = append(linkSlice, NewLink(source, dest))
	}

	return linkSlice
}

func (l *Link) GetFullPath() string {
	return "/confman/" + l.Destination
}

func (l *Link) CanBeLinked() bool {
	if !l.DestinationExists() {
		return false
	}

	if l.IsSourceSymlink() {
		return false
	}

	if l.IsLinked() {
		return false
	}

	return true
}

func (l *Link) SourceContentsIsSame() bool {
	if !l.SourceExists() || !l.DestinationExists() {
		return false
	}

	sourceHash, err := l.GetSourceHash()
	if err != nil {
		return false
	}

	destHash, err := l.GetDestinationHash()
	if err != nil {
		return false
	}

	return sourceHash == destHash
}

func (l *Link) GetSourceHash() (string, error) {
	return getHashForPath(l.Source)
}

func (l *Link) GetDestinationHash() (string, error) {
	return getHashForPath(l.GetFullPath())
}

func getHashForPath(path string) (string, error) {
	hasher := sha256.New()
	s, err := os.ReadFile(path)
	hasher.Write(s)
	if err != nil {
		return "", err
	}

	return string(hasher.Sum(nil)), nil
}

func (l *Link) DestinationExists() bool {
	if _, err := os.Stat(l.GetFullPath()); err != nil {
		return false
	}

	return true
}

func (l *Link) SourceExists() bool {
	if _, err := os.Stat(l.Source); err != nil {
		return false
	}

	return true
}

func (l *Link) IsSourceSymlink() bool {
	if _, err := os.Readlink(l.Source); err != nil {
		return false
	}

	return true
}

func (l *Link) GetSymlinkTarget() (string, error) {
	return os.Readlink(l.Source)
}

func (l *Link) IsLinked() bool {
	target, err := l.GetSymlinkTarget()
	if err != nil {
		return false
	}

	if target != l.GetFullPath() {
		return false
	}

	return true
}

func (l *Link) Link() error {
	if !l.CanBeLinked() {
		return errors.New("link cannot be created")
	}

	return os.Symlink(l.GetFullPath(), l.Source)
}

func (l *Link) Create() error {
	if l.DestinationExists() {
		return errors.New("destination already exists")
	}

	if !l.SourceExists() {
		return errors.New("source does not exist")
	}

	// open file
	file, err := os.Open(l.Source)
	if err != nil {
		return err
	}
	defer file.Close()

	// create directory if required
	destDir := filepath.Dir(l.GetFullPath())
	err = os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		return err
	}

	// create file
	destFile, err := os.Create(l.GetFullPath())
	if err != nil {
		return err
	}
	defer destFile.Close()

	// copy file
	_, err = io.Copy(destFile, file)
	if err != nil {
		return fmt.Errorf("could not copy file: %v", err)
	}

	// remove source file
	err = os.Remove(l.Source)
	if err != nil {
		os.Remove(l.GetFullPath()) // cleanup
		return fmt.Errorf("could not remove source file: %v", err)
	}

	if err := os.Symlink(l.GetFullPath(), l.Source); err != nil {
		return err
	}

	return nil
}

