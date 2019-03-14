package opts

import (
	"strings"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type CurlHeader struct {
	Name  string
	Value string
}

func (a *CurlHeader) UnmarshalFlag(data string) error {
	pieces := strings.SplitN(data, ": ", 2)
	if len(pieces) != 2 {
		return bosherr.Errorf("Expected header '%s' to be in format 'name: value'", data)
	}

	if len(pieces[0]) == 0 {
		return bosherr.Errorf("Expected header '%s' to specify non-empty name", data)
	}

	if len(pieces[1]) == 0 {
		return bosherr.Errorf("Expected header '%s' to specify non-empty value", data)
	}

	*a = CurlHeader{Name: pieces[0], Value: pieces[1]}

	return nil
}

type CurlOpts struct {
	Args CurlArgs `positional-args:"true"`

	Method  string       `long:"method" short:"X" description:"HTTP method" default:"GET"`
	Headers []CurlHeader `long:"header" short:"H" description:"HTTP header in 'name: value' format"`
	Body    FileBytesArg `long:"body"             description:"HTTP request body (path)"`

	ShowHeaders bool `long:"show-headers" short:"i"   description:"Show HTTP headers"`

	cmd
}

type CurlArgs struct {
	Path string `positional-arg-name:"PATH" description:"URL path which can include query string"`
}
