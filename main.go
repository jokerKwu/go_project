package main

import (
	"encoding/json"
	"fmt"
	mux2 "github.com/gorilla/mux"
	"net/http"
	"sort"
	"strconv"
)

type Post struct{
	Id int
	Title string
	Content string
	Author string
}

var posts map[int]Post	//게시물을 저장하는 변수
var lastId int
type Posts []Post

func MakeWebHandler() http.Handler{
	mux := mux2.NewRouter()
	//핸들러 등록 1. 전체 조회 , 2.조회, 3. 작성, 4. 삭제, 5. 수정
	mux.HandleFunc("/posts",GetPostListHandler).Methods("GET")
	mux.HandleFunc("/posts/{id:[0-9]+}",GetPostHandler).Methods("GET")
	mux.HandleFunc("/posts",PostPostHandler).Methods("POST")
	mux.HandleFunc("/posts/{id:[0-9]+}",DeletePostHandler).Methods("DELETE")
	mux.HandleFunc("/posts/{id:[0-9]+}",PutPostHandler).Methods("PUT")

	posts = make(map[int]Post)
	posts[1] = Post{1,"title1","content1","Ryan"}
	posts[2] = Post{2,"title2","content2","Ryan"}

	lastId = 2
	return mux
}

func (p Posts) Len() int{
	return len(p)
}
func (p Posts) Swap(i, j int){
	p[i], p[j] = p[i], p[i]
}
func (p Posts) Less(i, j int) bool{
	return p[i].Id < p[j].Id
}
//게시물 수정
func PutPostHandler(w http.ResponseWriter, r *http.Request){
	vars := mux2.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	_, ok := posts[id]
	if !ok{
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

//게시물 제거
func DeletePostHandler(w http.ResponseWriter, r *http.Request){
	vars := mux2.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	_, ok := posts[id]
	if !ok{
		w.WriteHeader(http.StatusNotFound)
		return
	}
	delete(posts, id)
	w.WriteHeader(http.StatusOK)
}

//게시물 추가
func PostPostHandler(w http.ResponseWriter,r *http.Request){
	var post Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	lastId++
	post.Id = lastId
	posts[lastId] = post
	w.WriteHeader(http.StatusCreated)
}

//게시물 조회
func GetPostHandler(w http.ResponseWriter, r *http.Request){
	vars := mux2.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	post, ok := posts[id]

	if !ok{
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(post)
}

// 게시물 리스트 조회
func GetPostListHandler(w http.ResponseWriter,r *http.Request){
	list := make(Posts, 0)
	for _, post := range posts{
		list = append(list, post)
	}
	sort.Sort(list)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(list)
}

func main(){
	fmt.Println("서버 시작...")
	http.ListenAndServe(":3000",MakeWebHandler())
}
