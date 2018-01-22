package apis

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/empijei/cli/lg"
)

// Command represents a packet of information sent or received by or from the server.
type Command struct {
	Channel UIChannel
	Action  Action
	Args    map[ArgName]string
	Payload []byte
}

// UnpackArgs is used to extract the value of the arguments from a command.
// cmd is the command to extract the values from, names is a list of ArgName
// that is used to access cmd.Args.
// vars can be pointers to either int, bool or string. This function will attempt
// to deserialize the arguments with the proper type and operations in order
// to fit the given vars types.
//
// WARNING: this function PANICS if len(names) != len(vars) since that surely
// means there is a bug in the code.
func (cmd *Command) UnpackArgs(names []ArgName, vars ...interface{}) (err error) {
	if nargs, nvars := len(cmd.Args), len(vars); nargs != nvars {
		err := fmt.Sprintf("wrong number of parameters, expected %d but got %d. Args: <%#v>", nvars, nargs, cmd.Args)
		lg.Error(err)
		return errors.New(err)
	}
	if nnames, nvars := len(names), len(vars); nnames != nvars {
		lg.Failuref("wrong call to ArgsUnpack: given %d names but got %d variables to store them", nnames, nvars)
	}
	for i := 0; i < len(vars); i++ {
		arg := cmd.Args[names[i]]
		switch _var := vars[i].(type) {
		case *int:
			*_var, err = strconv.Atoi(arg)
			if err != nil {
				lg.Infof("cannot read <%s> as int: %s", arg, err.Error())
				return err
			}
		case *bool:
			*_var = arg == ARG_TRUE
		case *string:
			*_var = arg
		default:
			err := fmt.Sprintf("unsupported type passed to ArgsUnpack: %s, only supports pointers to int, string, bool", reflect.TypeOf(_var))
			lg.Error(err)
			return errors.New(err)
		}
	}
	return nil
}

// PackArgs is used to set the value of the arguments of a command.
//
// WARNING: this function PANICS if len(names) != len(vars) since that surely
// means there is a bug in the code.
func (cmd *Command) PackArgs(names []ArgName, vars ...string) {
	if nnames, nvars := len(names), len(vars); nnames != nvars {
		lg.Failuref("wrong call to ArgsUnpack: given %d names but got %d variables to store them", nnames, nvars)
	}
	cmd.Args = make(map[ArgName]string)
	for i, name := range names {
		cmd.Args[name] = vars[i]
	}
}
