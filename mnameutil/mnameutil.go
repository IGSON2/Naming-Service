package mnameutil

import (
	"bytes"
	"encoding/gob"
	"log"
)

func Errchk(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func Encode(anystruct interface{}) []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(anystruct)
	Errchk(err)
	return buffer.Bytes()
}
