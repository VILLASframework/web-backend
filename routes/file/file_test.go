package file

import (
	"bytes"
	"encoding/json"
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/user"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var token string
var filecontent = "This is my testfile"
var filecontent_update = "This is my updated testfile with a dot at the end."
var filename = "testfile.txt"
var filename_update = "testfileupdate.txt"

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var cred = credentials{
	Username: "User_A",
	Password: "abc123",
}

var msgOK = common.ResponseMsg{
	Message: "OK.",
}

var FileA = common.FileResponse{
	Name:     "File_A",
	ID:       1,
	WidgetID: 1,
}

var FileB = common.FileResponse{
	Name:     "File_B",
	ID:       2,
	WidgetID: 1,
}

var myFiles = []common.FileResponse{
	FileA,
	FileB,
}

var msgFiles = common.ResponseMsgFiles{
	Files: myFiles,
}

// Test /models endpoints
func TestSignalEndpoints(t *testing.T) {

	db := common.DummyInitDB()
	defer db.Close()
	common.DummyPopulateDB(db)

	// create a testfile in local folder
	c1 := []byte(filecontent)
	c2 := []byte(filecontent_update)
	err := ioutil.WriteFile(filename, c1, 0644)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(filename_update, c2, 0644)
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	api := router.Group("/api")

	// All endpoints require authentication except when someone wants to
	// login (POST /authenticate)
	user.VisitorAuthenticate(api.Group("/authenticate"))

	api.Use(user.Authentication(true))

	RegisterFileEndpoints(api.Group("/files"))

	credjson, err := json.Marshal(cred)
	if err != nil {
		panic(err)
	}

	msgOKjson, err := json.Marshal(msgOK)
	if err != nil {
		panic(err)
	}

	msgFilesjson, err := json.Marshal(msgFiles)
	if err != nil {
		panic(err)
	}

	token = common.AuthenticateForTest(t, router, "/api/authenticate", "POST", credjson, 200)

	// test GET files
	common.TestEndpoint(t, router, token, "/api/files?objectID=1&objectType=widget", "GET", nil, 200, string(msgFilesjson))

	// test POST files
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, err := bodyWriter.CreateFormFile("file", "testuploadfile.txt")
	if err != nil {
		fmt.Println("error writing to buffer")
		panic(err)
	}

	// open file handle
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		panic(err)
	}
	defer fh.Close()

	// io copy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		fmt.Println("error on IO copy")
		panic(err)
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/api/files?objectID=1&objectType=widget", bodyBuf)
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", contentType)
	if err != nil {
		fmt.Println("error creating post request")
		panic(err)
	}

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	fmt.Println(w.Body.String())
	assert.Equal(t, string(msgOKjson), w.Body.String())

	// test GET files/:fileID
	common.TestEndpoint(t, router, token, "/api/files/5", "GET", nil, 200, filecontent)
	//common.TestEndpoint(t, router, token, "/api/files?objectID=1&objectType=widget", "GET", nil, 200, string(msgFilesjson))

	// test PUT files/:fileID
	bodyBuf_update := &bytes.Buffer{}
	bodyWriter_update := multipart.NewWriter(bodyBuf_update)
	fileWriter_update, err := bodyWriter_update.CreateFormFile("file", "testuploadfile.txt")
	if err != nil {
		fmt.Println("error writing to buffer")
		panic(err)
	}

	// open file handle
	fh_update, err := os.Open(filename_update)
	if err != nil {
		fmt.Println("error opening file")
		panic(err)
	}
	defer fh_update.Close()

	// io copy
	_, err = io.Copy(fileWriter_update, fh_update)
	if err != nil {
		fmt.Println("error on IO copy")
		panic(err)
	}

	contentType_update := bodyWriter_update.FormDataContentType()
	bodyWriter_update.Close()
	w_update := httptest.NewRecorder()
	req_update, err := http.NewRequest("PUT", "/api/files/5", bodyBuf_update)
	req_update.Header.Add("Authorization", "Bearer "+token)
	req_update.Header.Set("Content-Type", contentType_update)
	if err != nil {
		fmt.Println("error creating post request")
		panic(err)
	}

	router.ServeHTTP(w_update, req_update)

	assert.Equal(t, 200, w_update.Code)
	fmt.Println(w_update.Body.String())
	assert.Equal(t, string(msgOKjson), w_update.Body.String())

	common.TestEndpoint(t, router, token, "/api/files/5", "GET", nil, 200, filecontent_update)

	// test DELETE files/:fileID
	common.TestEndpoint(t, router, token, "/api/files/5", "DELETE", nil, 200, string(msgOKjson))
	common.TestEndpoint(t, router, token, "/api/files?objectID=1&objectType=widget", "GET", nil, 200, string(msgFilesjson))

	// TODO add testing for other return codes

	// clean up temporary file
	err = os.Remove(filename)
	if err != nil {
		panic(err)
	}

	err = os.Remove(filename_update)
	if err != nil {
		panic(err)
	}

}
