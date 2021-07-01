package main

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
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
	"time"
)

// jwtCustomClaims are custom claims extending default ones.
type jwtCustomClaims struct {
	Userid string `json:"userid"`
	Admin  bool   `json:"admin"`
	jwt.StandardClaims
}

type CustomValidator struct {
	validator *validator.Validate
}

type Success struct {
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

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

var rd *render.Render

//게시물 수정
func PutPostHandler(c echo.Context) (err error) {
	post := new(mongodb.Post)
	if err = c.Bind(post); err != nil {
		c.Logger().Printf("PutPostHandler() - Bind Fail : ", post)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err = c.Validate(post); err != nil {
		c.Logger().Printf("PutPostHandler() - Validate Fail : ", post)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	mdb, err := mongodb.GetClient()
	defer mdb.Disconnect(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	id, _ := strconv.Atoi(c.Param("id"))
	postUpdated := mongodb.UpdatePost(mdb, post, bson.M{"id": id})
	posts := mongodb.ReturnPostList(mdb, bson.M{})

	if postUpdated > 0 {
		return c.Render(http.StatusOK, "index.html", posts)
	} else {
		return c.Render(http.StatusBadRequest, "index.html", posts)
	}
}

//게시물 제거
func DeletePostHandler(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	mdb, err := mongodb.GetClient()
	defer mdb.Disconnect(context.Background())
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", Success{false})
	}
	posts := mongodb.ReturnPostList(mdb, bson.M{})
	if postRemoved := mongodb.RemoveOnePost(mdb, bson.M{"id": id}); postRemoved > 0 {
		return c.Render(http.StatusOK, "index.html", posts)
	} else {
		return c.Render(http.StatusNotFound, "error.html", Success{false})
	}
}

//게시물 추가
func PostPostHandler(c echo.Context) (err error) {
	post := new(mongodb.Post)
	if err = c.Bind(post); err != nil {
		c.Logger().Printf("PostPostHandler() - Bind Fail : ", post)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err = c.Validate(post); err != nil {
		c.Logger().Printf("PostPostHandler() - Validate Fail : ", post)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	p := mongodb.Post{post.Id, post.Title, post.Content, post.Author, post.Date}
	mdb, err := mongodb.GetClient()
	defer mdb.Disconnect(context.Background())
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", Success{false})
	}
	insertId := mongodb.InsertNewPost(mdb, p)
	c.Logger().Print("post create complete!! : ", insertId)
	posts := mongodb.ReturnPostList(mdb, bson.M{})
	return c.Render(http.StatusOK, "index.html", posts)
}

//게시물 조회

func GetPostHandler(c echo.Context) error {
	c.Logger().Printf("GET으로 오나")
	mdb, err := mongodb.GetClient()
	defer mdb.Disconnect(context.Background())
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", nil)
	}
	id, _ := strconv.Atoi(c.Param("id"))
	if post := mongodb.ReturnPostOne(mdb, bson.M{"id": id}); post.Id == 0 {
		return c.Render(http.StatusBadRequest, "error.html", nil)
	} else {
		return c.Render(http.StatusOK, "post_content.html", []mongodb.Post{post})
	}
}

//게시물 리스트 조회
func GetPostListHandler(c echo.Context) error {
	mdb, err := mongodb.GetClient()
	defer mdb.Disconnect(context.Background())
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", nil)
	}
	posts := mongodb.ReturnPostList(mdb, bson.M{})
	return c.Render(http.StatusOK, "index.html", posts)
}

func GetPostWriteHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "post_write.html", nil)
}

//게시글 수정 페이지 이동
func GetPostUpdateHandler(c echo.Context) error {
	mdb, err := mongodb.GetClient()
	defer mdb.Disconnect(context.Background())
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", nil)
	}
	id, _ := strconv.Atoi(c.Param("id"))
	if post := mongodb.ReturnPostOne(mdb, bson.M{"id": id}); post.Id == 0 {
		return c.Render(http.StatusBadRequest, "error.html", nil)
	} else {
		return c.Render(http.StatusOK, "post_write.html", []mongodb.Post{post})
		id, _ := strconv.Atoi(c.Param("id"))
		if post := mongodb.ReturnPostOne(mdb, bson.M{"id": id}); post.Id == 0 {
			return c.Render(http.StatusBadRequest, "error.html", nil)
		} else {
			return c.Render(http.StatusOK, "post_update.html", []mongodb.Post{post})
		}
	}
}

//로그인 페이지 이동
func GetLoginPageHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", nil)
}

//회원가입 페이지 이동
func GetJoinPageHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "join.html", nil)
}

//로그인 핸들러
func PostLoginHandler(c echo.Context) (err error) {
	//1.요청으로부터 받은 사용자 정보를 디비 체킹한다.
	user := new(mongodb.User)
	if err = c.Bind(user); err != nil {
		c.Logger().Printf("USER Bind Fail : ", user)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err = c.Validate(user); err != nil {
		c.Logger().Printf("User Validate Fail : ", user)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	mdb, err := mongodb.GetClient()
	defer mdb.Disconnect(context.Background())
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", nil)
	}
	dbUser := mongodb.ReturnUserOne(mdb, bson.M{"userid": user.Userid})
	//디비 체킹한다.
	if user.Userid != dbUser.Userid || user.Password != dbUser.Password {
		return echo.ErrUnauthorized
	}
	//2.액세스,리프레시 토큰을 생성한다.
	tokens, err := generateTokenPair()
	if err != nil {
		return err
	}

	// 토큰 쿠키에 저장
	access_cookie := new(http.Cookie)
	access_cookie.Name = "access_token"
	access_cookie.Value = tokens["access_token"]
	access_cookie.Expires = time.Now().Add(time.Hour * 1)

	refresh_cookie := new(http.Cookie)
	refresh_cookie.Name = "refresh_token"
	refresh_cookie.Value = tokens["refresh_token"]
	refresh_cookie.Expires = time.Now().Add(time.Hour * 24)

	//3. 클라이언트에게 엑세스, 리프레쉬 토큰을 발급해준다.
	c.SetCookie(access_cookie)
	c.SetCookie(refresh_cookie)

	posts := mongodb.ReturnPostList(mdb, bson.M{})
	return c.Render(http.StatusOK, "index.html", posts)
}

//리프레시 토큰 생성 핸들러
func token(c echo.Context) error {
	type tokenReqBody struct {
		RefreshToken string `json:"refresh_token"`
	}
	tokenReq := tokenReqBody{}
	c.Bind(&tokenReq)

	// Parse takes the token string and a function for looking up the key.
	// The latter is especially useful if you use multiple keys for your application.
	// The standard is to use 'kid' in the head of the token to identify
	// which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenReq.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("secret"), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Get the user record from database or
		// run through your business logic to verify if the user can log in
		if int(claims["sub"].(float64)) == 1 {

			newTokenPair, err := generateTokenPair()
			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, newTokenPair)
		}

		return echo.ErrUnauthorized
	}

	return err
}

//토큰 생성 함수
func generateTokenPair() (map[string]string, error) {
	//Create token
	token := jwt.New(jwt.SigningMethodHS256)

	//Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = 1
	claims["userid"] = "ryan"
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return nil, err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = 1
	rtClaims["userid"] = "ryan"
	rtClaims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()
	rtClaims["admin"] = true

	rt, err := refreshToken.SignedString([]byte("secret"))
	if err != nil {
		return nil, err
	}

	return map[string]string{"access_token": t, "refresh_token": rt}, nil
}

//회원가입 핸들러
func PostJoinHandler(c echo.Context) (err error) {
	//회원가입 처리
	user := new(mongodb.User)
	if err = c.Bind(user); err != nil {
		c.Logger().Printf("PostPostHandler() - Bind Fail : ", user)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err = c.Validate(user); err != nil {
		c.Logger().Printf("PostPostHandler() - Validate Fail : ", user)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	u := mongodb.User{user.Userid, user.Password}
	mdb, err := mongodb.GetClient()
	defer mdb.Disconnect(context.Background())
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", Success{false})
	}
	insertId := mongodb.InsertNewUser(mdb, u)
	c.Logger().Print("User create complete!! : ", insertId)
	return c.Render(http.StatusOK, "login.html", u)

	return nil
}

func main() {
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
	e.Static("/static/", "public")
	//CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	e.Static("/static/", "public")
	e.Renderer = t

	e.GET("/", GetPostListHandler)
	e.GET("/posts", GetPostListHandler)
	e.GET("/posts/:id", GetPostHandler)
	e.POST("/posts", PostPostHandler)
	e.PUT("/posts/:id", PutPostHandler)
	e.DELETE("/posts/:id", DeletePostHandler)
	//Custom JWT
	e.POST("/token", token)
	r := e.Group("/posts")
	r.Use(middleware.JWTWithConfig(
		middleware.JWTConfig{
			SigningMethod: "HS256",
			SigningKey:    []byte("secret"),
			TokenLookup:   "cookie:access_token",
		}))

	//권한이 필요하지 않는 핸들러
	// 로그인 및 회원가입
	e.POST("/login", PostLoginHandler)
	e.POST("/join", PostJoinHandler)
	// server start
	e.Logger.Fatal(e.Start(":8080"))
}
