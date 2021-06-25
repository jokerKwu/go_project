package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)
type Post struct{
		Id int	`json:"id" validate:"required"`
		Title string `json:"title" validate:"required"`
		Content string `json:"content" validate:"required"`
		Author string `json:"author" validate:"required"`
		Date string `json:"date" validate:"required"`
}
func GetClient() *mongo.Client{
	// 컨텍스트 수행시간 5초로 설정
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017/?connect=direct")
	client, err := mongo.NewClient(clientOptions)
	if err != nil{
		log.Println(err)
	}
	err = client.Connect(ctx)
	if err != nil{
		log.Println(err)
	}
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil{
		log.Println("Couldn't connect to the database",err)
	}else{
		log.Println("Connected!")
	}
	return client
}
//게시물 리스트 반환
func ReturnPostList(client *mongo.Client, filter bson.M) []*Post {
	ctx, _ := context.WithTimeout(context.Background(), 3 * time.Second)

	var posts []*Post
	collection := client.Database("webboard").Collection("posts")
	cur, err := collection.Find(ctx, filter)
	if err != nil{
		log.Println("Error on Finding all the documents", err)
	}
	for cur.Next(ctx){
		var post Post
		err = cur.Decode(&post)
		if err != nil{
			log.Println("Error on Decoding the document",err)
		}
		posts = append(posts, &post)
	}
	return posts
}
//특정 게시물 리턴
func ReturnPostOne(client *mongo.Client, filter bson.M) Post {
	var post Post
	collection := client.Database("webboard").Collection("posts")
	documentReturned := collection.FindOne(context.TODO(),filter)
	documentReturned.Decode(&post)
	return post
}
//게시물 생성
func InsertNewPost(client *mongo.Client, post Post) interface{}{
	collection := client.Database("webboard").Collection("posts")
	insertResult, err := collection.InsertOne(context.TODO(),post)
	if err != nil{
		log.Println("Error on inserting new post", err)
	}
	return insertResult.InsertedID
}

//게시물 삭제
func RemoveOnePost(client *mongo.Client, filter bson.M) int64{
	collection := client.Database("webboard").Collection("posts")
	deleteResult, err := collection.DeleteOne(context.TODO(),filter)
	if err != nil{
		log.Println("Error on deleting one post",err)
	}
	return deleteResult.DeletedCount
}
//게시물 수정
func UpdatePost(client *mongo.Client, updateData interface{}, filter bson.M) int64{
	collection := client.Database("webboard").Collection("posts")
	updateQuery := bson.D{{Key:"$set",Value:updateData}}
	updateResult, err := collection.UpdateOne(context.TODO(),filter, updateQuery)
	if err != nil{
		log.Println("Error on updating one post",err)
	}
	return updateResult.ModifiedCount
}