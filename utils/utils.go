package utils

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type StandardResponse struct {
	Success      bool
	ErrorMessage string
	Data         interface{}
}

var (
	TimeLocation, _ = time.LoadLocation("Asia/Jakarta")
)

type CommonRequest struct {
	StartIndex int
	EndIndex   int
	Filter     map[string]interface{}
	SortBy     string
	SortType   string
}
type CustomContext struct {
	echo.Context
	CommonRequest
}

func ParseCommonRequest(c echo.Context) (ret CommonRequest, errRet error) {
	sortQuery := c.QueryParam("sort")
	if sortQuery != "" {
		sortStr := strings.Split(between(sortQuery, "[", "]"), ",")
		if len(sortStr) != 2 {
			errRet = errors.New("Sort params must have length of 2")
			return
		}
		ret.SortBy = strings.ReplaceAll(sortStr[0], `"`, ``)
		ret.SortType = strings.ReplaceAll(sortStr[1], `"`, ``)
	}
	rangeQuery := c.QueryParam("range")
	if rangeQuery != "" {
		rangesStr := strings.Split(between(rangeQuery, "[", "]"), ",")
		if len(rangesStr) != 2 {
			errRet = errors.New("Range params must have length of 2")
			return
		}
		startID, err := strconv.Atoi(rangesStr[0])
		if err != nil {
			errRet = errors.New("Range params must be a number")
			return
		}
		endID, err := strconv.Atoi(rangesStr[1])
		if err != nil {
			errRet = errors.New("Range params must be a number")
			return
		}
		ret.StartIndex = startID
		ret.EndIndex = endID
	}
	filterQuery := c.QueryParam("filter")
	if filterQuery != "" {
		err := json.Unmarshal([]byte(filterQuery), &ret.Filter)
		if err != nil {
			errRet = errors.New("Invalid filter params format")
			return
		}
		for _, v := range ret.Filter {
			switch result := v.(type) {
			case string:
				log.Println("string:", result)
			case []string:
				log.Println("[]string:", result)
			case int:
				log.Println("int:", result)
			case []int:
				log.Println("[]int:", result)
			default:
				errRet = errors.New("Invalid filter inside params format")
				return
			}
		}

	}
	return
}
func between(value string, a string, b string) string {
	// Get substring between two strings.
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return value[posFirstAdjusted:posLast]
}
