package jwt

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go_project/utils"
	"time"
)

const(
	ACCESS_TOKEN_EXP = 1
	REFRESH_TOKEN_EXP = 24 * 7
)

//jwt session data
type TokenData struct{
	UserID string `json:"userid" validate:"required"`
}


type TokenMetaData struct{
	IssuedAt int64
	ExpiredAt int64
}

//claim 데이타
type jwtClaim struct{
	UID string `json:"uid" validate:"required"`
	jwt.StandardClaims
}


//토큰 생성 함수
func GenerateTokenPair(ctx context.Context, tokenData TokenData) (string, string,TokenMetaData, error){
	createTime := time.Now()
	//Create token
	accessToken, accessTokenMetaData, err := jwtCreate(ctx,tokenData,createTime,1)
	if err != nil{
		return "","",TokenMetaData{},err
	}

	refreshToken, _, err := jwtCreate(ctx,tokenData,createTime,REFRESH_TOKEN_EXP)
	if err != nil{
		return "","",TokenMetaData{},err
	}
	return accessToken,refreshToken,accessTokenMetaData,nil
}

func GenerateAccessToken(ctx context.Context,tokenData TokenData ) (string,TokenMetaData, error){
	createTime := time.Now()
	accessToken,accessTokenMetaData, err := jwtCreate(ctx,tokenData,createTime,ACCESS_TOKEN_EXP)
	if err != nil{
		return "",TokenMetaData{},err
	}
	return accessToken, accessTokenMetaData, nil
}
func GenerateRefreshToken(ctx context.Context,tokenData TokenData)(string,TokenMetaData,error){
	createTime := time.Now()
	refreshToken,refreshTokenMetaData, err := jwtCreate(ctx,tokenData,createTime,REFRESH_TOKEN_EXP)
	if err != nil{
		return "",TokenMetaData{},err
	}
	return refreshToken, refreshTokenMetaData, nil
}

func TokenVerifyBoth(ctx context.Context,accessTokenString string, refreshTokenString string, isForRefresh bool) (TokenData,TokenMetaData,TokenData,TokenMetaData,error){
	accessTokenData, accessTokenMetaData, err := TokenVerifyAccess(ctx, accessTokenString,isForRefresh)
	if err != nil{
		return TokenData{},TokenMetaData{},TokenData{},TokenMetaData{},err
	}
	refreshTokenData, refreshTokenMetaData, err := TokenVerifyAccess(ctx, refreshTokenString,isForRefresh)
	if err != nil{
		return TokenData{},TokenMetaData{},TokenData{},TokenMetaData{},err
	}
	if accessTokenData.UserID != refreshTokenData.UserID || accessTokenMetaData.IssuedAt != refreshTokenMetaData.IssuedAt{
		return TokenData{},TokenMetaData{},TokenData{},TokenMetaData{},err
	}
	return accessTokenData,accessTokenMetaData,refreshTokenData,refreshTokenMetaData,nil
}

func TokenVerifyAccess(ctx context.Context, tokenString string, isAllowExpire bool) (TokenData,TokenMetaData,error){
	tokenData, tokenMetaData, iErr := jwtVerify(ctx, tokenString, isAllowExpire)
	if iErr != nil {
		return TokenData{}, TokenMetaData{}, iErr
	}
	return tokenData, tokenMetaData, nil
}
func TokenVerifyRefresh(ctx context.Context, tokenString string) (TokenData,TokenMetaData,error){
	tokenData,tokenMetaData, err := jwtVerify(ctx, tokenString, false)
	if err != nil{
		return TokenData{},TokenMetaData{},err
	}
	return tokenData,tokenMetaData,nil
}

func jwtCreate(ctx context.Context,tokenData TokenData, createTime time.Time,expireHour time.Duration) (string,TokenMetaData,error){
	//validate struct
	if err := utils.ValidateStruct(tokenData) ;err != nil{
		return "",TokenMetaData{}, err
	}
	//create jwt
	tokenMetaData := TokenMetaData{
		IssuedAt: createTime.Unix(),
		ExpiredAt: createTime.Add(time.Hour * expireHour).Unix(),
	}
	claimStd := jwt.StandardClaims{
		IssuedAt: tokenMetaData.IssuedAt,
		ExpiresAt: tokenMetaData.ExpiredAt,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,jwtClaim{
		tokenData.UserID,
		claimStd,
	})

	tokenSigned, err := token.SignedString([]byte("secret"))
	if err != nil{
		return "",TokenMetaData{},err
	}
	return tokenSigned,tokenMetaData,nil
}
func jwtVerify(ctx context.Context,tokenString string, isAllowExpire bool) (TokenData,TokenMetaData,error){
	claims := &jwtClaim{}
	token, err := jwt.ParseWithClaims(tokenString,claims, func(token *jwt.Token)(interface{},error){
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok{
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		JWTkey := []byte("secret")
		return JWTkey , nil
	})
	//에러가 있거나 토큰이 유효하지 않는 경우
	if err != nil || !token.Valid{
		if err == nil{
			return TokenData{},TokenMetaData{},errors.New("jwt now valid")
		}
		if v, ok := err.(*jwt.ValidationError); !isAllowExpire || !ok || v.Errors != jwt.ValidationErrorExpired|| claims.ExpiresAt >= time.Now().Unix(){
			return TokenData{},TokenMetaData{},errors.New("jwt now valid")
		}
	}
	//validate claim format
	if err := utils.Val.Struct(claims); err != nil{
		return TokenData{},TokenMetaData{},err
	}
	tokenData := TokenData{
		UserID: claims.UID,
	}
	tokenMetaData := TokenMetaData{
		IssuedAt: claims.IssuedAt,
		ExpiredAt: claims.ExpiresAt,
	}
	return tokenData,tokenMetaData,nil
}