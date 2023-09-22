package main

import (
	"fmt"
	"me-offlinetx/config"
	"me-offlinetx/sdkclient"
	"me-offlinetx/tx/staking"
)

func main() {
	err := sdkclient.InitClient(config.DefaultRPCURI, config.DefaultGRPCURI)
	if err != nil {
		fmt.Println("InitClient err: ", err)
		return
	}

	err = sdkclient.ImportWallet(config.SuperAdminPriKey)
	if err != nil {
		fmt.Println("ImportWallet err: ", err)
		return
	}

	//---------------------------------------示例一：普通转账------------------------------------------------//

	//// 转账fromAddr
	//fromAddr, _ := sdkclient.MeWallet.Address(config.SuperAdminPriKey)
	//
	//// 转账toAddr
	//toAddr := config.ToAddress
	//
	//// 查询余额
	//response, err := query.QueryBalances(toAddr)
	//if err != nil {
	//	fmt.Println("查询余额失败：", err)
	//}
	//fmt.Println("转账前toAddr余额：", *response)
	//
	//// 转账金额 2mec (1mec=1000000umec)
	//amount := 2
	//// 构建转账交易
	//txBytes, err := bank.SendCoin(fromAddr, toAddr, amount)
	//if err != nil {
	//	fmt.Println("转账失败：", err)
	//	return
	//}
	//
	//// 发送交易
	//err = tx.SendTx(txBytes)
	//if err != nil {
	//	fmt.Println("发送交易失败", err)
	//	return
	//}
	//
	//// 等待10秒，交易上链
	//time.Sleep(10 * time.Second)
	//
	//// 查询余额
	//response2, err := query.QueryBalances(toAddr)
	//if err != nil {
	//	fmt.Println("查询余额失败：", err)
	//}
	//fmt.Println("转账后toAddr余额：", *response2)

	//---------------------------------------示例二：从国库给超管转账------------------------------------------//

	//// 超管地址
	//superAdmin, err := sdkclient.MeWallet.Address(config.SuperAdminPriKey)
	//if err != nil {
	//	fmt.Println("获取地址失败：", err)
	//}
	//
	//// 查询superAdmin余额
	//response, err := query.QueryBalances(superAdmin)
	//if err != nil {
	//	fmt.Println("查询余额失败：", err)
	//}
	//fmt.Println("转账前superAdmin余额：", *response)
	//
	//// 转账金额 2mec
	//amount := 2
	//// 构建从国库转账给超管交易
	//txBytes, err := bank.SendToAdmin(superAdmin, amount)
	//if err != nil {
	//	fmt.Println("从国库转账给超管失败：", err)
	//	return
	//}
	//
	//fmt.Println(base64.StdEncoding.EncodeToString(txBytes))

	//// 发送交易
	//err = tx.SendTx(txBytes)
	//if err != nil {
	//	fmt.Println("发送交易失败", err)
	//	return
	//}

	// 等待10秒，交易上链
	//time.Sleep(10 * time.Second)
	//
	//// 查询superAdmin余额
	//response2, err := query.QueryBalances(superAdmin)
	//if err != nil {
	//	fmt.Println("查询余额失败：", err)
	//}
	//fmt.Println("转账后superAdmin余额：", *response2)

	//---------------------------------------示例三：增加验证者质押--------------------------------------------//

	//// 超管地址
	//superAdmin, err := sdkclient.MeWallet.Address(config.SuperAdminPriKey)
	//if err != nil {
	//	fmt.Println("获取地址失败：", err)
	//}
	//// 质押金额
	//amount := 2
	//// 构建增加验证者质押交易
	//txBytes, err := staking.AddStake(config.ValidatorAddress, amount, superAdmin)
	//if err != nil {
	//	fmt.Println("增加质押失败：", err)
	//	return
	//}
	//
	//// 发送交易
	//err = tx.SendTx(txBytes)
	//if err != nil {
	//	fmt.Println("发送交易失败", err)
	//	return
	//}

	//-------------------------------------------------------------------------------------//

	//----------------------------------------创建验证者-------------------------------------//
	privateKey := config.SuperAdminPriKey
	msg := "{\"commission\":{\"max_change_rate\":\"0.01\",\"max_rate\":\"0.2\",\"rate\":\"0.1\"},\"description\":{\"moniker\":\"node1\"},\"min_self_stake\":\"1000000\",\"pubkey\":\"SWzLMwAGI0UqiYngQlRSmY4whI1ceywVDVDBKXVO4CE=\",\"validator_address\":\"mevaloper1p56mzra0p9anwg33uf96c5fwd59mltf85nnsly\",\"value\":{\"amount\":\"10000000000000\",\"denom\":\"umec\"}}"
	signer := "{\"AccountNumber\":1,\"Address\":\"me1csjcnf5xnskl2vstxyka04yuf437esjcv2a4vw\",\"ChainID\":\"me-chain\",\"Sequence\":5}"
	tx := staking.CreateValidator(msg, signer, privateKey, int64(100), int64(200000))
	fmt.Println(tx)

}
