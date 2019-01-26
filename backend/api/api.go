package api

import (
	"net/http"

	"github.com/jbrunsting/transient/backend/database"
)

type Api interface {
	SelfGet(w http.ResponseWriter, r *http.Request)
	UserGet(w http.ResponseWriter, r *http.Request)
	UserAuthenticatedGet(w http.ResponseWriter, r *http.Request)
	UserPost(w http.ResponseWriter, r *http.Request)
	UserLoginPost(w http.ResponseWriter, r *http.Request)
	UserLogoutPost(w http.ResponseWriter, r *http.Request)
	UserInvalidatePost(w http.ResponseWriter, r *http.Request)
	UserDeletePost(w http.ResponseWriter, r *http.Request)

	PostsGet(w http.ResponseWriter, r *http.Request)
	PostPost(w http.ResponseWriter, r *http.Request)
	PostDelete(w http.ResponseWriter, r *http.Request)

	FollowingsGet(w http.ResponseWriter, r *http.Request)
	FollowingPost(w http.ResponseWriter, r *http.Request)
	FollowingDelete(w http.ResponseWriter, r *http.Request)
}

type api struct {
	userApi
	postApi
	followingApi
}

func NewApi(db database.DatabaseHandler) Api {
    return &api{userApi: userApi{db: db}, postApi: postApi{db: db}, followingApi: followingApi{db: db}}
}
