package domain

import "time"

type Image struct {
	ID           int
	FilePath     string
	OriginalName string
	ContentType  string
	ByteSize     int64
	OwnerID      int
	CreatedAt    time.Time
	DeletedAt    *time.Time
}
