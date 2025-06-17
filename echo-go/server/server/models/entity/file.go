package entity

import "time"

type File struct {
	Name      string    `bson:"name"`
	Path      string    `bson:"path"`
	Size      int       `bson:"size"`
	Type      string    `bson:"type"`
	CreatedAt time.Time `bson:"created_at"`
}
