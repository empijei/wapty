package main

import (
	"github.com/empijei/wapty/config"
	"github.com/empijei/wapty/intercept"
)

var project = config.NewProject(intercept.GetStatus())
