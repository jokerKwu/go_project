package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

//전체 조회
func TestJsonHandler(t *testing.T){
	assert := assert.New(t)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET","/posts",nil)

	mux := MakeWebHandler()
	mux.ServeHTTP(res, req)
	assert.Equal(http.StatusOK,res.Code)

	var list []Post

	err := json.NewDecoder(res.Body).Decode(&list) // JSON 데이터를 list로 변환한다.
	assert.Nil(err)	//이렇게 변환한 객체의 값이다.
	assert.Equal("Ryan",list[0].Author)
	assert.Equal(1,list[0].Id)
	assert.Equal("title1",list[0].Title)
	assert.Equal("content1",list[0].Content)

	assert.Equal("Ryan",list[1].Author)
	assert.Equal(2,list[1].Id)
	assert.Equal("title2",list[1].Title)
	assert.Equal("content2",list[1].Content)
}

//특정 게시물 조회
func TestJsonHandler2(t *testing.T){
	assert := assert.New(t)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET","/posts/1",nil)

	mux := MakeWebHandler()
	mux.ServeHTTP(res, req)
	assert.Equal(http.StatusOK, res.Code)

	var post Post
	err := json.NewDecoder(res.Body).Decode(&post)
	assert.Nil(err)
	assert.Equal(1,post.Id)
	assert.Equal("title1",post.Title)
	assert.Equal("content1",post.Content)
	assert.Equal("Ryan",post.Author)

}
//추가 테스트
func TestJsonHandler3(t *testing.T){
	assert := assert.New(t)
	mux := MakeWebHandler()

	res := httptest.NewRecorder()
	req := httptest.NewRequest("POST","/posts",strings.NewReader(`{"Id":0,"Author":"Ryan","Title":"title3","Content":"content3"}`))

	mux.ServeHTTP(res, req)
	//assert.Equal(http.StatusOK, res.Code)

	res = httptest.NewRecorder()
	req = httptest.NewRequest("GET","/posts/3",nil)

	mux.ServeHTTP(res, req)
	assert.Equal(http.StatusOK,res.Code)
	var post Post
	err := json.NewDecoder(res.Body).Decode(&post)
	assert.Nil(err)
	assert.Equal("content3",post.Content)

}
//삭제
func TestJsonHandler4(t *testing.T){
	assert := assert.New(t)
	mux := MakeWebHandler()
	res := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE","/posts/1",nil)
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusOK,res.Code)
	res = httptest.NewRecorder()
	req = httptest.NewRequest("GET","/posts",nil)
	mux.ServeHTTP(res, req)
	assert.Equal(http.StatusOK, res.Code)
	var list []Post
	err := json.NewDecoder(res.Body).Decode(&list)
	assert.Nil(err)
	assert.Equal(1,len(list))
	assert.Equal("title2",list[0].Title)

}