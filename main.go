package main

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/unrolled/render"
	"go.mongodb.org/mongo-driver/bson"
	mongodb "go_project/mongodb"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)
type CustomValidator struct{
		validator *validator.Validate
	}

type Success struct{
	Success bool `json:"success"`
}

type Template struct {
	templates *template.Template
}
//GetTempFilesFromFolders is scans file path
func GetTempFilesFromFolders(folders []string) []string {
	var filepaths []string
	for _, folder := range folders {
		files, _ := ioutil.ReadDir(folder)

		for _, file := range files {
			if strings.Contains(file.Name(), ".html") {
				filepaths = append(filepaths, folder+file.Name())
			}
		}
	}
	return filepaths
}
//Render is function to use template function
/*
 template 기능을 사용하기 위한 함수
 w : http.status
 name : html 명
 data : 전달하고자 하는 데이터
*/
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}


func (cv *CustomValidator) Validate(i interface{}) error{
	if err := cv.validator.Struct(i); err != nil{
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return nil
}
var rd *render.Render


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
		return c.Render(http.StatusInternalServerError,"error.html",Success{false})
	}
	posts := mongodb.ReturnPostList(mdb,bson.M{})
	if postRemoved := mongodb.RemoveOnePost(mdb,bson.M{"id":id}); postRemoved > 0 {
		return c.Render(http.StatusOK,"index.html",posts)
	}else{
		return c.Render(http.StatusNotFound,"error.html",Success{false})
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
		return c.Render(http.StatusInternalServerError,"error.html",Success{false})
	}
	insertId := mongodb.InsertNewPost(mdb,p)
	c.Logger().Print("post create complete!! : ", insertId)
	posts := mongodb.ReturnPostList(mdb,bson.M{})
	return c.Render(http.StatusOK,"index.html",posts)
}
//게시물 조회
func GetPostHandler(c echo.Context) error{
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
// 게시물 리스트 조회
func GetPostListHandler(c echo.Context) error{
	mdb,err := mongodb.GetClient()
	defer mdb.Disconnect(context.Background())
	if err != nil{
		return c.Render(http.StatusInternalServerError,"error.html",nil)
	}
	posts := mongodb.ReturnPostList(mdb,bson.M{})
	return c.Render(http.StatusOK,"index.html",posts)
}

func GetPostWriteHandler(c echo.Context) error{
	//아이디 체크.
	//...
	return c.Render(http.StatusOK,"post_write.html",nil)
}
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
		return c.Render(http.StatusOK,"post_write.html",[]mongodb.Post{post})
	}
}

func main(){
	dirs := []string{"./public/", "./public/static/include/"}
	tempfiles := GetTempFilesFromFolders(dirs)
	t := &Template{
		templates: template.Must(template.ParseFiles(tempfiles...)),
	}

	//Echo Instance create
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/static/","public")
	//CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET,echo.HEAD,echo.PUT,echo.PATCH,echo.POST,echo.DELETE},
	}))
	e.Static("/static/", "public")
	e.Renderer = t

	/*
		Handler Register
		1. 전체 조회 2. 조회 3. 생성 4. 수정 5. 삭제
	 */
	e.GET("/",GetPostListHandler)
	e.GET("/posts", GetPostListHandler)
	e.GET("/posts/:id",GetPostHandler)
	e.POST("/posts",PostPostHandler)
	e.PUT("/posts/:id",PutPostHandler)
	e.DELETE("/posts/:id",DeletePostHandler)

	//글작성 페이지 이동
	e.GET("/posts/write",GetPostWriteHandler)
	e.GET("/posts/write/:id",GetPostUpdateHandler)
	// server start
	e.Logger.Fatal(e.Start(":8080"))
}