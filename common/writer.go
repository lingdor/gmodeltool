package common

import "io"

type StringWriter interface {
	io.Writer
	io.StringWriter
}
