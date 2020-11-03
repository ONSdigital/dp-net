package request

import (
	"fmt"
)

// PatchOp - iota enum of possible patch operations
type PatchOp int

// Possible patch operations
const (
	OpAdd PatchOp = iota
	OpRemove
	OpReplace
	OpMove
	OpCopy
	OpTest
)

var validOps = []string{"add", "remove", "replace", "move", "copy", "test"}

// ErrInvalidOp is an error returned when a patch contains a wrong 'op'
var ErrInvalidOp = fmt.Errorf("operation is missing or not valid. Please, provide one of the following: %v", validOps)

// ErrMissingMember generates an error for a missing member
func ErrMissingMember(members []string) error {
	return fmt.Errorf("missing member(s) in patch: %v", members)
}

func (o PatchOp) String() string {
	return validOps[o]
}

// Patch models an HTTP patch operation request, according to RFC 6902
type Patch struct {
	Op    string   `json:"op"`
	Path  string   `json:"path"`
	From  string   `json:"from"`
	Value []string `json:"value"`
}

// Validate checks that the provided operation is correct and the expected members are provided
func (p *Patch) Validate() error {
	missing := []string{}
	switch p.Op {
	case OpAdd.String(), OpReplace.String(), OpTest.String():
		if p.Path == "" {
			missing = append(missing, "path")
		}
		if p.Value == nil {
			missing = append(missing, "value")
		}
	case OpRemove.String():
		if p.Path == "" {
			missing = append(missing, "path")
		}
	case OpMove.String(), OpCopy.String():
		if p.Path == "" {
			missing = append(missing, "path")
		}
		if p.From == "" {
			missing = append(missing, "from")
		}
	default:
		return ErrInvalidOp
	}

	if len(missing) > 0 {
		return ErrMissingMember(missing)
	}
	return nil
}
