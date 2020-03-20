package etherscan

import (
	"github.com/gin-gonic/gin"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"net/http"
)

func (p *Platform) getBalance(c *gin.Context) {
	token := c.Query("token")
	address := c.Param("address")
	var srcPage BalancePage
	var err error

	if token != "" {
		//srcPage, err = p.client.GetTokenBalance(address, token)
		srcPage, err = p.client.GetBalance(address)
	} else {
		srcPage, err = p.client.GetBalance(address)
	}

	if apiError(c, err) {
		return
	}
	logger.Info("Balance", srcPage)

	c.JSON(http.StatusOK, &srcPage)
}
