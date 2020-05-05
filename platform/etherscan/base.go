package etherscan

import (
	"github.com/gin-gonic/gin"
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
)

type Platform struct {
	CoinIndex uint
	RpcURL    string
	client    Client
}

func Init(coin uint, api, rpc string) *Platform {
	return &Platform{
		CoinIndex: coin,
		RpcURL:    rpc,
		client:    Client{blockatlas.InitClient(api), 10},
	}
}

func (p *Platform) Coin() coin.Coin {
	return coin.Coins[p.CoinIndex]
}

func (p *Platform) RegisterRoutes(router gin.IRouter) {
	router.GET("/balance/:address", func(c *gin.Context) {
		p.getBalance(c)
	})
	router.GET("/address/:address", func(c *gin.Context) {
		p.getTransactions(c)
	})
	router.GET("/current_block", func(c *gin.Context) {
		p.getCurrentBlockNumber(c)
	})
	router.GET("/block/:block", func(c *gin.Context) {
		p.getBlockByNumber(c)
	})
	router.GET("/gas/price", func(c *gin.Context) {
		p.getGasPrice(c)
	})
	router.POST("/transaction/estimate", func(c *gin.Context) {
		p.getEstimatedGas(c)
	})
	router.POST("/transaction/send", func(c *gin.Context) {
		p.sendTransaction(c)
	})

}
