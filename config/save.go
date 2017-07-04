package config

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
)

// Project is a collection of SaveLoadStringers
type Project []SaveLoadStringer

// NewProject returns a Project as a collection of the given SaveLoadStringers
func NewProject(SLSs ...SaveLoadStringer) Project {
	return Project(SLSs)
}

// SaveLoadStringer represents a unit of wapty that allows resuming a previous state
type SaveLoadStringer interface {
	Save(io.Writer) error
	Load(io.Reader) error
	// for debug purposes
	fmt.Stringer
}

// SaveAll invokes all the "Save" methods of the project, creating a zip file containing the status.
// The old file will be removed only on successful save.
func (p Project) SaveAll(workspace string) error {
	var errorlist []error
	out, err := os.OpenFile(workspace+".status.zip", os.O_CREATE|os.O_WRONLY, 0660)
	if err != nil {
		return err
	}
	w := zip.NewWriter(out)

	for _, sls := range p {
		f, err := w.Create(sls.String())
		if err != nil {
			errorlist = append(errorlist, err)
			continue
		}

		err = sls.Save(f)
		if err != nil {
			errorlist = append(errorlist, err)
			continue
		}

		err = w.Flush()
		if err != nil {
			errorlist = append(errorlist, err)
			continue
		}
	}

	if len(errorlist) > 0 {
		buf := bytes.NewBuffer(nil)
		for _, err := range errorlist {
			buf.WriteString(err.Error() + "\n")
		}
		return errors.New(string(buf.Bytes()))
	}

	return os.Rename(workspace+".status.zip", workspace+"status.zip")
}

func (p Project) LoadAll(workspace string) error {
	return nil
}
