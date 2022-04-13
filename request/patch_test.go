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
	Convey("Given a patch request body", t, func() {
		body := strings.NewReader(`[
			{ "op": "test", "path": "/a/b/c", "value": "foo" },
			{ "op": "remove", "path": "/a/b/c" },
			{ "op": "add", "path": "/a/b/c", "value": [ "foo", "bar" ] },
			{ "op": "replace", "path": "/a/b/c", "value": 42 },
			{ "op": "move", "from": "/a/b/c", "path": "/a/b/d" },
			{ "op": "copy", "from": "/a/b/d", "path": "/a/b/e" }
		]`)
		req := httptest.NewRequest(http.MethodPatch, "http://localhost:21800/jobs/12345", body)

		Convey("And all patch operations are supported", func() {
			supportedOps := []PatchOp{
				OpTest,
				OpRemove,
				OpAdd,
				OpReplace,
				OpMove,
				OpCopy,
			}

			Convey("When GetPatches is called", func() {
				patches, err := GetPatches(req.Body, supportedOps)

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
	})

	Convey("Given empty list of supported patch operations", t, func() {
		emptySupportedOps := []PatchOp{}

		Convey("And valid patch request body", func() {
			body := strings.NewReader(`[
				{ "op": "test", "path": "/a/b/c", "value": "foo" }
			]`)
			req := httptest.NewRequest(http.MethodPatch, "http://localhost:21800/jobs/12345", body)

			Convey("When GetPatches is called", func() {
				patches, err := GetPatches(req.Body, emptySupportedOps)

				Convey("Then an error should be returned ", func() {
					So(err, ShouldNotBeNil)
					So(err, ShouldResemble, fmt.Errorf("empty list of support patch operations given"))

					Convey("And an empty patch array should be returned", func() {
						So(patches, ShouldBeEmpty)
					})
				})
			})
		})
	})

	Convey("Given an empty request body", t, func() {
		emptyBody := strings.NewReader("")
		req := httptest.NewRequest(http.MethodPatch, "http://localhost:21800/jobs/12345", emptyBody)

		Convey("And valid supported patch operations given", func() {
			supportedOps := []PatchOp{OpAdd}

			Convey("When GetPatches is called", func() {
				patches, err := GetPatches(req.Body, supportedOps)

				Convey("Then an error should be returned ", func() {
					So(err, ShouldNotBeNil)
					So(err, ShouldResemble, fmt.Errorf("empty request body given"))

					Convey("And an empty patch array should be returned", func() {
						So(patches, ShouldBeEmpty)
					})
				})
			})
		})
	})

	Convey("Given a patch request with invalid patch request body", t, func() {
		body := strings.NewReader(`{}`)
		req := httptest.NewRequest(http.MethodPatch, "http://localhost:21800/jobs/12345", body)

		Convey("And valid supported patch operations given", func() {
			supportedOps := []PatchOp{OpAdd}

			Convey("When GetPatches is called", func() {
				patches, err := GetPatches(req.Body, supportedOps)

				Convey("Then an error should be returned ", func() {
					So(err, ShouldNotBeNil)
					So(err, ShouldResemble, fmt.Errorf("failed to unmarshal patch request body"))

					Convey("And an empty patch array should be returned", func() {
						So(patches, ShouldBeEmpty)
					})
				})
			})
		})
	})

	Convey("Given a patch request with no patch operation given in request body", t, func() {
		body := strings.NewReader(`[]`)
		req := httptest.NewRequest(http.MethodPatch, "http://localhost:21800/jobs/12345", body)

		Convey("And valid supported patch operations given", func() {
			supportedOps := []PatchOp{OpAdd}

			Convey("When GetPatches is called", func() {
				patches, err := GetPatches(req.Body, supportedOps)

				Convey("Then an error should be returned ", func() {
					So(err, ShouldNotBeNil)
					So(err, ShouldResemble, fmt.Errorf("no patches given in request body"))

					Convey("And an empty patch array should be returned", func() {
						So(patches, ShouldBeEmpty)
					})
				})
			})
		})
	})

	Convey("Given a patch request with unknown patch operation given in request body", t, func() {
		body := strings.NewReader(`[
			{ "op": "invalid", "path": "/a/b/c", "value": "foo" }
		]`)
		req := httptest.NewRequest(http.MethodPatch, "http://localhost:21800/jobs/12345", body)

		Convey("And valid supported patch operations given", func() {
			supportedOps := []PatchOp{OpAdd}

			Convey("When GetPatches is called", func() {
				patches, err := GetPatches(req.Body, supportedOps)

				Convey("Then an error should be returned ", func() {
					So(err, ShouldNotBeNil)

					supportedOpsStringSlice := getPatchOpsStringSlice(supportedOps)
					So(err, ShouldResemble, fmt.Errorf("patch operation is missing or invalid. Please, provide one of the following: %v", supportedOpsStringSlice))

					Convey("And an empty patch array should be returned", func() {
						So(patches, ShouldBeEmpty)
					})
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
		supportedOps := []PatchOp{OpAdd}
		So(patch.Validate(supportedOps), ShouldBeNil)
	})

	Convey("Validating a valid patch with a supported op and float64 value is successful", t, func() {
		patch := Patch{
			Op:    "add",
			Path:  "/a/b/c",
			Value: float64(123.321),
		}
		supportedOps := []PatchOp{OpAdd}
		So(patch.Validate(supportedOps), ShouldBeNil)
	})

	Convey("Validating a patch struct with an invalid op fails with the expected error", t, func() {
		patch := Patch{
			Op:    "wrong",
			Path:  "/a/b/c",
			Value: []string{"foo"},
		}
		emptySupportedOps := []PatchOp{}
		So(patch.Validate(emptySupportedOps), ShouldResemble, ErrInvalidOp(emptySupportedOps))
	})

	Convey("Validating a valid patch with an unsupported op fails with the expected error", t, func() {
		patch := Patch{
			Op:    "add",
			Path:  "/a/b/c",
			Value: []string{"foo"},
		}
		supportedOps := []PatchOp{OpRemove}
		So(patch.Validate(supportedOps), ShouldResemble, ErrUnsupportedOp("add", []PatchOp{OpRemove}))
	})

	Convey("Validating a patch struct with missing members for an operation results in the expected error being returned", t, func() {
		supportedOps := []PatchOp{OpAdd, OpReplace, OpTest, OpRemove, OpMove, OpCopy}
		patch := Patch{
			Op:   "add",
			Path: "/a/b/c",
		}
		So(patch.Validate(supportedOps), ShouldResemble, ErrMissingMember([]string{"value"}))
		patch = Patch{
			Op:    "replace",
			Value: []string{"foo"},
		}
		So(patch.Validate(supportedOps), ShouldResemble, ErrMissingMember([]string{"path"}))
		patch = Patch{
			Op: "test",
		}
		So(patch.Validate(supportedOps), ShouldResemble, ErrMissingMember([]string{"path", "value"}))
		patch = Patch{
			Op: "remove",
		}
		So(patch.Validate(supportedOps), ShouldResemble, ErrMissingMember([]string{"path"}))
		patch = Patch{
			Op:   "move",
			Path: "/a/b/c",
		}
		So(patch.Validate(supportedOps), ShouldResemble, ErrMissingMember([]string{"from"}))
		patch = Patch{
			Op:   "copy",
			From: "/c/b/a",
		}
		So(patch.Validate(supportedOps), ShouldResemble, ErrMissingMember([]string{"path"}))
	})
}
