package core

import "errors"

var ErrFailedMigration = errors.New("failed migration")
var ErrUnknownDbDriver = errors.New("unknown database driver")
var ErrTokenInvalid = errors.New("invalid token")
var ErrTokenExpired = errors.New("token expired")
var ErrFailedToParseToken = errors.New("failed to parse token")
var ErrInvalidSignature = errors.New("invalid signature")
