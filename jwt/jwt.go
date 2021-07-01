package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const(
	ACCESS_TOKEN_EXP = 1			// HOUR
	REFRESH_TOKEN_EXP = 24 * 7
)

// RefreshToken 토큰갱신: 성공하면 access 토큰과 refresh 토큰을 모두 다시 출력한다.
func PostAccessToken(c echo.Context) error {
	accessToken :=c.Request().Header.Get("access_token")
	refreshToken := c.Request().Header.Get("refresh_token")
	aToken, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		// 클라이언트로 받은 토큰이 HMAC 알고리즘이 맞는지 확인
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return err
	}
	// 토큰을 분석하고 유효한지 검증하고 *jwt.Token 객체로 반환한다.
	// KeyFunc 은 sceret key를 반환해야 한다.
	rToken, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// 클라이언트로 받은 토큰이 HMAC 알고리즘이 맞는지 확인
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return err
	}

	// Claims 생성
	aClaims, ok := aToken.Claims.(jwt.MapClaims)
	rClaims, ok := rToken.Claims.(jwt.MapClaims)

	//처음 만들때 true로 설정하기 때문에 확인 가능
	if !ok || !rToken.Valid || !aToken.Valid{
		if !rToken.Valid{
			fmt.Println(rToken.Valid)
		}else {
			return echo.ErrInternalServerError
		}
	}
	// 발급할때의 토큰제목과 일치하는지 체크
	if aClaims["sub"].(string) != rClaims["sub"].(string) {
		return echo.ErrBadRequest
	}
	// 재발급
	tokens, err := GenerateTokenPair(rClaims["userid"].(string))
	if err != nil {
		return err
	}
	// 출력
	return c.JSON(http.StatusOK, map[string]string{
		"access_token":  tokens["access_token"],
		"refresh_token": tokens["refresh_token"],
	})
}

//토큰 생성 함수
func GenerateTokenPair(name string) (map[string]string, error){
	timeSource := rand.NewSource(time.Now().UnixNano())
	random := rand.New(timeSource)
	value := random.Intn(100000000)

	//Create token
	token := jwt.New(jwt.SigningMethodHS256)

	//Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = name + strconv.Itoa(value)
	claims["userid"] = name
	claims["admin"] = true
//	claims["exp"] = time.Now().Add(time.Hour * ACCESS_TOKEN_EXP).Unix()
	claims["exp"] = time.Now().Add(time.Second * 30).Unix()
	fmt.Println(claims["exp"])
	t, err := token.SignedString([]byte("secret"))

	if err != nil{
		return nil, err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = name + strconv.Itoa(value)
	rtClaims["exp"] = time.Now().Add(time.Hour * REFRESH_TOKEN_EXP).Unix()
	rtClaims["admin"] = true
	rtClaims["userid"] = name
	rt, err := refreshToken.SignedString([]byte("secret"))

	if err != nil{
		return nil, err
	}

	return map[string]string{"access_token": t, "refresh_token":rt},nil
}
