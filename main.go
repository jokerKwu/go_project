package main

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/unrolled/render"
	j "go_project/jwt"
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
	e.POST("/token",j.PostAccessToken)
	//글작성 페이지 이동
	e.GET("/loginpage", user.GetLoginPageHandler)
	e.GET("/joinpage", user.GetJoinPageHandler)

	//Custom JWT
	r := e.Group("/posts",func(h echo.HandlerFunc) echo.HandlerFunc{
		return func(c echo.Context) error{
			//토큰 가져오고
			tokenString := c.Request().Header.Get("access_token")
			if tokenString == ""{
				return errors.New("토큰이 비어잇다.")
			}
			//여기서 토큰이 유효한지 체크
			token, err := jwt.Parse(tokenString, func(token *jwt.Token)(interface{}, error){
				claims := token.Claims.(jwt.MapClaims)
				fmt.Println(claims["exp"])
				fmt.Println(claims["userid"])
				if claims["userid"] != "test01"{
					return nil, errors.New("User isvalid")
				}
				return []byte("secret"), nil
			})
			//토큰이 유효하지 않다면 (만료시간 및 signature 체크)
			//리프레쉬 토큰을 요청한다.
			if err != nil{
				if err == jwt.ErrSignatureInvalid{
					return c.JSON(http.StatusUnauthorized,nil)
				}else{
					return c.JSON(http.StatusNotAcceptable, false)
				}
			}
			//컨텍스트에 사용자 아이디 저장
			c.Set("userid",token.Claims.(jwt.MapClaims)["userid"])
			return h(c)
		}
	})
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