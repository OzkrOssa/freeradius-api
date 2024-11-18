package domain

import "errors"

var (
	DataNotFoundError  = errors.New("data not found")
	InternalError      = errors.New("internal server error")
	ConflictDataError  = errors.New("data conflicts with existing data")
	NoUpdatedDataError = errors.New("no data to update")
)
