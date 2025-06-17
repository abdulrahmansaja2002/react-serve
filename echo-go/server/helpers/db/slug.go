package db

import (
	"context"
	"echo-react-serve/config"
	"fmt"
	"log"

	slugify "github.com/gosimple/slug"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func GenerateSlug(collection, title string) string {
	slug := slugify.MakeLang(title, "id")
	coll := config.MongoDB.Collection(collection)
	count, err := coll.CountDocuments(context.TODO(), bson.M{"slug": slug})
	if err != nil {
		log.Println("[SLUG] find count error: ", err)
		return slug
	}
	if count > 0 {
		slug = fmt.Sprintf("%s-%d", slug, count)
	}
	return slug
}
