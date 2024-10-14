package util

import (
	"fmt"
	"geoserver/api/internal/types"
	"strings"
)

func ParseErrorCode(err error) *types.ErrorResponse {
	if err == nil {
		return &types.ErrorResponse{Code: 200, Message: "Success"}
	}
	fmt.Println(err.Error())
	switch {
	case strings.Contains(err.Error(), "500"):
		return &types.ErrorResponse{Code: 500, Message: "Internal Server Error: " + "Data source already exists"}
	case strings.Contains(err.Error(), "404"):
		return &types.ErrorResponse{Code: 404, Message: "Internal Server Error: " + "Resource not found"}
	default:
		return &types.ErrorResponse{Code: 400, Message: "Bad Request"}
	}
}
