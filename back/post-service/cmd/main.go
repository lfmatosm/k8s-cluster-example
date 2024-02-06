package main

import (
	"encoding/json"
	"fmt"
	"io"
	"lfmatosm/k8s-cluster-example/back/post-service/pkg/controllers"
	"lfmatosm/k8s-cluster-example/back/post-service/pkg/repositories"
	"lfmatosm/k8s-cluster-example/back/post-service/pkg/services"
	"lfmatosm/k8s-cluster-example/back/post-service/pkg/utils"
	"net/http"
	"os"
	"time"
)

var repository repositories.Repository
var service *services.PostService
var controller *controllers.PostController

func init() {
	var MONGODB_URI = os.Getenv("MONGODB_URI")
	var MONGODB_DATABASE = os.Getenv("MONGODB_DATABASE")

	repository = repositories.NewMongoRepository(MONGODB_URI, MONGODB_DATABASE)
	service = services.NewPostService(repository, "post")
	controller = controllers.NewPostController(service)
}

func logRequest(req *http.Request, resp *utils.HttpResponse) {
	fmt.Printf("[%s] %s %d %s %s\n", time.Now().UTC().String(), req.RemoteAddr, resp.Status, req.Method, req.RequestURI)
}

func getResponse(req *http.Request) *utils.HttpResponse {
	if req.Method == "OPTIONS" {
		return utils.Ok(nil)
	} else if req.Method == "GET" {
		return controller.List()
	} else if req.Method == "POST" {
		if req.Body == nil {
			return utils.BadRequest(map[string]string{"error": "Empty body received"})
		}

		defer req.Body.Close()

		bytes, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Printf("Error reading body: %s\n", err.Error())
			return utils.InternalServerError(map[string]interface{}{"error": err.Error()})
		}

		return controller.Save(req.Header, bytes)
	}
	fmt.Printf("Method not allowed: %s\n", req.Method)
	return utils.MethodNotAllowed()
}

func writeResponse(w http.ResponseWriter, resp *utils.HttpResponse) {
	for k, v := range resp.Headers {
		w.Header().Set(k, v)
	}

	w.WriteHeader(resp.Status)

	if resp.Body == nil {
		return
	}

	bytes, err := json.Marshal(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(bytes))
}

func handle(w http.ResponseWriter, req *http.Request) {
	var resp *utils.HttpResponse = getResponse(req)
	writeResponse(w, resp)
	logRequest(req, resp)
}

func main() {
	http.HandleFunc("/posts", handle)
	http.ListenAndServe(":8090", nil)
}
