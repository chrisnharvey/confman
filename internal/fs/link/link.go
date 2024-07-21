package link

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Symlink struct {
	Source      string
	Destination string
}

func NewLink(source, destination string) *Symlink {
	return &Symlink{
		Source:      source,
		Destination: destination,
	}
}

func NewLinks(links map[string]string) []*Symlink {
	var linkSlice []*Symlink

	for source, dest := range links {
		linkSlice = append(linkSlice, NewLink(source, dest))
	}

	return linkSlice
}

func (l *Symlink) GetFullPath() string {
	return "/confman/" + l.Destination
}

func (l *Symlink) CanBeLinked() bool {
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

func (l *Symlink) SourceContentsIsSame() bool {
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

func (l *Symlink) GetSourceHash() (string, error) {
	return getHashForPath(l.Source)
}

func (l *Symlink) GetDestinationHash() (string, error) {
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

func (l *Symlink) DestinationExists() bool {
	if _, err := os.Stat(l.GetFullPath()); err != nil {
		return false
	}

	return true
}

func (l *Symlink) SourceExists() bool {
	if _, err := os.Stat(l.Source); err != nil {
		return false
	}

	return true
}

func (l *Symlink) IsSourceSymlink() bool {
	if _, err := os.Readlink(l.Source); err != nil {
		return false
	}

	return true
}

func (l *Symlink) GetSymlinkTarget() (string, error) {
	return os.Readlink(l.Source)
}

func (l *Symlink) IsLinked() bool {
	target, err := l.GetSymlinkTarget()
	if err != nil {
		return false
	}

	if target != l.GetFullPath() {
		return false
	}

	return true
}

func (l *Symlink) Link() error {
	if !l.CanBeLinked() {
		return errors.New("link cannot be created")
	}

	return os.Symlink(l.GetFullPath(), l.Source)
}

func (l *Symlink) Unlink() error {
	if !l.IsLinked() {
		return errors.New("link does not exist")
	}

	return os.Remove(l.Source)
}

func (l *Symlink) Create() error {
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

func (l *Symlink) Restore() error {
	if !l.IsLinked() {
		return errors.New("link does not exist")
	}

	if !l.DestinationExists() {
		return errors.New("destination does not exist")
	}

	if !l.SourceExists() {
		return errors.New("source does not exist")
	}

	// remove symlink
	err := os.Remove(l.Source)
	if err != nil {
		return err
	}

	// create file for restore
	destFile, err := os.Create(l.Source)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// open file to copy
	file, err := os.Open(l.GetFullPath())
	if err != nil {
		return err
	}
	defer file.Close()

	// copy file back to its original location
	_, err = io.Copy(destFile, file)
	if err != nil {
		return fmt.Errorf("could not copy file: %v", err)
	}

	// remove source file
	err = os.Remove(l.GetFullPath())
	if err != nil {
		return fmt.Errorf("could not remove source file: %v", err)
	}

	return nil
}
