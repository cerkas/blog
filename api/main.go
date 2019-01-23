package main

import (
  "encoding/json"
  "fmt"
  "github.com/go-martini/martini"
  "io/ioutil"
  "log"
  "net/http"
  "os"
  "time"

  "gopkg.in/mgo.v2"

  "html/template"
)

type Post struct {
	Text      string    `json:"text" bson:"text"`
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
}

var posts *mgo.Collection

func main() {
	// Connect to mongo
	session, err := mgo.Dial("mongo:27017")
	if err != nil {
		log.Fatalln(err)
		log.Fatalln("mongo err")
		os.Exit(1)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	// Get posts collection
	posts = session.DB("app").C("posts")

	// Set up routes
	/*r := mux.NewRouter()
	r.HandleFunc("/posts", createPost).
		Methods("POST")
	r.HandleFunc("/posts", readPosts).
		Methods("GET")*/
  m := martini.Classic()
  //posts = make(map[string]*models.Post,0)
 /* m.Post("/assets/",http.StripPrefix("./web/assets/",http.FileServer(http.Dir("./assets"))))
  staticOptions := martini.StaticOptions{Prefix:"./web/assets"}
  m.Use(martini.Static("/web/assets",staticOptions))*/
  m.Get("/", indexHandler)
  m.Post("/posts",createPost)
	/*http.ListenAndServe(":8080", cors.AllowAll().Handler(r))
	log.Println("Listening on port 8080...")*/
	m.Run()
}
func indexHandler(w http.ResponseWriter, r *http.Request) {
 t, err := template.ParseFiles("web/templates/index.html","web/templates/header.html","web/templates/footer.html")
  if err !=nil {
    fmt.Printf( err.Error())
  }
  t.ExecuteTemplate(w,"index",nil)
}
func createPost(w http.ResponseWriter, r *http.Request) {
	// Read body
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Read post
	post := &Post{}
	err = json.Unmarshal(data, post)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	post.CreatedAt = time.Now().UTC()

	// Insert new post
	if err := posts.Insert(post); err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseJSON(w, post)
}

func readPosts(w http.ResponseWriter, r *http.Request) {
	result := []Post{}
	if err := posts.Find(nil).Sort("-created_at").All(&result); err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
	} else {
		responseJSON(w, result)
	}
}

func responseError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func responseJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
