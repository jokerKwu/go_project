package jwt
/*
import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"net/http"
	"os"
)
type AuthApi struct {
}


// RefreshToken 토큰갱신: 성공하면 access 토큰과 refresh 토큰을 모두 다시 출력한다.
func (a *AuthApi) RefreshToken(c echo.Context) error {

	refreshToken := c.FormValue("refresh_token")

	// 토큰을 분석하고 유효한지 검증하고 *jwt.Token 객체로 반환한다.
	// KeyFunc 은 sceret key를 반환해야 한다.
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// 클라이언트로 받은 토큰이 HMAC 알고리즘이 맞는지 확인
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return err
	}

	// Claims 생성
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return echo.ErrInternalServerError
	}

	// 발급할때의 토큰제목과 일치하는지 체크
	if claims["sub"].(string) != os.Getenv("API_TOKEN_SUB") {
		return echo.ErrBadRequest
	}

	// 재발급
	cl := map[string]interface{}(claims)
	accessToken, refreshToken, err := h.generateToken(c, cl)
	if err != nil {
		return err
	}

	// 출력
	return c.JSON(http.StatusOK, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
*/