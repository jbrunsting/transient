package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/jbrunsting/transient/backend/api"
	"github.com/jbrunsting/transient/backend/database"
)

type response struct {
	Message string `json:"message"`
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var err error
	r := mux.NewRouter()

	databaseHandler, err := database.NewDatabaseHandler()
	if err != nil {
		panic(err)
	}
	defer databaseHandler.Close()

	a := api.NewApi(databaseHandler)

	r.HandleFunc("/user", a.SelfGet).Methods("GET")
	r.HandleFunc("/user/{id}", a.UserGet).Methods("GET")
	r.HandleFunc("/user", a.UserPost).Methods("POST")
	r.HandleFunc("/user/login", a.UserLoginPost).Methods("POST")
	r.HandleFunc("/user/logout", a.UserLogoutPost).Methods("POST")
	r.HandleFunc("/user/invalidate", a.UserInvalidatePost).Methods("POST")
	r.HandleFunc("/user/delete", a.UserDeletePost).Methods("POST")
	r.HandleFunc("/user/password", a.UserPasswordPost).Methods("POST")
	r.HandleFunc("/users/search", a.UsersSearchGet).Methods("GET")
    r.HandleFunc("/users/exact/{username}", a.UsersExactGet).Methods("GET")
	r.HandleFunc("/authenticated", a.UserAuthenticatedGet).Methods("GET")

	r.HandleFunc("/posts/{id}", a.PostsGet).Methods("GET")
	r.HandleFunc("/post", a.PostPost).Methods("POST")
	r.HandleFunc("/post/{id}", a.PostDelete).Methods("DELETE")
	r.HandleFunc("/post/vote/{id}", a.PostVotePost).Methods("POST")
	r.HandleFunc("/post/{id}/comment", a.PostCommentPost).Methods("POST")
	r.HandleFunc("/post/{id}/comments", a.PostCommentsGet).Methods("GET")

	r.HandleFunc("/followings", a.FollowingsGet).Methods("GET")
	r.HandleFunc("/followings/posts", a.FollowingsPostsGet).Methods("GET")
	r.HandleFunc("/following/{id}", a.FollowingPost).Methods("POST")
	r.HandleFunc("/following/{id}", a.FollowingDelete).Methods("DELETE")

	r.HandleFunc("/recommends/posts", a.RecommendsPostsGet).Methods("GET")
	r.HandleFunc("/recommends/followings", a.RecommendsFollowingsGet).Methods("GET")

	log.Println("Listening on port 3000")
	http.ListenAndServe(":3000", r)
}
