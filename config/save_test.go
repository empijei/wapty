package config

import (
	"io"
	"io/ioutil"
	"os"
	"testing"
)

var (
	workspace string
	p         *Project
	data      string
)

type mockSaveLoadStringer struct {
	data, name string
}

func (m *mockSaveLoadStringer) Save(w io.Writer) error {
	_, err := w.Write([]byte(m.data))
	return err
}

func (m *mockSaveLoadStringer) Load(r io.Reader) error {
	tmp, err := ioutil.ReadAll(r)
	m.data = string(tmp)
	return err
}

func (m *mockSaveLoadStringer) String() string {
	return m.name
}

func TestSaveAll(t *testing.T) {
	// FIXME debug "not a valid zip file" error
	workspace = os.TempDir() + string(os.PathSeparator)

	pSave := Project{
		&mockSaveLoadStringer{
			name: "package1",
			data: "data of package1",
		},
		&mockSaveLoadStringer{
			name: "package2",
			data: "data of package2",
		},
	}

	err := pSave.SaveAll(workspace)
	if err != nil {
		t.Error(err)
	}

	pLoad := Project{
		&mockSaveLoadStringer{
			name: "package1",
		},
		&mockSaveLoadStringer{
			name: "package2",
		},
	}

	err = pLoad.LoadAll(workspace)
	if err != nil {
		t.Error(err)
	}

	for _, slsl := range pLoad {
		var found bool
		for _, slss := range pSave {
			if slss.String() == slsl.String() {
				found = true
				break
			}
		}
		if !found {
			t.Error("Laoded alien package " + slsl.String())
		}
	}

	for _, slss := range pSave {
		var found bool
		for _, slsl := range pLoad {
			if slss.String() == slsl.String() {
				found = true
				break
			}
		}
		if !found {
			t.Error("Failed to load package " + slss.String())
		}
	}

}
