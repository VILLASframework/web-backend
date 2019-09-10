package helper

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetIDOfElement(c *gin.Context, elementName string, source string, providedID int) (int, error) {

	if source == "path" {
		id, err := strconv.Atoi(c.Param(elementName))
		if err != nil {
			BadRequestError(c, fmt.Sprintf("No or incorrect format of path parameter"))
			return -1, err
		}
		return id, nil
	} else if source == "query" {
		id, err := strconv.Atoi(c.Request.URL.Query().Get(elementName))
		if err != nil {
			BadRequestError(c, fmt.Sprintf("No or incorrect format of query parameter"))
			return -1, err
		}
		return id, nil
	} else if source == "body" {
		id := providedID
		return id, nil
	} else {
		return -1, fmt.Errorf("invalid source of element ID")
	}
}
