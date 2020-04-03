package etherscan

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (p *Platform) getGasPrice(c *gin.Context) {
	gasPrice, err := p.client.GasPrice()
	if apiError(c, err) {
		return
	}

	c.JSON(http.StatusOK, GasPage{strconv.FormatInt(gasPrice, 10)})
}

func (p *Platform) getEstimatedGas(c *gin.Context) {
	gasPrice, err := p.client.GasPrice()
	if apiError(c, err) {
		return
	}

	c.JSON(http.StatusOK, GasPage{strconv.FormatInt(gasPrice, 10)})
}
