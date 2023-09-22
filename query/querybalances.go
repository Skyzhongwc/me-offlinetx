package query

import (
	"context"
	"fmt"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"me-offlinetx/config"
	"me-offlinetx/sdkclient"
)

func QueryBalances(addr string) (*banktypes.QueryBalanceResponse, error) {
	bankRequest := banktypes.QueryBalanceRequest{
		Address: addr,
		Denom:   config.UDenom,
	}

	client2 := banktypes.NewQueryClient(sdkclient.MeClient.Conn)
	response, err := client2.Balance(context.Background(), &bankRequest)
	if err != nil {
		fmt.Println("query balances err：", err)
		return nil, err
	}

	return response, nil
}
