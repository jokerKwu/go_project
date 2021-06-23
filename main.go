package main

import (
	"encoding/json"
	mux2 "github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type Post struct{
	Id int	`json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
	Author string `json:"author"`
	Date string `json:"date"`
}
type Success struct{
	Success bool `json:"success"`
}
var posts map[int]Post	//게시물을 저장하는 맵형태 변수
var lastId int
type Posts []Post
var rd *render.Render
func MakeWebHandler() http.Handler{
	mux := mux2.NewRouter()
	mux.Handle("/",http.FileServer(http.Dir("public")))
	//핸들러 등록 1. 전체 조회 , 2.조회, 3. 작성, 4. 삭제, 5. 수정
	mux.HandleFunc("/posts",GetPostListHandler).Methods("GET")
	mux.HandleFunc("/posts/{id:[0-9]+}",GetPostHandler).Methods("GET")
	mux.HandleFunc("/posts",PostPostHandler).Methods("POST")
	mux.HandleFunc("/posts/{id:[0-9]+}",DeletePostHandler).Methods("DELETE")
	mux.HandleFunc("/posts/{id:[0-9]+}",PutPostHandler).Methods("PUT")

	posts = make(map[int]Post)
	posts[1] = Post{1,"title1","content1","Ryan",time.Now().Format("2006-01-02 15:05:05")}
	posts[2] = Post{2,"title2","content2","Ryan",time.Now().Format("2006-01-02 15:05:05")}

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
	var newPost Post
	err := json.NewDecoder(r.Body).Decode(&newPost)
	if err != nil{
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	vars := mux2.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	if post, ok := posts[id]; ok{
		post.Content = newPost.Content
		post.Title = newPost.Title
		post.Author = newPost.Author
		posts[id] = post
		rd.JSON(w, http.StatusOK, Success{true})
	}else{
		rd.JSON(w, http.StatusBadRequest,Success{false})
	}

}

//게시물 제거
func DeletePostHandler(w http.ResponseWriter, r *http.Request){
	vars := mux2.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	if _, ok := posts[id]; ok{
		delete(posts,id)
		rd.JSON(w,http.StatusOK,Success{true})
	}else{
		rd.JSON(w,http.StatusNotFound,Success{false})
	}
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
	rd.JSON(w,http.StatusCreated,post)
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
	rd.JSON(w, http.StatusOK, post)
}

// 게시물 리스트 조회
func GetPostListHandler(w http.ResponseWriter,r *http.Request){
	list := make(Posts, 0)
	for _, post := range posts{
		list = append(list, post)
	}
	sort.Sort(list)
	rd.JSON(w, http.StatusOK, list)
}

func main(){
	rd = render.New()
	m := MakeWebHandler()
	n := negroni.Classic()
	n.UseHandler(m)

	log.Println("Server start...")
	err := http.ListenAndServe(":8080",MakeWebHandler())
	if err != nil{
		panic(err)
	}
}
