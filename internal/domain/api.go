package domain

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type API struct {
	Id         string      `bson:"_id" json:"id"`
	Spec       *openapi3.T `json:"spec" bson:",inline"`
	Selections []Selection `json:"selections,omitempty"`
}
type Selection struct {
	Path     string `json:"path"`
	Method   string `json:"method"`
	Selected bool   `json:"selected"`
}

type CustomizableT struct {
	Paths CustomizedPaths `json:"paths,omitempty" bson:",omitempty"`
}
type CustomizableAPI struct {
	ApiId      string         `bson:"_ref" json:"id"`
	Spec       *CustomizableT `json:"spec" bson:",inline"`
	Selections []Selection    `json:"selections,omitempty" bson:",omitempty"`
	Username   string         `json:"username,omitempty"`
}
type CustomizedPaths map[string]*PathItem

type PathItem struct {
	Delete  *Operation `json:"delete,omitempty" bson:",omitempty"`
	Get     *Operation `json:"get,omitempty" bson:",omitempty"`
	Head    *Operation `json:"head,omitempty" bson:",omitempty"`
	Options *Operation `json:"options,omitempty" bson:",omitempty"`
	Patch   *Operation `json:"patch,omitempty" bson:",omitempty"`
	Post    *Operation `json:"post,omitempty" bson:",omitempty"`
	Put     *Operation `json:"put,omitempty" bson:",omitempty"`
	Trace   *Operation `json:"trace,omitempty" bson:",omitempty"`
}
type Operation struct {
	Description *string `json:"description" bson:",omitempty"`
}
