package user

import (
	"context"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson"
	j "go_project/jwt"
	"go_project/mongodb"
	"net/http"
	"time"
)

//로그인 페이지 이동
func GetLoginPageHandler(c echo.Context) error{
	return c.Render(http.StatusOK,"login.html",nil)
}

//회원가입 페이지 이동
func GetJoinPageHandler(c echo.Context) error{
	return c.Render(http.StatusOK,"join.html",nil)
}

//로그인 핸들러
func PostLoginHandler(c echo.Context) (err error) {
	//1.요청으로부터 받은 사용자 정보를 디비 체킹한다.
	user := new(mongodb.User)
	if err = c.Bind(user); err != nil{
		c.Logger().Printf("USER Bind Fail : " , user )
		return echo.NewHTTPError(http.StatusBadRequest,err.Error())
	}
	if err = c.Validate(user); err != nil{
		c.Logger().Printf("User Validate Fail : ",user)
		return echo.NewHTTPError(http.StatusBadRequest,err.Error())
	}
	mdb,err := mongodb.GetClient()
	defer mdb.Disconnect(context.Background())
	if err != nil{
		return c.Render(http.StatusInternalServerError,"error.html" ,nil)
	}
	dbUser := mongodb.ReturnUserOne(mdb,bson.M{"userid":user.Userid})
	//디비 체킹한다.
	if user.Userid != dbUser.Userid || user.Password != dbUser.Password {
		return echo.ErrUnauthorized
	}
	//2.액세스,리프레시 토큰을 생성한다.
	tokens, err := j.GenerateTokenPair(user.Userid)
	if err != nil{
		return err
	}

	// 토큰 쿠키에 저장
	access_cookie := new(http.Cookie)
	access_cookie.Name = "access_token"
	access_cookie.Value = tokens["access_token"]
	access_cookie.Expires = time.Now().Add(time.Hour * 10)

	refresh_cookie := new(http.Cookie)
	refresh_cookie.Name = "refresh_token"
	refresh_cookie.Value = tokens["refresh_token"]
	refresh_cookie.Expires = time.Now().Add(time.Hour * 24*7)

	//3. 클라이언트에게 엑세스, 리프레쉬 토큰을 발급해준다.
	c.SetCookie(access_cookie)
	c.SetCookie(refresh_cookie)

	posts := mongodb.ReturnPostList(mdb,bson.M{})
	return c.Render(http.StatusOK,"index.html",posts)
}


//회원가입 핸들러
func PostJoinHandler(c echo.Context) (err error){
	//회원가입 처리
	user := new(mongodb.User)
	if err = c.Bind(user); err != nil{
		c.Logger().Printf("PostPostHandler() - Bind Fail : " , user )
		return echo.NewHTTPError(http.StatusBadRequest,err.Error())
	}
	if err = c.Validate(user); err != nil{
		c.Logger().Printf("PostPostHandler() - Validate Fail : ",user)
		return echo.NewHTTPError(http.StatusBadRequest,err.Error())
	}
	u := mongodb.User{user.Userid,user.Password}
	mdb, err := mongodb.GetClient()
	defer mdb.Disconnect(context.Background())
	if err != nil{
		return c.Render(http.StatusInternalServerError,"error.html",nil)
	}
	insertId := mongodb.InsertNewUser(mdb,u)
	c.Logger().Print("User create complete!! : ", insertId)
	return c.Render(http.StatusOK,"login.html",u)

	return nil
}

