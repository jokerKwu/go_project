package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/unrolled/render"
	m "go_project/middleware"
	p "go_project/post"
	"go_project/user"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
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

	//권한이 필요하지 않는 핸들러
	// 로그인 및 회원가입
	e.GET("/", p.GetPostListHandler)
	e.POST("/login", user.PostLoginHandler)
	e.POST("/join", user.PostJoinHandler)
	//글작성 페이지 이동
	e.GET("/loginpage", user.GetLoginPageHandler)
	e.GET("/joinpage", user.GetJoinPageHandler)

	//Custom JWT
	r := e.Group("/posts",m.AuthToken)
	{	//권한이 필요한 핸들러
		r.GET("/:id", p.GetPostHandler)
		r.POST("", p.PostPostHandler)
		r.POST("/:id", p.PutPostHandler)
		r.DELETE("/:id", p.DeletePostHandler)
		r.GET("/write", p.GetPostWriteHandler)
		r.GET("/write/:id", p.GetPostUpdateHandler)
	}


	// server start
	e.Logger.Fatal(e.Start(":8080"))
}