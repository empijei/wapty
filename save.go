package main

import (
	"fmt"
	"io"

	"github.com/empijei/wapty/intercept"
)

type Saver interface {
	Save(io.Writer) error
}

type Loader interface {
	Load(io.Reader) error
}

type SaveLoadStringer interface {
	Saver
	Loader
	// for debug purposes
	fmt.Stringer
}

var saveLoaders = []SaveLoadStringer{
	intercept.GetStatus(),
}
