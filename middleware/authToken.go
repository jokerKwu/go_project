package middleware

import (
	"errors"
	"github.com/labstack/echo"
	j "go_project/jwt"
	"go_project/utils"
	"net/http"
)
func AuthToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := utils.CtxGenerate(c.Request(), "", "")
		// get jwt Token
		accessToken := c.Request().Header.Get("access_token")
		if accessToken == "" {
			return c.JSON(http.StatusBadRequest,errors.New( "no access code in header"))
		}

		// verify & get Data
		tokenData, _, err := j.TokenVerifyAccess(ctx, accessToken, false)
		if err != nil {
			return c.JSON(http.StatusBadRequest, errors.New("token invalid"))
		}
		c.Set("userid",tokenData.UserID)
		return next(c)
	}
}
/*
func AuthToekn(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		//토큰 가져오고
		tokenString := c.Request().Header.Get("access_token")
		if tokenString == "" {
			return c.JSON(http.StatusBadRequest, errors.New("token empty"))
		}
		//여기서 토큰이 유효한지 체크
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			claims := token.Claims.(jwt.MapClaims)
			fmt.Println(claims["exp"])
			fmt.Println(claims["userid"])
			if claims["userid"] != "test01" {
				return c.JSON(http.StatusBadRequest, errors.New("user isvalid")), nil
			}
			return []byte("secret"), nil
		})
		//토큰이 유효하지 않다면 (만료시간 및 signature 체크)
		//리프레쉬 토큰을 요청한다.
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				return c.JSON(http.StatusUnauthorized, nil)
			} else {
				return c.JSON(http.StatusNotAcceptable, false)
			}
		}
		//컨텍스트에 사용자 아이디 저장
		c.Set("userid", token.Claims.(jwt.MapClaims)["userid"])
		return next(c)
	}
}

 */