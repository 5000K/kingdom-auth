package core

import "errors"

var ErrFailedMigration = errors.New("failed migration")
var ErrUnknownDbDriver = errors.New("unknown database driver")
