package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type PostsList struct {
	Posts []struct {
		Author     string   `json:"author"`
		AuthorID   int      `json:"authorId"`
		ID         int      `json:"id"`
		Likes      int      `json:"likes"`
		Popularity float64  `json:"popularity"`
		Reads      int      `json:"reads"`
		Tags       []string `json:"tags"`
	} `json:"posts"`
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/posts", getPosts)

	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

//  ! global variable
var BaseURL = "https://api.hatchways.io/assessment/blog/posts?tag="
var posts PostsList

func main() {
	fmt.Println("Rest API using hatchways API")
	handleRequests()

}

func getPosts(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	tag := queryParams["tag"]
	fmt.Println(tag)
	var tempPosts PostsList
	// now we do a GET request to the API
	// we use the tag variable to get the posts
	// we use the BaseURL variable to get the posts
	if len(tag) == 1 {
		fmt.Println(BaseURL + tag[0])
		resp, err := http.Get(BaseURL + tag[0])
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()
		// we use the json.Decoder to decode the response
		fmt.Println(resp.Body)
		err = json.NewDecoder(resp.Body).Decode(&tempPosts)
		if err != nil {
			fmt.Println(err)
		}
		posts = tempPosts
	} else if len(tag) > 1 {
		for _, value := range tag {
			resp, err := http.Get(BaseURL + value)
			if err != nil {
				fmt.Println(err)
			}
			defer resp.Body.Close()
			// we use the json.Decoder to decode the response

			err = json.NewDecoder(resp.Body).Decode(&tempPosts)
			if err != nil {
				fmt.Println(err)
			}
			// we check if author id is already in the posts array
			// if not we append the author id to the posts array
			for _, post := range tempPosts.Posts {
				if !contains(posts.Posts, post) {
					posts.Posts = append(posts.Posts, post)
				}
			}
		}
	}

	json.NewEncoder(w).Encode(posts)

}

func contains(s []struct {
	Author     string   `json:"author"`
	AuthorID   int      `json:"authorId"`
	ID         int      `json:"id"`
	Likes      int      `json:"likes"`
	Popularity float64  `json:"popularity"`
	Reads      int      `json:"reads"`
	Tags       []string `json:"tags"`
}, e struct {
	Author     string   `json:"author"`
	AuthorID   int      `json:"authorId"`
	ID         int      `json:"id"`
	Likes      int      `json:"likes"`
	Popularity float64  `json:"popularity"`
	Reads      int      `json:"reads"`
	Tags       []string `json:"tags"`
}) bool {
	for _, a := range s {
		if a.AuthorID == e.AuthorID {
			return true
		}
	}
	return false
}
