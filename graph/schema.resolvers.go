package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.64

import (
	"backoffice/db"
	"backoffice/graph/model"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, input model.NewPost) (*model.Post, error) {
	counters := r.Resolver.MongoClient.Client.Database("db").Collection("counters")

	newIndex, err := db.GetNextCollectionIndex(counters, "blogPostID")
	if err != nil {
		log.Fatalf("Error getting next sequence: %v", err)
	}

	post := &model.Post{
		ID:        fmt.Sprintf("%d", newIndex),
		Published: false,
		Title:     input.Title,
		Text:      input.Text,
	}

	bsonPost, err := bson.Marshal(post)
	if err != nil {
		log.Fatalf("Error marshalling to BSON: %v", err)
	}

	collection := r.Resolver.MongoClient.Client.Database("db").Collection("blog")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, bson.Raw(bsonPost))
	if err != nil {
		log.Fatalf("DB Insert Failed: %v", err)
	}

	return post, nil
}

// Attach is the resolver for the attach field.
func (r *mutationResolver) Attach(ctx context.Context, files []string) (string, error) {
	panic(fmt.Errorf("not implemented: Attach - attach"))
}

// EditPost is the resolver for the editPost field.
func (r *mutationResolver) EditPost(ctx context.Context, input model.EditPost) (*model.Post, error) {
	log.Println("Updating blog post")
	collection := r.Resolver.MongoClient.Client.Database("db").Collection("blog")

	filter := bson.M{"id": input.ID}
	update := bson.M{"$set": bson.M{}}

	// Conditionally add fields to update only if they are provided
	if input.Title != "" {
		update["$set"].(bson.M)["title"] = input.Title
	}
	if input.Text != "" {
		update["$set"].(bson.M)["text"] = input.Text
	}
	if input.Published != false {
		update["$set"].(bson.M)["published"] = input.Published
	}

	// Execute update operation
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedPost model.Post
	err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedPost)
	if err != nil {
		return nil, err
	}

	return &updatedPost, nil
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context) ([]*model.Post, error) {
	// Use MongoDB client from the Resolver struct
	collection := r.Resolver.MongoClient.Client.Database("db").Collection("blog")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatalf("DB Lookup Failed: %v", err)
	}

	var posts []*model.Post

	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var post model.Post
		if err := cur.Decode(&post); err != nil {
			log.Fatal(err)
		}

		posts = append(posts, &post)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return posts, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
