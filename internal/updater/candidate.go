package updater

import (
	"time"
)

type ReleaseCandidate interface {
	Version() string
	Description() string
	DownloadUrl() string
	CreatedAt() time.Time
	IsPrerelease() bool
}
