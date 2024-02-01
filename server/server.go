package server

import (
	"DistriAI-Node/middleware"
	"DistriAI-Node/server/template"
	dbutils "DistriAI-Node/utils/db_utils"
	logs "DistriAI-Node/utils/log_utils"
	"fmt"
	"net/http"

	"github.com/gagliardetto/solana-go"
	"github.com/gin-gonic/gin"
)

func StartServer(serverPost string) error {
	logs.Normal("Start server")

	r := gin.Default()
	r.Use(middleware.Cors())
	workspace := r.Group(template.WORKSPACE)
	workspace.POST(template.TOKEN, getToken)

	// order := r.Group(template.ORDER)
	// order.POST(template.RT, getToken)

	err := r.Run("127.0.0.1:" + serverPost)
	if err != nil {
		logs.Error(fmt.Sprintf("gin run error: %v", err))
		return err
	}
	return nil
}

func getToken(c *gin.Context) {
	logs.Normal("GetToken API")

	var requestToken template.RequestToken
	if err := c.ShouldBindJSON(&requestToken); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db, err := dbutils.NewDB()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	buyerPublicKey, err := db.Get([]byte("buyer"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := db.Get([]byte("token"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Close()

	publicKeyStr := string(buyerPublicKey)
	publicKey, err := solana.PublicKeyFromBase58(publicKeyStr)
	if err != nil {
		logs.Error(fmt.Sprintf("PublicKeyFromBase58 error: %v", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	signature, err := solana.SignatureFromBase58(requestToken.Signature)

	if err != nil {
		logs.Error(fmt.Sprintf("SignatureFromBase58 error: %v", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if publicKey.Verify([]byte("workspace/token/"+publicKeyStr), signature) {
		c.JSON(http.StatusOK, gin.H{"token": string(token)})
	} else {
		logs.Error("Verify failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verify failed"})
		return
	}
}
