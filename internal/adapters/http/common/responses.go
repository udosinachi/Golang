package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Status  int         `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Success bool        `json:"success,omitempty"`
}

func SendBadRequest(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusBadRequest, APIResponse{
		Status:  http.StatusBadRequest,
		Message: message,
		Success: false,
	})
}

func SendUnauthorized(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusUnauthorized, APIResponse{
		Status:  http.StatusUnauthorized,
		Message: message,
		Success: false,
	})
}

func SendForbidden(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusForbidden, APIResponse{
		Status:  http.StatusForbidden,
		Message: message,
		Success: false,
	})
}

func SendPreconditionFailed(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusPreconditionFailed, APIResponse{
		Status:  http.StatusPreconditionFailed,
		Message: message,
		Success: false,
	})
}

func SendServerError(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusInternalServerError, APIResponse{
		Status:  http.StatusInternalServerError,
		Message: message,
		Success: false,
	})
}

func SendNotFound(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusNotFound, APIResponse{
		Status:  http.StatusNotFound,
		Message: message,
		Success: false,
	})
}

func SendCreated(ctx *gin.Context, payload interface{}, message string) {
	ctx.JSON(http.StatusCreated, APIResponse{
		Status:  http.StatusCreated,
		Data:    payload,
		Message: message,
		Success: true,
	})
}

func SendOk(ctx *gin.Context, payload interface{}, message string) {
	ctx.JSON(http.StatusOK, APIResponse{
		Status:  http.StatusOK,
		Data:    payload,
		Message: message,
		Success: true,
	})
}
