package etherscan

import (
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"github.com/trustwallet/blockatlas/pkg/numbers"
	"net/url"
	"strconv"
)

type Client struct {
	blockatlas.Request
}

func (c *Client) GetTxs(address string) (*Page, error) {
	return c.getTxs(url.Values{"address": {address}, "module": {"account"}, "action": {"txlist"}})
}

func (c *Client) GetTxsWithContract(address, contract string) (*Page, error) {
	return c.getTxs(url.Values{"address": {address}, "contract": {contract}})
}

func (c *Client) getTxs(query url.Values) (page *Page, err error) {
	err = c.Get(&page, "api", query)
	return
}

func (c *Client) GetBlockByNumber(num int64) (page *BlockPage, err error) {
	values := url.Values{"module": {"proxy"}, "action": {"eth_getBlockByNumber"}, "tag": {"0x" + strconv.FormatInt(num, 16)}, "boolean": {"true"}}
	err = c.Get(&page, "api", values)
	return
}

func (c *Client) CurrentBlockNumber() (int64, error) {
	var nodeInfo NodeInfo
	values := url.Values{"module": {"proxy"}, "action": {"eth_blockNumber"}}
	err := c.Get(&nodeInfo, "api", values)
	if err != nil {
		return 0, err
	}
	if block_number, err := hexToInt(nodeInfo.LatestBlock); err != nil {
		return 0, err
	} else {
		return block_number, nil
	}

}

func (c *Client) GasPrice() (int64, error) {
	var gasInfo GasPage
	values := url.Values{"module": {"proxy"}, "action": {"eth_gasPrice"}}
	err := c.Get(&gasInfo, "api", values)
	if err != nil {
		return 0, err
	}
	if block_number, err := hexToInt(gasInfo.Gas); err != nil {
		return 0, err
	} else {
		return block_number, nil
	}

}

func (c *Client) GetTokens(address string) (tp *TokenPage, err error) {
	query := url.Values{
		"address": {address},
	}
	err = c.Get(&tp, "tokens", query)
	return
}

func hexToInt(hex string) (int64, error) {
	nonceStr, err := numbers.HexToDecimal(hex)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(nonceStr, 10, 64)
}

func (c *Client) GetBalance(address string) (page BalancePage, err error) {
	values := url.Values{"module": {"account"}, "action": {"balance"}, "address": {address}, "tag": {"latest"}}
	//var page BalancePage
	err = c.Get(&page, "api", values)
	return
}
