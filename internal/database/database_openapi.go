package database

import (
	"bytes"
	"encoding/gob"

	"github.com/getkin/kin-openapi/openapi3"
)

func (db *Db) InsertAPI(api *openapi3.T) error {
	var output bytes.Buffer
	encoder := gob.NewEncoder(&output)
	err := encoder.Encode(api)
	if err != nil {
		return err
	}
	return db.Create(&API{OpenAPI: output.Bytes()}).Error
}
