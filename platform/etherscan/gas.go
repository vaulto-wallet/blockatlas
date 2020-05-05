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

	c.JSON(http.StatusOK, InterfaceResultPage{strconv.FormatInt(gasPrice, 10)})
}
