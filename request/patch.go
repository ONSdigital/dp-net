package request

import (
	"encoding/json"
	"fmt"
	"io"
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

// ErrInvalidOp generates an error when a patch contains a wrong 'op'
func ErrInvalidOp(supportedOps []PatchOp) error {
	if len(supportedOps) < 1 {
		return fmt.Errorf("patch operation is invalid as no patch operations are supported")
	}

	validSupportedOps := getPatchOpsStringSlice(supportedOps)
	return fmt.Errorf("patch operation is missing or invalid. Please, provide one of the following: %v", validSupportedOps)
}

// ErrMissingMember generates an error for a missing member
func ErrMissingMember(members []string) error {
	return fmt.Errorf("missing member(s) in patch: %v", members)
}

// ErrUnsupportedOp generates an error for unsupported ops
func ErrUnsupportedOp(op string, supportedOps []PatchOp) error {
	supported := getPatchOpsStringSlice(supportedOps)
	return fmt.Errorf("patch operation '%s' not supported. Supported op(s): %v", op, supported)
}

func (o PatchOp) String() string {
	return validOps[o]
}

func getPatchOpsStringSlice(ops []PatchOp) []string {
	opsStringSlice := []string{}
	for _, op := range ops {
		opsStringSlice = append(opsStringSlice, op.String())
	}
	return opsStringSlice
}

// Patch models an HTTP patch operation request, according to RFC 6902
type Patch struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	From  string      `json:"from"`
	Value interface{} `json:"value"`
}

// GetPatches gets the patches from the request body and returns it in the form of []Patch.
// An error will be returned if request body cannot be read, unmarshalling the requets body is unsuccessful,
// no patches are provided in the request or any of the provided patches are invalid
func GetPatches(requestBody io.ReadCloser, supportedOps []PatchOp) ([]Patch, error) {
	patches := []Patch{}

	if len(supportedOps) < 1 {
		return []Patch{}, fmt.Errorf("empty list of support patch operations given")
	}

	bytes, err := io.ReadAll(requestBody)
	if err != nil {
		return []Patch{}, fmt.Errorf("failed to read and get patch request body")
	}

	if len(bytes) == 0 {
		return []Patch{}, fmt.Errorf("empty request body given")
	}

	err = json.Unmarshal(bytes, &patches)
	if err != nil {
		return []Patch{}, fmt.Errorf("failed to unmarshal patch request body")
	}

	if len(patches) < 1 {
		return []Patch{}, fmt.Errorf("no patches given in request body")
	}

	for _, patch := range patches {
		if err := patch.Validate(supportedOps); err != nil {
			return []Patch{}, err
		}
	}
	return patches, nil
}

// Validate checks that the provided operation is correct and the expected members are provided
func (p *Patch) Validate(supportedOps []PatchOp) error {
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
		return ErrInvalidOp(supportedOps)
	}

	if !p.isOpSupported(supportedOps) {
		return ErrUnsupportedOp(p.Op, supportedOps)
	}

	if len(missing) > 0 {
		return ErrMissingMember(missing)
	}
	return nil
}

// isOpSupported checks that the patch op is in the provided list of supported Ops
func (p *Patch) isOpSupported(supportedOps []PatchOp) bool {
	for _, op := range supportedOps {
		if p.Op == op.String() {
			return true
		}
	}
	return false
}
