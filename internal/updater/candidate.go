package updater

import (
	"time"
)

type ReleaseCandidate interface {
	Version() string
	Branch() Branch
	Description() string
	DownloadUrl() string
	CreatedAt() time.Time
	IsPrerelease() bool
}
