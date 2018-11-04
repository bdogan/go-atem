package types

import (
	"bytes"
	"strings"
)

type NullTerminatedString struct {
	Body []byte
	nullPoint int
}

func (nts *NullTerminatedString) String() string  {
	if nts.nullPoint == 0 {
		nts.nullPoint = bytes.IndexByte(nts.Body, 0)
	}
	if nts.nullPoint == -1 {
		return strings.TrimSpace(string(nts.Body))
	}
	return strings.TrimSpace(string(nts.Body[0:nts.nullPoint]))
}