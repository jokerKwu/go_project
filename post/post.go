package post

import (
	"context"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson"
	"go_project/mongodb"
	"net/http"
	"strconv"
)

//게시물 수정
func PutPostHandler(c echo.Context) (err error){
	post := new(mongodb.Post)
	if err = c.Bind(post); err != nil{
		c.Logger().Printf("PutPostHandler() - Bind Fail : " , post )
		return echo.NewHTTPError(http.StatusBadRequest,err.Error())
	}
	if err = c.Validate(post); err != nil{
		c.Logger().Printf("PutPostHandler() - Validate Fail : ",post)
		return echo.NewHTTPError(http.StatusBadRequest,err.Error())
	}
	mdb, err:= mongodb.GetClient()
	defer mdb.Disconnect(context.Background())
	if err != nil{
		return c.JSON(http.StatusInternalServerError, nil)
	}
	id, _ := strconv.Atoi(c.Param("id"))
	postUpdated := mongodb.UpdatePost(mdb, post, bson.M{"id":id})
	posts := mongodb.ReturnPostList(mdb,bson.M{})

	if postUpdated > 0{
		return c.Render(http.StatusOK,"index.html",posts)
	}else{
		return c.Render(http.StatusBadRequest,"index.html",posts)
	}
}

//게시물 제거
func DeletePostHandler(c echo.Context) error{
	id, _ :=strconv.Atoi(c.Param("id"))
	mdb, err := mongodb.GetClient()
	defer mdb.Disconnect(context.Background())
	if err != nil{
		return c.Render(http.StatusInternalServerError,"error.html",nil)
	}
	posts := mongodb.ReturnPostList(mdb,bson.M{})
	if postRemoved := mongodb.RemoveOnePost(mdb,bson.M{"id":id}); postRemoved > 0 {
		return c.Render(http.StatusOK,"index.html",posts)
	}else{
		return c.Render(http.StatusNotFound,"error.html",nil)
	}
}

//게시물 추가
func PostPostHandler(c echo.Context) (err error) {
	post := new(mongodb.Post)
	if err = c.Bind(post); err != nil{
		c.Logger().Printf("PostPostHandler() - Bind Fail : " , post )
		return echo.NewHTTPError(http.StatusBadRequest,err.Error())
	}
	if err = c.Validate(post); err != nil{
		c.Logger().Printf("PostPostHandler() - Validate Fail : ",post)
		return echo.NewHTTPError(http.StatusBadRequest,err.Error())
	}
	p := mongodb.Post{post.Id,post.Title,post.Content,post.Author,post.Date}
	mdb, err := mongodb.GetClient()
	defer mdb.Disconnect(context.Background())
	if err != nil{
		return c.Render(http.StatusInternalServerError,"error.html",nil)
	}
	insertId := mongodb.InsertNewPost(mdb,p)
	c.Logger().Print("post create complete!! : ", insertId)
	posts := mongodb.ReturnPostList(mdb,bson.M{})
	return c.Render(http.StatusOK,"index.html",posts)
}

//게시물 조회
func GetPostHandler(c echo.Context) error{
	c.Logger().Printf("GET으로 오나")
	mdb,err := mongodb.GetClient()
	defer mdb.Disconnect(context.Background())
	if err != nil{
		return c.Render(http.StatusInternalServerError,"error.html" ,nil)
	}
	id,_ := strconv.Atoi(c.Param("id"))
	if post := mongodb.ReturnPostOne(mdb, bson.M{"id":id}); post.Id == 0{
		return c.Render(http.StatusBadRequest,"error.html" ,nil)
	}else{
		return c.Render(http.StatusOK,"post_content.html",[]mongodb.Post{post})
	}
}

//게시물 리스트 조회
func GetPostListHandler(c echo.Context) error{
	mdb,err := mongodb.GetClient()
	defer mdb.Disconnect(context.Background())
	if err != nil{
		return c.Render(http.StatusInternalServerError,"error.html",nil)
	}
	posts := mongodb.ReturnPostList(mdb,bson.M{})
	return c.Render(http.StatusOK,"index.html",posts)
}

//게시글 작성 페이지 이동
func GetPostWriteHandler(c echo.Context) error{
	return c.Render(http.StatusOK,"post_write.html",nil)
}

//게시글 수정 페이지 이동
func GetPostUpdateHandler(c echo.Context) error{
	mdb,err := mongodb.GetClient()
	defer mdb.Disconnect(context.Background())
	if err != nil{
		return c.Render(http.StatusInternalServerError,"error.html" ,nil)
	}
	id,_ := strconv.Atoi(c.Param("id"))
	if post := mongodb.ReturnPostOne(mdb, bson.M{"id":id}); post.Id == 0{
		return c.Render(http.StatusBadRequest,"error.html" ,nil)
	}else{
		return c.Render(http.StatusOK,"post_update.html",[]mongodb.Post{post})
	}
}
