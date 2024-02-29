package server

import (
	"DistriAI-Node/config"
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
	workspace.GET("/debugToken/:signature", getDebugToken)

	err := r.Run("127.0.0.1:" + serverPost)
	if err != nil {
		logs.Error(fmt.Sprintf("gin run error: %v", err))
		return err
	}
	return nil
}

func getDebugToken(c *gin.Context) {
	logs.Normal("GetDebugToken API")

	signature := c.Param("signature")
	logs.Normal(fmt.Sprintf("signature: %v", signature))

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

	out, err := solana.SignatureFromBase58(signature)
	if err != nil {
		logs.Error(fmt.Sprintf("SignatureFromBase58 error: %v", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if publicKey.Verify([]byte("workspace/token/"+publicKeyStr), out) {
		workspaceURL := fmt.Sprintf("http://%v:%v?token=%v",
			config.GlobalConfig.Console.OuterNetIP,
			config.GlobalConfig.Console.OuterNetPort,
			string(token))

		c.Redirect(http.StatusFound, workspaceURL)
	} else {
		logs.Error("Verify failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verify failed"})
		return
	}
}
