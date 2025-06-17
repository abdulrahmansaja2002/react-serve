package errors

import "errors"

var (
	ErrConnectionFailed      = errors.New("connect to Elasticsearch failed")
	ErrNoIdFieldInMapData    = errors.New("no id field in map data")
	ErrNoIdFieldInStructData = errors.New("no id field in struct data")
	ErrNoIdFieldInData       = errors.New("failed to get id field in data")
)
