package shared

import (
	"net/http"
	"strconv"
	"strings"
)

const DefaultPage = 1
const DefaultSort = "desc"
const DefaultOffset = 20

func GetPage(req *http.Request) int {
	if pageParam := req.URL.Query().Get("page"); pageParam != "" {
		page, err := strconv.ParseInt(pageParam, 10, 64)
		if err == nil {
			return int(page)
		}
	}

	return DefaultPage
}

func GetSort(req *http.Request) string {
	if sortParam := strings.ToLower(req.URL.Query().Get("sort")); sortParam != "" {
		valid := map[string]bool{"desc": true, "asc": true}
		if valid[sortParam] {
			return sortParam
		}
	}

	return DefaultSort
}

func GetOffset(req *http.Request) int {
	if offsetParam := req.URL.Query().Get("offset"); offsetParam != "" {
		offset, err := strconv.ParseInt(offsetParam, 10, 64)
		if err == nil {
			return int(offset)
		}
	}

	return DefaultOffset
}
