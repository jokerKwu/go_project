package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sort"
	"time"
)

type Post struct {
	Id      int       `json:"id" bson:"id" validate:"required"`
	Title   string    `json:"title" bson:"title" validate:"required"`
	Content string    `json:"content" bson:"content" validate:"required"`
	Author  string    `json:"author" bson:"author" validate:"required"`
	Date    time.Time `json:"date" bson:"date" validate:"omitempty"'`
}
type User struct {
	Userid   string `json:"userid" bson:"userid" validate:"required"`
	Password string `json:"password" bson:"password" validate:"required"`
}
type PostID struct {
	Seq int `json:"seq" validate:"required"`
}
type Posts []Post

func (p Posts) Len() int {
	return len(p)
}
func (p Posts) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p Posts) Less(i, j int) bool {
	return p[i].Id > p[j].Id
}

func GetClient() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017/?connect=direct"))
	if err != nil {
		log.Println(err)
	}
	return client, err
}

//게시물 리스트 반환
func ReturnPostList(client *mongo.Client, filter bson.M) []Post {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var posts Posts
	collection := client.Database("webboard").Collection("posts")
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		log.Println("Error on Finding all the documents", err)
	}
	for cur.Next(ctx) {
		var post Post
		err = cur.Decode(&post)
		if err != nil {
			log.Println("Error on Decoding the document", err)
		}
		posts = append(posts, post)
	}
	sort.Sort(posts)
	return posts
}

//특정 게시물 리턴
func ReturnPostOne(client *mongo.Client, filter bson.M) Post {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var post Post
	collection := client.Database("webboard").Collection("posts")
	documentReturned := collection.FindOne(ctx, filter)
	documentReturned.Decode(&post)
	if err := documentReturned.Decode(&post); err != nil {
		log.Println(err)
	}
	return post
}

//게시물 생성
func InsertNewPost(client *mongo.Client, post Post) interface{} {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := client.Database("webboard").Collection("posts")
	post.Date = time.Now()

	var postid PostID
	collection2 := client.Database("webboard").Collection("counters")
	res := collection2.FindOne(ctx, bson.M{})
	res.Decode(&postid)
	post.Id = postid.Seq + 1
	collection2.UpdateOne(ctx, bson.M{"id": "postid"}, bson.M{"$set": bson.M{"seq": postid.Seq + 1}})

	insertResult, err := collection.InsertOne(ctx, post)
	if err != nil {
		log.Println("Error on inserting new post", err)
	}
	return insertResult.InsertedID
}

//게시물 삭제
func RemoveOnePost(client *mongo.Client, filter bson.M) int64 {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := client.Database("webboard").Collection("posts")
	deleteResult, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Println("Error on deleting one post", err)
	}
	return deleteResult.DeletedCount
}

//게시물 수정
func UpdatePost(client *mongo.Client, updateData interface{}, filter bson.M) int64 {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := client.Database("webboard").Collection("posts")
	updateQuery := bson.D{{Key: "$set", Value: updateData}}
	updateResult, err := collection.UpdateOne(ctx, filter, updateQuery)
	if err != nil {
		log.Println("Error on updating one post", err)
	}
	return updateResult.ModifiedCount
}

//게시물 생성
func InsertNewUser(client *mongo.Client, user User) interface{} {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := client.Database("webboard").Collection("users")
	insertResult, err := collection.InsertOne(ctx, user)
	if err != nil {
		log.Println("Error on inserting new post", err)
	}
	return insertResult.InsertedID
}

//유저 아이디 조회
func ReturnUserOne(client *mongo.Client, filter bson.M) User {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var user User
	collection := client.Database("webboard").Collection("users")
	documentReturned := collection.FindOne(ctx, filter)
	if err := documentReturned.Decode(&user); err != nil {
		log.Println(err)
	}
	return user
}
