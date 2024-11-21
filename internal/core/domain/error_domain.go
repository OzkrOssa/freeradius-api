package domain

import "errors"

var (
	DataNotFoundError               = errors.New("data not found")
	InternalError                   = errors.New("internal server error")
	ConflictDataError               = errors.New("data conflicts with existing data")
	NoUpdatedDataError              = errors.New("no data to update")
	TokenDurationError              = errors.New("invalid token duration format")
	TokenCreationError              = errors.New("error creating token")
	ExpiredTokenError               = errors.New("access token has expired")
	InvalidTokenError               = errors.New("access token is invalid")
	InvalidCredentialsError         = errors.New("invalid email or password")
	EmptyAuthorizationHeaderError   = errors.New("authorization header is empty")
	InvalidAuthorizationHeaderError = errors.New("invalid authorization header")
	InvalidAuthorizationTypeError   = errors.New("invalid authorization type")
)
