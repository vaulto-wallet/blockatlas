package etherscan

import (
	"github.com/gin-gonic/gin"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"net/http"
)

func (p *Platform) getBalance(c *gin.Context) {
	token := c.Query("token")
	address := c.Param("address")
	var err error
	var balance string

	if token != "" {
		//srcPage, err = p.client.GetTokenBalance(address, token)
		balance, err = p.client.GetBalance(address)
	} else {
		balance, err = p.client.GetBalance(address)
	}

	if apiError(c, err) {
		return
	}
	logger.Info("Balance", balance)

	c.JSON(http.StatusOK, &balance)
}
