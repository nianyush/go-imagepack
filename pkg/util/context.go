package util

import (
	"context"

	"github.com/twpayne/go-vfs/v4"
)

type Contenxt struct {
	Ctx context.Context
	FS vfs.FS
}