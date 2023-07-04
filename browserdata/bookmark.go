package browserdata

import (
	"time"
)

type Bookmark struct {
	ID        int64
	Name      string
	Type      string
	URL       string
	DateAdded time.Time
}
