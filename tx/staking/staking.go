package staking

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	cosmoscli "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	staketypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/pkg/errors"
	"me-offlinetx/config"
	"me-offlinetx/sdkclient"
	"strconv"
	"strings"
)

/*
AddStake
示例一：增加验证者质押
*/
func AddStake(validatorAddress string, amount int, superAdmin string) ([]byte, error) {
	if sdkclient.MeClient.HTTPClient == nil {
		return nil, errors.New("client is nil")
	}

	stakerAddr, err := sdk.AccAddressFromBech32(superAdmin)
	if err != nil {
		return nil, errors.Wrap(err, "err staker address")
	}

	valAddr, err := sdk.ValAddressFromBech32(validatorAddress)
	if err != nil {
		return nil, errors.Wrap(err, "err validator address")
	}

	coin, err := sdk.ParseCoinNormalized(strconv.Itoa(amount) + config.Denom)
	if err != nil {
		return nil, err
	}

	msg := staketypes.NewMsgStake(stakerAddr, valAddr, coin)
	if msg.ValidateBasic() != nil {
		return nil, errors.Wrap(err, "err NewMsgStake")
	}

	privKey, err := sdkclient.MeWallet.PrivKey(superAdmin)
	if err != nil {
		return nil, errors.Wrap(err, "err superAdmin private key")
	}
	seq := sdkclient.MeWallet.IncrementNonce(superAdmin)
	num := sdkclient.MeWallet.AccountNum(superAdmin)

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
CreateValidator
示例二：创建验证者
*/
func CreateValidator(msg, signer, privateKey string, fee, feeLimit int64) string {
	txBuilder, signerData := GetTxBuilderAndSigner(msg, signer, fee, feeLimit)
	signature := Sign(privateKey, txBuilder, signerData, signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON)
	sigV2, _ := sdkclient.MeClient.TxConfig.UnmarshalSignatureJSON([]byte(signature))
	txBuilder.SetSignatures(sigV2[0])
	bytes, _ := sdkclient.MeClient.TxConfig.TxJSONEncoder()(txBuilder.GetTx())
	return string(bytes)
}

func Sign(priKeyStr string, txBuilder cosmoscli.TxBuilder, signerData authsigning.SignerData, mode signing.SignMode) string {
	// 创建3个私钥
	priBytes1, _ := hex.DecodeString(priKeyStr)
	privKey := &secp256k1.PrivKey{Key: priBytes1}
	// First round: we gather all the signer infos. We use the "set empty signature" hack to do that.
	sig := signing.SignatureV2{
		PubKey: privKey.PubKey(),
		Data: &signing.SingleSignatureData{
			SignMode:  mode,
			Signature: nil,
		},
		Sequence: signerData.Sequence,
	}
	sig, _ = tx.SignWithPrivKey(mode, signerData, txBuilder, privKey, sdkclient.MeClient.TxConfig, signerData.Sequence)
	bytes, _ := sdkclient.MeClient.TxConfig.MarshalSignatureJSON([]signing.SignatureV2{sig})
	return string(bytes)
}

func GetTxBuilderAndSigner(msg, signerData string, fee, gasLimit int64) (cosmoscli.TxBuilder, authsigning.SignerData) {
	txMsg := ParseMsg(msg)
	signer := authsigning.SignerData{}
	json.Unmarshal([]byte(signerData), &signer)
	txBuilder := sdkclient.MeClient.TxConfig.NewTxBuilder()
	txBuilder.SetMsgs(txMsg...)
	if fee > 0 {
		txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewInt64Coin(config.UDenom, fee)))
	}
	txBuilder.SetGasLimit(uint64(gasLimit))
	return txBuilder, signer
}

func ParseMsg(msg string) []sdk.Msg {
	msgBytes := []byte(msg)
	realMsg := &staketypes.MsgCreateValidator{}
	err := json.Unmarshal(msgBytes, realMsg)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	pk := &CreateValidatorPk{}
	json.Unmarshal(msgBytes, pk)
	validatorAddress := realMsg.ValidatorAddress
	stakerAddress := realMsg.StakerAddress
	if len(strings.TrimSpace(validatorAddress)) == 0 {
		validatorAddress = GetValAddr(stakerAddress)
	}
	if len(strings.TrimSpace(stakerAddress)) == 0 {
		stakerAddress = GetAccAddr(validatorAddress)
	}
	createmsg := &staketypes.MsgCreateValidator{
		Description:      realMsg.Description,
		Commission:       realMsg.Commission,
		StakerAddress:    stakerAddress,
		ValidatorAddress: validatorAddress,
		Pubkey:           pk.GetPubkeyAny(),
		MinSelfStake:     realMsg.MinSelfStake,
		Value:            realMsg.Value,
	}
	return []sdk.Msg{createmsg}
}

func (pk *CreateValidatorPk) GetPubkeyAny() *codectypes.Any {
	bytes, _ := base64.StdEncoding.DecodeString(pk.Pubkey)
	pkAny, _ := codectypes.NewAnyWithValue(&ed25519.PubKey{Key: bytes})
	return pkAny
}

func GetValAddr(accAddr string) string {
	acc, _ := sdk.AccAddressFromBech32(accAddr)
	address := sdk.ValAddress(acc)
	return address.String()
}
func GetAccAddr(valAddr string) string {
	val, _ := sdk.ValAddressFromBech32(valAddr)
	address := sdk.AccAddress(val)
	return address.String()
}

func Kyc(kycAddress, regionId string, account string, accSeq, accNum uint64) ([]byte, error) {
	priBytes, _ := hex.DecodeString(config.SuperAdminPriKey)
	priv := &secp256k1.PrivKey{Key: priBytes}
	txBuilder := sdkclient.MeClient.TxConfig.NewTxBuilder()
	msg1 := &staketypes.MsgNewKyc{
		Creator:  kycAddress,
		Account:  "me1wnz7x68mcjh4y6zwd6mnq9zsze6zvcqmrehxss",
		RegionId: regionId,
	}

	msg2 := &staketypes.MsgNewKyc{
		Creator:  kycAddress,
		Account:  "me1t0vjavhe9q3dzmgscq27fs9s65020gvkmnf9rr",
		RegionId: regionId,
	}
	msgs := []sdk.Msg{msg1, msg2}
	err := txBuilder.SetMsgs(msgs...)
	if err != nil {
		return nil, err
	}
	fees := sdk.NewCoins(sdk.NewInt64Coin(config.UDenom, 100))
	txBuilder.SetGasLimit(uint64(400000))
	txBuilder.SetFeeAmount(fees)
	if err = txBuilder.SetSignatures(signing.SignatureV2{
		PubKey: priv.PubKey(),
		Data: &signing.SingleSignatureData{
			SignMode:  sdkclient.MeClient.TxConfig.SignModeHandler().DefaultMode(),
			Signature: nil,
		},
		Sequence: accSeq,
	}); err != nil {
		return nil, err
	}
	// Second round: all signer infos are set, so each signer can sign.
	signerData := authsigning.SignerData{
		ChainID:       config.ChainID,
		AccountNumber: accNum,
		Sequence:      accSeq,
	}
	sigV2, err := tx.SignWithPrivKey(sdkclient.MeClient.TxConfig.SignModeHandler().DefaultMode(), signerData, txBuilder, priv, sdkclient.MeClient.TxConfig, accSeq)
	if err != nil {
		return nil, err
	}
	if err = txBuilder.SetSignatures(sigV2); err != nil {
		return nil, err
	}
	tx := txBuilder.GetTx()
	bytes, err := sdkclient.MeClient.TxConfig.TxEncoder()(tx)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

type CreateValidatorPk struct {
	Pubkey string `json:"pubkey"`
}
