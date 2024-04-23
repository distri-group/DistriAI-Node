package server

import (
	"DistriAI-Node/config"
	"DistriAI-Node/middleware"
	"DistriAI-Node/server/template"
	dbutils "DistriAI-Node/utils/db_utils"
	logs "DistriAI-Node/utils/log_utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gin-gonic/gin"
)

func StartServer(serverPort string) error {
	logs.Normal("Start server")

	r := gin.Default()
	r.Use(middleware.Cors())
	r.SetTrustedProxies([]string{"127.0.0.1"})
	workspace := r.Group(template.WORKSPACE)
	workspace.GET("/debugToken/:signature", getDebugToken)

	err := r.Run("127.0.0.1:" + serverPort)
	if err != nil {
		logs.Error(fmt.Sprintf("gin run error: %v", err))
		return err
	}
	return nil
}

func getDebugToken(c *gin.Context) {
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
		db.Close()
		return
	}
	token, err := db.Get([]byte("token"))
	db.Close()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	currentTime := time.Now().Unix() / 100
	msg := fmt.Sprintf("workspace/token/%v/%v", currentTime, publicKeyStr)

	if publicKey.Verify([]byte(msg), out) {
		RedirectWorkspace(c, string(token))
	} else {
		currentTime -= 1
		msg = fmt.Sprintf("workspace/token/%v/%v", currentTime, publicKeyStr)

		if publicKey.Verify([]byte(msg), out) {
			RedirectWorkspace(c, string(token))
		} else {
			if publicKey.Verify([]byte("deploy/token/"+publicKeyStr), out) {
				deployURL := fmt.Sprintf("http://%v:%v",
					config.GlobalConfig.Console.OuterNetIP,
					config.GlobalConfig.Console.OuterNetPort)

				logs.Normal(fmt.Sprintf("Redirect to: %v", deployURL))

				c.Redirect(http.StatusFound, deployURL)
			} else {
				logs.Error("Verify failed")
				c.JSON(http.StatusBadRequest, gin.H{"error": "Verify failed"})
			}
		}
	}
}

func RedirectWorkspace(c *gin.Context, token string) {
	workspaceURL := fmt.Sprintf("http://%v:%v?token=%v",
		config.GlobalConfig.Console.OuterNetIP,
		config.GlobalConfig.Console.OuterNetPort,
		token)

	logs.Normal(fmt.Sprintf("Redirect to: %v", workspaceURL))

	c.Redirect(http.StatusFound, workspaceURL)
}
