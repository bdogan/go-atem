package types

import (
	"fmt"
)

type Version struct {
	Major uint16
	Minor uint16
}

func (v *Version) String() string {
	return fmt.Sprintf("%d.%d", v.Major, v.Minor)
}