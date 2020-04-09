package rice_migrator

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	rice "github.com/GeertJohan/go.rice"
	"github.com/golang-migrate/migrate/v4/source"
)

func init() {
	source.Register("rice", &Box{})
}

// Box is a migrate driver
type Box struct {
	box        *rice.Box
	migrations *source.Migrations
}

// NewBox initializes a Box driver
func NewBox(migrations *rice.Box) *Box {
	return &Box{
		box:        migrations,
		migrations: source.NewMigrations(),
	}
}

// Initialize reads the rice-box for all migrations
func (b *Box) Initialize() error {
	initial := true
	return b.box.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if initial {
			initial = false
			return nil
		}
		if info.IsDir() {
			return filepath.SkipDir
		}
		// skip files we can't parse
		if m, err := source.DefaultParse(path); err == nil {
			if !b.migrations.Append(m) {
				return fmt.Errorf("unable to parse file %v", path)
			}
		}

		return nil
	})
}

// Open just throws an error, we need to use NewWithSourceInstance
func (b *Box) Open(url string) (source.Driver, error) {
	return nil, fmt.Errorf("not yet implemented")
}

// Close is a no-op
func (b *Box) Close() error {
	return nil
}

// First returns the very first migration version available to the driver.
// Migrate will call this function multiple times.
func (b *Box) First() (version uint, err error) {
	if v, ok := b.migrations.First(); !ok {
		return 0, &os.PathError{Op: "first", Path: "<rice>", Err: os.ErrNotExist}
	} else {
		return v, nil
	}
}

// Prev returns the previous version for a given version available to the driver.
// Migrate will call this function multiple times.
func (b *Box) Prev(version uint) (prevVersion uint, err error) {
	if v, ok := b.migrations.Prev(version); !ok {
		return 0, &os.PathError{Op: fmt.Sprintf("prev for version %v", version), Path: "<rice>", Err: os.ErrNotExist}
	} else {
		return v, nil
	}
}

// Next returns the next version for a given version available to the driver.
// Migrate will call this function multiple times.
func (b *Box) Next(version uint) (nextVersion uint, err error) {
	if v, ok := b.migrations.Next(version); !ok {
		return 0, &os.PathError{Op: fmt.Sprintf("next for version %v", version), Path: "<rice>", Err: os.ErrNotExist}
	} else {
		return v, nil
	}
}

// ReadUp returns the UP migration body and an identifier that helps
// finding this migration in the source for a given version.
// If there is no up migration available for this version,
// it must return os.ErrNotExist.
func (b *Box) ReadUp(version uint) (r io.ReadCloser, identifier string, err error) {
	if m, ok := b.migrations.Up(version); ok {
		file, err := b.box.Open(m.Raw)
		if err != nil {
			return nil, "", err
		}
		return file, m.Identifier, nil
	}
	return nil, "", &os.PathError{Op: fmt.Sprintf("read version %v", version), Path: "<rice>", Err: os.ErrNotExist}
}

// ReadDown returns the DOWN migration body and an identifier that helps
// finding this migration in the source for a given version.
// If there is no down migration available for this version,
// it must return os.ErrNotExist.
func (b *Box) ReadDown(version uint) (r io.ReadCloser, identifier string, err error) {
	if m, ok := b.migrations.Down(version); ok {
		file, err := b.box.Open(m.Raw)
		if err != nil {
			return nil, "", err
		}
		return file, m.Identifier, nil
	}
	return nil, "", &os.PathError{Op: fmt.Sprintf("read version %v", version), Path: "<rice>", Err: os.ErrNotExist}
}
