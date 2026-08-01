package minio

import "io"

type UploadReq struct {
	BucketName  string
	ObjectName  string
	Data        io.Reader
	Size        int64
	ContentType string
	DataType    string
}

type UploadRes struct {
}
