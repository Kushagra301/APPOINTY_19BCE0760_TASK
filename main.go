package main

// Importing the packages
import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/crypto/bcrypt"
)

// Defining the structure of the User and post.
type User struct {
	Id       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
}

type Post struct {
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID    string             `json:"UserID,omitempty" bson:"userID,omitempty"`
	Caption   string             `json:"caption,omitempty" bson:"caption,omitempty"`
	ImageURL  string             `json:"imageurl,omitempty" bson:"imageurl,omitempty"`
	Timestamp time.Time          `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
}

// Main function to be run
func main() {
	connectToDB()
	handleRequest()
}

// Code to connect to the MongoDB Database
var client *mongo.Client

func connectToDB() {
	atlasURI := "mongodb+srv://kushagra:1234@cluster0.wpm6r.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"
	clientOptions := options.Client().ApplyURI(atlasURI)

	client, _ = mongo.NewClient(clientOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//setup connection to MongoDB Atlas Cluster
	err := client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Connection Error: ", err)
	} else {
		log.Println("Connection Successful.")
	}

}

// function to handle all the requests
func handleRequest() {

	http.HandleFunc("/", root)
	http.HandleFunc("/posts", getAllPosts)
	http.HandleFunc("/newpost", newpost)
	http.HandleFunc("/users", getAllUsers)
	http.HandleFunc("/newuser", newuser)
	http.HandleFunc("/posts/", getpostbyID)
	http.HandleFunc("/users/", getuserbyID)
	http.HandleFunc("/posts/users/", getpostbyuserID)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}

// --------------------- CONNECTION ------------------------- //

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Connection Established.\n")
	fmt.Fprintf(w, "APPOINTY SUMMER INTERN TASK\n")
	fmt.Fprintf(w, "KUSHAGRA GUPTA\n")
	fmt.Fprintf(w, "19BCE0760\n")

}

// ----------------------- POSTS ---------------------------- //

// This function gets all the posts uploaded by a user by referencing its user-id
func getpostbyuserID(response http.ResponseWriter, request *http.Request) {

	if request.Method == "GET" {
		id := strings.TrimPrefix(request.URL.Path, "/posts/users/")
		var posts []Post
		collection := client.Database("test").Collection("Post")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		cursor, err := collection.Find(ctx, Post{UserID: id})
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var post Post
			cursor.Decode(&post)
			posts = append(posts, post)
		}
		if err = cursor.Err(); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		fmt.Println("Endpoint Hit: Post Article")
		json.NewEncoder(response).Encode(posts)
	}
}

// This function gets all the posts
func getAllPosts(response http.ResponseWriter, request *http.Request) {

	if request.Method == "GET" {
		var posts []Post
		collection := client.Database("test").Collection("Post")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var post Post
			cursor.Decode(&post)
			posts = append(posts, post)
		}
		if err = cursor.Err(); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		fmt.Println("Endpoint Hit: Post Article")
		json.NewEncoder(response).Encode(posts)
	}

}

// This is used to post a new post and it references the insertPost function
func newpost(response http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		request.ParseForm()
		decoder := json.NewDecoder(request.Body)
		var newPost Post

		newPost.Timestamp = time.Now()

		err := decoder.Decode(&newPost)
		if err != nil {
			panic(err)
		}
		log.Println(newPost.Id)
		fmt.Println("Endpoint Hit: Post Article")
		insertPost(newPost)
	}
}

//  This function is referenced by the newpost function and it inserts a new post to the database
func insertPost(post Post) {
	collection := client.Database("test").Collection("Post")
	insertResult, err := collection.InsertOne(context.TODO(), post)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Post Inserted with ID:", insertResult.InsertedID)
}

// This function is used to get all the posts by using their IDs.
func getpostbyID(response http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		id := strings.TrimPrefix(request.URL.Path, "/posts/")
		objID, _ := primitive.ObjectIDFromHex(id)
		var post Post
		collection := client.Database("test").Collection("Post")
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		err := collection.FindOne(ctx, Post{Id: objID}).Decode(&post)

		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		fmt.Println("Endpoint Hit: returnAllArticles")
		json.NewEncoder(response).Encode(post)
	}
}

// ---------------------------- USERS --------------------------------//

// This function gets all the posts uploaded by a user by referencing its user-id
func getAllUsers(response http.ResponseWriter, request *http.Request) {

	if request.Method == "GET" {
		var users []User
		collection := client.Database("test").Collection("User")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var user User
			cursor.Decode(&user)
			users = append(users, user)
		}
		if err = cursor.Err(); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		fmt.Println("Endpoint Hit: returnAllArticles")
		json.NewEncoder(response).Encode(users)
	}

}

// This is used to add a new user and it references the insertPost function
func newuser(response http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		request.ParseForm()
		decoder := json.NewDecoder(request.Body)

		var newUser User

		err := decoder.Decode(&newUser)
		if err != nil {
			panic(err)
		}

		log.Println(newUser.Id)
		fmt.Println("Endpoint Hit: Post Article")
		newUser.Password = hash(newUser.Password)
		insertUser(newUser)
	}
}

// This function converts the password entered by the user in the form of string and converts it into a hash function which cannot be reversed to get the password back.
func hash(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hash)
}

//  This function is referenced by the newuser function and it inserts a new post to the database
func insertUser(user User) {
	collection := client.Database("test").Collection("User")
	insertResult, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted user with ID:", insertResult.InsertedID)
}

// This function is used to get all the Users by using their IDs.
func getuserbyID(response http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		id := strings.TrimPrefix(request.URL.Path, "/users/")
		objID, _ := primitive.ObjectIDFromHex(id)
		var user User
		collection := client.Database("test").Collection("User")
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		err := collection.FindOne(ctx, User{Id: objID}).Decode(&user)

		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}

		fmt.Println("Endpoint Hit: Post Article")
		json.NewEncoder(response).Encode(user)
	}
}
