package etherscan

import (
	"github.com/gin-gonic/gin"
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"net/http"
	"strconv"
)

func (p *Platform) getCurrentBlockNumber(c *gin.Context) {
	blockNumber, err := p.CurrentBlockNumber()

	if apiError(c, err) {
		return
	}

	c.JSON(http.StatusOK, BlockNumberPage{blockNumber})
}

func (p *Platform) getBlockByNumber(c *gin.Context) {
	blockNumber := c.Param("block")
	block, err := strconv.ParseInt(blockNumber, 10, 64)
	srcPage, err := p.GetBlockByNumber(block)

	if apiError(c, err) {
		return
	}

	c.JSON(http.StatusOK, srcPage)
}

func (p *Platform) CurrentBlockNumber() (int64, error) {
	return p.client.CurrentBlockNumber()
}

func (p *Platform) GetBlockByNumber(num int64) (*blockatlas.Block, error) {
	if srcPage, err := p.client.GetBlockByNumber(num); err == nil {
		var txs []blockatlas.Tx
		for _, srcTx := range srcPage.Block.Docs {
			txs = AppendTxs(txs, &srcTx, p.CoinIndex)
		}
		return &blockatlas.Block{
			Number: num,
			ID:     strconv.FormatInt(num, 10),
			Txs:    txs,
		}, nil
	} else {
		return nil, err
	}
}
