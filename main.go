package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/gorilla/mux"
)

//  ! global variable
var BaseURL = "https://api.hatchways.io/assessment/blog/posts?tag="
var posts PostsList

type Post struct {
	Author     string   `json:"author"`
	AuthorID   int      `json:"authorId"`
	ID         int      `json:"id"`
	Likes      int      `json:"likes"`
	Popularity float64  `json:"popularity"`
	Reads      int      `json:"reads"`
	Tags       []string `json:"tags"`
}

type PostsList struct {
	Posts []Post `json:"posts"`
}

func homePage(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr
	xff := r.Header.Get("X-Forwarded-For")
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
	fmt.Println("IP: ", ip)
	fmt.Println("X-Forwarded-For: ", xff)
}

func getPosts(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	tag := queryParams["tag"]
	sort := queryParams["sortBy"]
	direction := queryParams["direction"]
	// fmt.Println(tag)
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
		// fmt.Println(resp.Body)
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
	// print type of posts
	// fmt.Println(reflect.TypeOf(posts))
	if sort != nil {
		if len(direction) == 1 {
			sortBy(posts.Posts, sort[0], direction[0])
		} else {
			sortBy(posts.Posts, sort[0], "asc")
		}
	}
	json.NewEncoder(w).Encode(posts)

}

func contains(s []Post, e Post) bool {
	for _, a := range s {
		if a.ID == e.ID {
			return true
		}
	}
	return false
}

func sortBy(s []Post, field string, direction string) {
	if direction == "asc" {
		var less func(i, j int) bool
		switch field {
		case "id":
			less = func(i, j int) bool { return s[i].ID < s[j].ID }
		case "reads":
			less = func(i, j int) bool { return s[i].Reads < s[j].Reads }
		case "likes":
			less = func(i, j int) bool { return s[i].Likes < s[j].Likes }
		case "popularity":
			less = func(i, j int) bool { return s[i].Popularity < s[j].Popularity }
		}
		sort.Slice(s, less)
	} else if direction == "desc" {
		var less func(i, j int) bool
		switch field {
		case "id":
			less = func(i, j int) bool { return s[i].ID > s[j].ID }
		case "reads":
			less = func(i, j int) bool { return s[i].Reads > s[j].Reads }
		case "likes":
			less = func(i, j int) bool { return s[i].Likes > s[j].Likes }
		case "popularity":
			less = func(i, j int) bool { return s[i].Popularity > s[j].Popularity }
		}
		sort.Slice(s, less)
	}
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/posts", getPosts)

	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	fmt.Println("Rest API using hatchways API")
	handleRequests()

}
