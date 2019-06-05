package core

import (
	"time"
)

type Directories struct {
	Path	string
	Files	[]Files
}

type Files struct {
	Name	string
}

type Stats struct {
	StartedAt  time.Time
	FinishedAt time.Time
}
