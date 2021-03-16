package variables

import (
	"flag"

	"github.com/earthly/earthly/domain"
	"github.com/pkg/errors"
)

type ArgMeta struct {
	Name         string
	Description  string
	DefaultValue string
	VarType      Type
	Optional     bool
}

type ArgsParser struct {
	Ref        domain.Reference
	Flags      []ArgMeta
	Positional []ArgMeta
}

func (ap *ArgsParser) Parse(args []string) (*Scope, error) {
	// Parse flags.
	fs := flag.NewFlagSet(ap.Ref.StringCanonical(), flag.ContinueOnError)
	boolValues := make(map[int]*bool)
	stringValues := make(map[int]*string)
	for flagIndex, meta := range ap.Flags {
		err := ValidateArgType(meta.VarType, meta.DefaultValue)
		if err != nil {
			return nil, errors.Wrapf(err, "processing default value for flag %s", meta.Name)
		}
		switch meta.VarType {
		case StringType:
			stringValues[flagIndex] = fs.String(
				meta.Name, meta.DefaultValue, meta.Description)
		case BoolType:
			boolValues[flagIndex] = fs.Bool(
				meta.Name, meta.DefaultValue == "true", meta.Description)
		default:
			return nil, errors.Errorf("invalid arg type %v", meta.VarType)
		}
	}
	err := fs.Parse(args)
	if err != nil {
		return nil, err
	}

	// Process flag values.
	ret := NewScope()
	for flagIndex, meta := range ap.Flags {
		var value string
		switch meta.VarType {
		case StringType:
			value = *stringValues[flagIndex]
		case BoolType:
			if *boolValues[flagIndex] {
				value = "true"
			} else {
				value = "false"
			}
		default:
			return nil, errors.Errorf("invalid arg type %v", meta.VarType)
		}
		ret.AddActive(meta.Name, Var{
			Value: value,
			Type:  meta.VarType,
		})
	}

	// Positional validation.
	minArgs := 0
	maxArgs := 0
	allowOptional := true
	for i := len(ap.Positional) - 1; i >= 0; i-- {
		meta := ap.Positional[i]
		if !meta.Optional {
			minArgs++
		}
		maxArgs++
		if meta.Optional == allowOptional {
			continue
		}
		if meta.Optional {
			return nil, errors.Errorf("only the rightmost positional args may be optional")
		}
		allowOptional = false
	}
	if fs.NArg() > maxArgs {
		return nil, errors.Errorf("too many args provided %+v", fs.Args())
	}
	if fs.NArg() < minArgs {
		return nil, errors.Errorf("not enough args provided %+v", fs.Args())
	}

	// Positional interpretation.
	for index, arg := range fs.Args() {
		meta := ap.Positional[index]
		err := ValidateArgType(meta.VarType, arg)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid value for positional arg %d %s", index, meta.Name)
		}
		ret.AddActive(meta.Name, Var{
			Type:  meta.VarType,
			Value: arg,
		})
	}
	return ret, nil
}

/*
Syntax scribblings...

FROM VAR1=abc DEF=def +target
BUILD VAR1=abc +target
COPY VAR1=abc +target/smthg +another/smthg-else ./
WITH DOCKER --load="VAR1=abc +something"
FROM DOCKERFILE VAR1=abc +target

FROM VAR1=abc DEF=def +target --flag1 --flag2=abc arg1 arg2
BUILD VAR1=abc +target --flag1 --flag2=abc arg1 arg2
COPY "VAR1=abc +target/smthg --flag1 --flag2=abc arg1 arg2" +another/smthg-else ./
WITH DOCKER --load="VAR1=abc +something --flag1 --flag2=abc arg1 arg2"
FROM DOCKERFILE VAR1=abc +target --flag1 --flag2=abc arg1 arg2

target:
	ARG --flag --bool flag1
	ARG --flag flag2
	ARG --positional var1
	ARG --positional var2
	RUN ...
	RUN ...
	ARG OLD_ARG=bla
	ARG ANOTHER_OLD_ARG

earthly OLD_ARG=foo ANOTHER_OLD_ARG=bar +target --flag1 --flag2=abc arg1 arg2

BUILD OLD_ARG=foo ANOTHER_OLD_ARG=bar +target --flag1 --flag2=abc arg1 arg2
FROM OLD_ARG=foo ANOTHER_OLD_ARG=bar +target --flag1 --flag2=abc arg1 arg2
WITH DOCKER --load="OLD_ARG=foo ANOTHER_OLD_ARG=bar +target --flag1 --flag2=abc arg1 arg2"
WITH DOCKER --load="mycontainer:latest=(OLD_ARG=foo ANOTHER_OLD_ARG=bar +target --flag1 --flag2=abc arg1 arg2)"
WITH DOCKER --load="(OLD_ARG=foo ANOTHER_OLD_ARG=bar +target --flag1 --flag2=abc arg1 arg2)"

COPY "OLD_ARG=foo ANOTHER_OLD_ARG=bar +target/smthg --flag1 --flag2=abc arg1 arg2" +another/smthg-else ./
COPY OLD_ARG=foo ANOTHER_OLD_ARG=bar flag1=true flag2=abc var1=arg1 var2=arg2 +target/smthg +another/smthg-else ./
COPY (OLD_ARG=foo ANOTHER_OLD_ARG=bar +target/smthg --flag1 --flag2=abc arg1 arg2) +another/smthg-else ./
*/
