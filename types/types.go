package types

import "mime/multipart"

type SetData struct {
	File     *multipart.FileHeader `form:"fileData" binding:"required"`
	Keywords string                `form:"keywords"`
}

type ParametersType struct {
	Datatype string
	Value    string
	Array    []string
	Bytes32  [32]byte
}

type SetDataOutput struct {
	Hash      string
	Timestamp string
}
