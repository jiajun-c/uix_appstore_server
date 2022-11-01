package handler

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name     string `json:"name" bson:"name"`
	Password string `json:"password" bson:"password"`
	Phone    string `json:"phone" bson:"phone"`
}

var Client *mongo.Client

var SECRET_KEY = []byte("gosecretkey")

func getHash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}
func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString(SECRET_KEY)
	if err != nil {
		log.Println("Error in JWT token generation")
		return "", err
	}
	return tokenString, nil
}
func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("Key: ", k)
		fmt.Println("Value: ", v)
	}
	fmt.Fprintf(w, "Hello astac")
}

func UserSignup(response http.ResponseWriter, request *http.Request) {
	fmt.Println("[LOG INFO] login up")
	response.Header().Set("Content-Type", "application/json")
	var user User
	json.NewDecoder(request.Body).Decode(&user)
	fmt.Println("The user name is ", user.Name)
	user.Password = getHash([]byte(user.Password))
	collection := Client.Database("GODB").Collection("user")
	ctx, err := context.WithTimeout(context.Background(), 10*time.Second)
	if err != nil {
		fmt.Println(err)
	}
	result, _ := collection.InsertOne(ctx, user)
	json.NewEncoder(response).Encode(result)
}

func UserLogin(response http.ResponseWriter, request *http.Request) {
	fmt.Println("[LOG INFO] login in")
	response.Header().Set("Content-Type", "application/json")
	var user User
	var dbUser User
	json.NewDecoder(request.Body).Decode(&user)
	collection := Client.Database("GODB").Collection("user")
	ctx, err1 := context.WithTimeout(context.Background(), 10*time.Second)
	if err1 != nil {
		fmt.Println("Timeout error:")
	}
	err := collection.FindOne(ctx, bson.M{"phone": user.Phone}).Decode(&dbUser)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	userPass := []byte(user.Password)
	dbPass := []byte(dbUser.Password)

	passErr := bcrypt.CompareHashAndPassword(dbPass, userPass)

	if passErr != nil {
		log.Println(passErr)
		response.Write([]byte(`{"response":"Wrong Password!"}`))
		return
	}
	jwtToken, err := GenerateJWT()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	response.Write([]byte(`{"token":"` + jwtToken + `"}`))
}

func Upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Method: ", r.Method)
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))
		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, token)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		// fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666) // 此处假设当前目录下已存在test目录
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}
