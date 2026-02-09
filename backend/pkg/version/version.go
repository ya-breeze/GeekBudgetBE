package version

import "time"

// These variables are populated at runtime.
var (
	// BuildTime is the time when the application was built.
	// It's extracted from Go's build info (vcs.time).
	BuildTime string

	// Commit is the git commit hash from which the application was built.
	// It's extracted from Go's build info (vcs.revision).
	Commit string

	// StartTime is the time when the application started.
	StartTime time.Time
)

// Init initializes the version information from build metadata.
func Init() {
	StartTime = time.Now()
}
