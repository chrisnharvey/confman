package link

type Factory struct {
}

type Link interface {
	CanBeLinked() bool
	Create() error
	DestinationExists() bool
	GetDestinationHash() (string, error)
	GetFullPath() string
	GetSourceHash() (string, error)
	GetSymlinkTarget() (string, error)
	IsLinked() bool
	IsSourceSymlink() bool
	Link() error
	Restore() error
	SourceContentsIsSame() bool
	SourceExists() bool
	Unlink() error
}

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) NewLink(source, destination string) Link {
	return &Symlink{
		Source:      source,
		Destination: destination,
	}
}
