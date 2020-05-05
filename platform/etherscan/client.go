package etherscan

import (
	"errors"
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/blockatlas/pkg/numbers"
	"net/url"
	"strconv"
	"time"
)

type Client struct {
	blockatlas.Request
	MaxRetries int
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

func (c *Client) EstimateGas(tx string, to string, value int64) (int64, error) {
	var gasInfo StringResultPage
	values := url.Values{"module": {"proxy"},
		"action": {"eth_estimateGas"}, "to": {to},
		"value": {"0x" + strconv.FormatInt(value, 16)},
		"data":  {tx}}
	err := c.Get(&gasInfo, "api", values)
	if err != nil {
		return 0, err
	}
	retryCounter := 0

	for gasInfo.Message == "NOTOK" &&
		gasInfo.Result == "Max rate limit reached, please use API Key for higher rate limit" &&
		retryCounter < c.MaxRetries {
		logger.Info("Sleeping 1 second before retry")
		time.Sleep(3 * time.Second)
		err := c.Get(&gasInfo, "api", values)
		if err != nil {
			return 0, err
		}
		retryCounter += 1
	}

	if block_number, err := hexToInt(gasInfo.Result); err != nil {
		return 0, err
	} else {
		return block_number, nil
	}

}

func (c *Client) SendTransaction(tx string) (string, error) {
	var txInfo StringResultPage
	values := url.Values{"module": {"proxy"},
		"action": {"eth_sendRawTransaction"},
		"hex":    {"0x" + tx}}
	err := c.Get(&txInfo, "api", values)
	if err != nil {
		return "", err
	}
	retryCounter := 0
	for txInfo.Message == "NOTOK" &&
		txInfo.Result == "Max rate limit reached, please use API Key for higher rate limit" {
		if retryCounter < c.MaxRetries {
			logger.Info("Sleeping 3 second before retry")
			time.Sleep(3 * time.Second)
			err := c.Get(&txInfo, "api", values)
			if err != nil {
				return "", err
			}
			retryCounter += 1
		} else {
			return "", errors.New("Error sending transaction")
		}
	}

	return txInfo.Result, nil

}

func (c *Client) GasPrice() (int64, error) {
	var gasInfo StringResultPage
	values := url.Values{"module": {"proxy"}, "action": {"eth_gasPrice"}}
	err := c.Get(&gasInfo, "api", values)
	if err != nil {
		return 0, err
	}

	retryCounter := 0
	for gasInfo.Message == "NOTOK" &&
		gasInfo.Result == "Max rate limit reached, please use API Key for higher rate limit" &&
		retryCounter < c.MaxRetries {
		logger.Info("Sleeping 1 second before retry")
		time.Sleep(1 * time.Second)
		err := c.Get(&gasInfo, "api", values)
		if err != nil {
			return 0, err
		}
	}

	if block_number, err := hexToInt(gasInfo.Result); err != nil {
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

func (c *Client) GetBalance(address string) (balance string, err error) {
	var balanceInfo StringResultPage
	values := url.Values{"module": {"account"}, "action": {"balance"}, "address": {address}, "tag": {"latest"}}
	err = c.Get(&balanceInfo, "api", values)
	if err != nil {
		return
	}
	retryCounter := 0
	for balanceInfo.Message == "NOTOK" &&
		balanceInfo.Result == "Max rate limit reached, please use API Key for higher rate limit" &&
		retryCounter < c.MaxRetries {
		logger.Info("Sleeping 1 second before retry")
		time.Sleep(1 * time.Second)
		err := c.Get(&balanceInfo, "api", values)
		if err != nil {
			return "0", err
		}
	}

	return balanceInfo.Result, nil
}
