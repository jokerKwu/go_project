package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
)
type Post struct{
		Id int	`json:"id" validate:"required"`
		Title string `json:"title" validate:"required"`
		Content string `json:"content" validate:"required"`
		Author string `json:"author" validate:"required"`
		Date string `json:"date" validate:"required"`
}
func GetClient() *mongo.Client{
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017/?connect=direct")
	client, err := mongo.NewClient(clientOptions)
	if err != nil{
		log.Fatal(err)
	}
	err = client.Connect(context.Background())
	if err != nil{
		log.Fatal(err)
	}
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil{
		log.Fatal("Couldn't connect to the database",err)
	}else{
		log.Println("Connected!")
	}
	return client
}
//게시물 리스트 반환
func ReturnPostList(client *mongo.Client, filter bson.M) []*Post {
	var posts []*Post
	collection := client.Database("webboard").Collection("posts")
	cur, err := collection.Find(context.TODO(), filter)
	if err != nil{
		log.Fatal("Error on Finding all the documents", err)
	}
	for cur.Next(context.TODO()){
		var post Post
		err = cur.Decode(&post)
		if err != nil{
			log.Fatal("Error on Decoding the document",err)
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
		log.Fatalln("Error on inserting new post", err)
	}
	return insertResult.InsertedID
}

//게시물 삭제
func RemoveOnePost(client *mongo.Client, filter bson.M) int64{
	collection := client.Database("webboard").Collection("posts")
	deleteResult, err := collection.DeleteOne(context.TODO(),filter)
	if err != nil{
		log.Fatal("Error on deleting one post",err)
	}
	return deleteResult.DeletedCount
}
//게시물 수정
func UpdatePost(client *mongo.Client, updateData interface{}, filter bson.M) int64{
	collection := client.Database("webboard").Collection("posts")
	updateQuery := bson.D{{Key:"$set",Value:updateData}}
	updateResult, err := collection.UpdateOne(context.TODO(),filter, updateQuery)
	if err != nil{
		log.Fatal("Error on updating one post",err)
	}
	return updateResult.ModifiedCount
}