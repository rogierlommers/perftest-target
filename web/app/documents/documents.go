package documents

import (
	"math/rand"

	"github.com/gin-gonic/gin"
)

func GETDocuments(ctx *gin.Context) {
	document := generateDocumentID()
	ctx.JSON(200, gin.H{
		"documentID": document,
	})
}

func POSTDocuments(ctx *gin.Context) {
	ctx.JSON(201, gin.H{
		"result":     "document_created",
		"documentID": generateDocumentID(),
	})
}

func generateDocumentID() string {
	return "doc-" + generateRandomString(8)
}

func generateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, n)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
