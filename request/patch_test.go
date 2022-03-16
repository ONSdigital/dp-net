package request

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetPatches(t *testing.T) {
	Convey("Given a patch request", t, func() {
		body := strings.NewReader(`[
			{ "op": "test", "path": "/a/b/c", "value": "foo" },
			{ "op": "remove", "path": "/a/b/c" },
			{ "op": "add", "path": "/a/b/c", "value": [ "foo", "bar" ] },
			{ "op": "replace", "path": "/a/b/c", "value": 42 },
			{ "op": "move", "from": "/a/b/c", "path": "/a/b/d" },
			{ "op": "copy", "from": "/a/b/d", "path": "/a/b/e" }
		]`)
		req := httptest.NewRequest(http.MethodPatch, "http://localhost:21800/jobs/12345", body)

		Convey("When GetPatches is called", func() {
			patches, err := GetPatches(req.Body)

			Convey("Then all the patches are in the form of []Patch", func() {
				So(err, ShouldBeNil)
				So(patches, ShouldHaveLength, 6)

				So(patches[0].Op, ShouldEqual, "test")
				So(patches[0].Path, ShouldEqual, "/a/b/c")
				So(patches[0].From, ShouldBeEmpty)
				So(patches[0].Value, ShouldEqual, "foo")

				So(patches[1].Op, ShouldEqual, "remove")
				So(patches[1].Path, ShouldEqual, "/a/b/c")
				So(patches[1].From, ShouldBeEmpty)
				So(patches[1].Value, ShouldBeEmpty)

				So(patches[2].Op, ShouldEqual, "add")
				So(patches[2].Path, ShouldEqual, "/a/b/c")
				So(patches[2].From, ShouldBeEmpty)
				So(patches[2].Value, ShouldResemble, []interface{}{"foo", "bar"})

				So(patches[3].Op, ShouldEqual, "replace")
				So(patches[3].Path, ShouldEqual, "/a/b/c")
				So(patches[3].From, ShouldBeEmpty)
				So(patches[3].Value, ShouldEqual, 42)

				So(patches[4].Op, ShouldEqual, "move")
				So(patches[4].Path, ShouldEqual, "/a/b/d")
				So(patches[4].From, ShouldEqual, "/a/b/c")
				So(patches[4].Value, ShouldBeEmpty)

				So(patches[5].Op, ShouldEqual, "copy")
				So(patches[5].Path, ShouldEqual, "/a/b/e")
				So(patches[5].From, ShouldEqual, "/a/b/d")
				So(patches[5].Value, ShouldBeEmpty)
			})
		})
	})

	Convey("Given a patch request with unknown patch operation given in request body", t, func() {
		body := strings.NewReader(`[
			{ "op": "invalid", "path": "/a/b/c", "value": "foo" }
		]`)
		req := httptest.NewRequest(http.MethodPatch, "http://localhost:21800/jobs/12345", body)

		Convey("When GetPatches is called", func() {
			patches, err := GetPatches(req.Body)

			Convey("Then an error should be returned ", func() {
				So(err, ShouldNotBeNil)
				So(err, ShouldResemble, fmt.Errorf("failed to validate patch - error: operation is missing or not valid. Please, provide one of the following: %v", validOps))

				Convey("And an empty patch array should be returned", func() {
					So(patches, ShouldBeEmpty)
				})
			})
		})
	})

	Convey("Given a patch request with an array of unknown patch operation given in request body", t, func() {
		body := strings.NewReader(`[]`)
		req := httptest.NewRequest(http.MethodPatch, "http://localhost:21800/jobs/12345", body)

		Convey("When GetPatches is called", func() {
			patches, err := GetPatches(req.Body)

			Convey("Then an error should be returned ", func() {
				So(err, ShouldNotBeNil)
				So(err, ShouldResemble, fmt.Errorf("no patches given in request body"))

				Convey("And an empty patch array should be returned", func() {
					So(patches, ShouldBeEmpty)
				})
			})
		})
	})

	Convey("Given a patch request with unknown patch operation given in request body", t, func() {
		body := strings.NewReader(`{}`)
		req := httptest.NewRequest(http.MethodPatch, "http://localhost:21800/jobs/12345", body)

		Convey("When GetPatches is called", func() {
			patches, err := GetPatches(req.Body)

			Convey("Then an error should be returned ", func() {
				So(err, ShouldNotBeNil)
				So(err, ShouldResemble, fmt.Errorf("failed to unmarshal patch request body - error: json: cannot unmarshal object into Go value of type []request.Patch"))

				Convey("And an empty patch array should be returned", func() {
					So(patches, ShouldBeEmpty)
				})
			})
		})
	})
}

func TestValidate(t *testing.T) {

	Convey("Validating a valid patch with a supported op and array of strings value is successful", t, func() {
		patch := Patch{
			Op:    "add",
			Path:  "/a/b/c",
			Value: []string{"foo"},
		}
		So(patch.Validate(OpAdd), ShouldBeNil)
	})

	Convey("Validating a valid patch with a supported op and float64 value is successful", t, func() {
		patch := Patch{
			Op:    "add",
			Path:  "/a/b/c",
			Value: float64(123.321),
		}
		So(patch.Validate(OpAdd), ShouldBeNil)
	})

	Convey("Validating a patch struct with an invalid op fails with the expected error", t, func() {
		patch := Patch{
			Op:    "wrong",
			Path:  "/a/b/c",
			Value: []string{"foo"},
		}
		So(patch.Validate(), ShouldResemble, ErrInvalidOp)
	})

	Convey("Validating a valid patch with an unsupported op fails with the expected error", t, func() {
		patch := Patch{
			Op:    "add",
			Path:  "/a/b/c",
			Value: []string{"foo"},
		}
		So(patch.Validate(OpRemove), ShouldResemble, ErrUnsupportedOp("add", []PatchOp{OpRemove}))
	})

	Convey("Validating a patch struct with missing members for an operation results in the expected error being returned", t, func() {
		patch := Patch{
			Op:   "add",
			Path: "/a/b/c",
		}
		So(patch.Validate(OpAdd), ShouldResemble, ErrMissingMember([]string{"value"}))
		patch = Patch{
			Op:    "replace",
			Value: []string{"foo"},
		}
		So(patch.Validate(OpReplace), ShouldResemble, ErrMissingMember([]string{"path"}))
		patch = Patch{
			Op: "test",
		}
		So(patch.Validate(OpTest), ShouldResemble, ErrMissingMember([]string{"path", "value"}))
		patch = Patch{
			Op: "remove",
		}
		So(patch.Validate(OpRemove), ShouldResemble, ErrMissingMember([]string{"path"}))
		patch = Patch{
			Op:   "move",
			Path: "/a/b/c",
		}
		So(patch.Validate(OpMove), ShouldResemble, ErrMissingMember([]string{"from"}))
		patch = Patch{
			Op:   "copy",
			From: "/c/b/a",
		}
		So(patch.Validate(OpCopy), ShouldResemble, ErrMissingMember([]string{"path"}))
	})

}
