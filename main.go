package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/unrolled/render"
	"go.mongodb.org/mongo-driver/bson"
	mongodb "go_project/mongodb"
	"net/http"
	"sort"
	"strconv"
)
type (
	Post struct{
		Id int	`json:"id" validate:"required"`
		Title string `json:"title" validate:"required"`
		Content string `json:"content" validate:"required"`
		Author string `json:"author" validate:"required"`
		Date string `json:"date" validate:"required"`
	}
	CustomValidator struct{
		validator *validator.Validate
	}
)
func (cv *CustomValidator) Validate(i interface{}) error{
	if err := cv.validator.Struct(i); err != nil{
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return nil
}
type Success struct{
	Success bool `json:"success"`
}
type Posts []Post
var rd *render.Render
func (p Posts) Len() int{
	return len(p)
}
func (p Posts) Swap(i, j int){
	p[i], p[j] = p[j], p[i]
}
func (p Posts) Less(i, j int) bool{
	return p[i].Id < p[j].Id
}
//게시물 수정
func PutPostHandler(c echo.Context) (err error){
	post := new(Post)
	if err = c.Bind(post); err != nil{
		c.Logger().Printf("PutPostHandler() - Bind Fail : " , post )
		return echo.NewHTTPError(http.StatusBadRequest,err.Error())
	}
	if err = c.Validate(post); err != nil{
		c.Logger().Printf("PutPostHandler() - Validate Fail : ",post)
		return echo.NewHTTPError(http.StatusBadRequest,err.Error())
	}
	mdb := mongodb.GetClient()
	//post 에 값이 들어갔다.
	id, _ := strconv.Atoi(c.Param("id"))
	postUpdated := mongodb.UpdatePost(mdb, post, bson.M{"id":id})
	if postUpdated > 0{
		return c.JSON(http.StatusOK,Success{true})
	}else{
		return c.JSON(http.StatusBadRequest,Success{false})
	}
}

//게시물 제거
func DeletePostHandler(c echo.Context) error{
	id, _ :=strconv.Atoi(c.Param("id"))

	mdb := mongodb.GetClient()
	postRemoved := mongodb.RemoveOnePost(mdb,bson.M{"id":id})

	if postRemoved > 0{
		return c.JSON(http.StatusOK, Success{true})
	}else{
		return c.JSON(http.StatusNotFound,Success{false})
	}
}

//게시물 추가
func PostPostHandler(c echo.Context) (err error) {
	//var post Post
	post := new(Post)

	if err = c.Bind(post); err != nil{
		c.Logger().Printf("PostPostHandler() - Bind Fail : " , post )
		return echo.NewHTTPError(http.StatusBadRequest,err.Error())
	}
	if err = c.Validate(post); err != nil{
		c.Logger().Printf("PostPostHandler() - Validate Fail : ",post)
		return echo.NewHTTPError(http.StatusBadRequest,err.Error())
	}

	p := mongodb.Post{post.Id,post.Title,post.Content,post.Author,post.Date}
	mdb := mongodb.GetClient()
	insertId := mongodb.InsertNewPost(mdb,p)
	c.Logger().Print("생성 완료 : ", insertId)
	return c.JSON(http.StatusOK, p)
}
//게시물 조회
func GetPostHandler(c echo.Context) error{
	mdb := mongodb.GetClient()
	id,_ := strconv.Atoi(c.Param("id"))
	post := mongodb.ReturnPostOne(mdb, bson.M{"id":id})
	if post.Id == 0{
		return c.JSON(http.StatusBadRequest,nil)
	}else {
		return c.JSON(http.StatusOK, post)
	}
}
// 게시물 리스트 조회
func GetPostListHandler(c echo.Context) error{
	mdb := mongodb.GetClient()
	posts := mongodb.ReturnPostList(mdb,bson.M{})
	//게시물 데이터 가져와서 정렬 후 리스트 전달
	list := make(Posts, 0)
	for _, post := range posts{
		p := Post{
			post.Id,
			post.Title,
			post.Content,
			post.Author,
			post.Date,
		}
		list = append(list, p)
	}
	sort.Sort(list)
	return c.JSON(http.StatusOK, list)
}
func main(){
	//Echo Instance create
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET,echo.HEAD,echo.PUT,echo.PATCH,echo.POST,echo.DELETE},
	}))

	// Route -> handler register
	// 1. 전체 조회 2. 조회 3. 생성 4. 수정 5. 삭제
	e.GET("/posts", GetPostListHandler)
	e.GET("/posts/:id",GetPostHandler)
	e.POST("/posts",PostPostHandler)
	e.PUT("/posts/:id",PutPostHandler)
	e.DELETE("/posts/:id",DeletePostHandler)

	// server start
	e.Logger.Fatal(e.Start(":8080"))
}