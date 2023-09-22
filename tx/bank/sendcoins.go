package bank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/pkg/errors"
	"me-offlinetx/config"
	"me-offlinetx/sdkclient"
	"strconv"
)

/*
SendCoin
示例一：普通转账交易
*/
func SendCoin(fromAddress string, toAddress string, amount int) ([]byte, error) {
	if sdkclient.MeClient.HTTPClient == nil {
		return nil, errors.New("client is nil")
	}
	fromAddr, err := sdk.AccAddressFromBech32(fromAddress)
	if err != nil {
		return nil, errors.Wrap(err, "err from address")
	}
	toAddr, err := sdk.AccAddressFromBech32(toAddress)
	if err != nil {
		return nil, errors.Wrap(err, "err to address")
	}

	coin, err := sdk.ParseCoinNormalized(strconv.Itoa(amount) + config.Denom)
	if err != nil {
		return nil, err
	}
	coins := sdk.NewCoins(coin)
	msg := banktypes.NewMsgSend(fromAddr, toAddr, coins)
	if msg.ValidateBasic() != nil {
		return nil, errors.Wrap(err, "err NewMsgSend")
	}

	privKey, err := sdkclient.MeWallet.PrivKey(fromAddress)
	if err != nil {
		return nil, errors.Wrap(err, "err address private key")
	}
	seq := sdkclient.MeWallet.IncrementNonce(fromAddress)
	num := sdkclient.MeWallet.AccountNum(fromAddress)

	tx, err := sdkclient.MeClient.BuildTx(msg, privKey, seq, num)
	if err != nil {
		return nil, errors.Wrap(err, "err BuildTx")
	}

	txBytes, err := sdkclient.MeClient.TxConfig.TxEncoder()(tx)
	if err != nil {
		return nil, errors.Wrap(err, "err TxEncoder")
	}

	return txBytes, nil
}

/*
SendToAdmin
示例二：从国库转账给超级管理员
*/
func SendToAdmin(fromAddress string, amount int) ([]byte, error) {
	if sdkclient.MeClient.HTTPClient == nil {
		return nil, errors.New("client is nil")
	}
	fromAddr, err := sdk.AccAddressFromBech32(fromAddress)
	if err != nil {
		return nil, errors.Wrap(err, "err from address")
	}

	coin, err := sdk.ParseCoinNormalized(strconv.Itoa(amount) + config.Denom)
	if err != nil {
		return nil, err
	}
	coins := sdk.NewCoins(coin)
	msg := banktypes.NewMsgSendToAdmin(fromAddr, coins)
	if msg.ValidateBasic() != nil {
		return nil, errors.Wrap(err, "err NewMsgSend")
	}

	privKey, err := sdkclient.MeWallet.PrivKey(fromAddress)
	if err != nil {
		return nil, errors.Wrap(err, "err address private key")
	}
	seq := sdkclient.MeWallet.IncrementNonce(fromAddress)
	num := sdkclient.MeWallet.AccountNum(fromAddress)

	tx, err := sdkclient.MeClient.BuildTx(msg, privKey, seq, num)
	if err != nil {
		return nil, errors.Wrap(err, "err BuildTx")
	}

	txBytes, err := sdkclient.MeClient.TxConfig.TxEncoder()(tx)
	if err != nil {
		return nil, errors.Wrap(err, "err TxEncoder")
	}

	return txBytes, nil
}
