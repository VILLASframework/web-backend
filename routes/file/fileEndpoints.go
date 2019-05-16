package file

import (
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/json"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

func FilesRegister(r *gin.RouterGroup) {
	r.GET("/", filesReadEp)
	r.POST("/", fileRegistrationEp) // NEW in API
	r.PUT("/:FileID", fileUpdateEp) // NEW in API
	r.GET("/:FileID", fileReadEp)
	r.DELETE("/:FileID", fileDeleteEp)
}

func filesReadEp(c *gin.Context)  {
	// Database query
	allFiles, _, err := FindAllFiles()

	if common.ProvideErrorResponse(c, err) == false {
		serializer := FilesSerializerNoAssoc{c, allFiles}
		c.JSON(http.StatusOK, gin.H{
			"files": serializer.Response(),
		})
	}

}

func fileRegistrationEp(c *gin.Context) {
	var m map[string]interface{}

	decoder := json.NewDecoder(c.Request.Body)
	defer c.Request.Body.Close()

	if err := decoder.Decode(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad request. Invalid body.",
		})
		return;
	}

	// Database query
	err := AddFile(m)

	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK.",
		})
	}
}

func fileUpdateEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func fileReadEp(c *gin.Context) {
	var err error
	var file common.File
	fileID := c.Param("FileID")
	desc := c.GetHeader("X-Request-FileDesc")
	desc_b, _ := strconv.ParseBool(desc)

	userID := 1 // TODO obtain ID of user making the request

	//check if description of file or file itself shall be returned
	if desc_b {
		file, err = FindFile(userID, fileID)
		if common.ProvideErrorResponse(c, err) == false {
			serializer := FileSerializerNoAssoc{c, file}
			c.JSON(http.StatusOK, gin.H{
				"file": serializer.Response(),
			})
		}


	} else {
		//TODO: return file itself
	}
}

func fileDeleteEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}
