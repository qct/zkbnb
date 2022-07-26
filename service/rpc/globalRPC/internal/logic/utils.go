/*
 * Copyright © 2021 Zkbas Protocol
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package logic

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"

	"github.com/bnb-chain/zkbas/common/commonConstant"
	"github.com/bnb-chain/zkbas/common/commonTx"
	"github.com/bnb-chain/zkbas/common/model/mempool"
	"github.com/bnb-chain/zkbas/common/model/sysconfig"
	"github.com/bnb-chain/zkbas/common/sysconfigName"
	"github.com/bnb-chain/zkbas/common/util"
	"github.com/bnb-chain/zkbas/service/rpc/globalRPC/internal/logic/errcode"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

func GetTxTypeArray(txType uint) ([]uint8, error) {
	switch txType {
	case L2TransferType:
		return []uint8{commonTx.TxTypeTransfer}, nil
	case LiquidityType:
		return []uint8{commonTx.TxTypeAddLiquidity, commonTx.TxTypeRemoveLiquidity}, nil
	case L2SwapType:
		return []uint8{commonTx.TxTypeSwap}, nil
	case WithdrawAssetsType:
		return []uint8{commonTx.TxTypeWithdraw}, nil
	default:
		errInfo := fmt.Sprintf("[GetTxTypeArray] txType error: %v", txType)
		logx.Error(errInfo)
		return []uint8{}, errors.New(errInfo)
	}
}

func ComputeL2TxTxHash(txInfo string) string {
	hFunc := mimc.NewMiMC()
	hFunc.Write([]byte(txInfo))
	return base64.StdEncoding.EncodeToString(hFunc.Sum(nil))
}

func ConstructMempoolTx(
	txType int64,
	gasFeeAssetId int64,
	gasFeeAssetAmount string,
	nftIndex int64,
	pairIndex int64,
	assetId int64,
	txAmount string,
	toAddress string,
	txInfo string,
	memo string,
	accountIndex int64,
	nonce int64,
	expiredAt int64,
	txDetails []*mempool.MempoolTxDetail,
) (txId string, mempoolTx *mempool.MempoolTx) {
	txHash := ComputeL2TxTxHash(txInfo)
	return txHash, &mempool.MempoolTx{
		TxHash:         txHash,
		TxType:         txType,
		GasFeeAssetId:  gasFeeAssetId,
		GasFee:         gasFeeAssetAmount,
		NftIndex:       nftIndex,
		PairIndex:      pairIndex,
		AssetId:        assetId,
		TxAmount:       txAmount,
		NativeAddress:  toAddress,
		MempoolDetails: txDetails,
		TxInfo:         txInfo,
		ExtraInfo:      "",
		Memo:           memo,
		AccountIndex:   accountIndex,
		Nonce:          nonce,
		ExpiredAt:      expiredAt,
		L2BlockHeight:  commonConstant.NilBlockHeight,
		Status:         mempool.PendingTxStatus,
	}
}

func CreateMempoolTx(
	nMempoolTx *mempool.MempoolTx,
	redisConnection *redis.Redis,
	mempoolModel mempool.MempoolModel,
) (err error) {
	var keys []string
	for _, mempoolTxDetail := range nMempoolTx.MempoolDetails {
		keys = append(keys, util.GetAccountKey(mempoolTxDetail.AccountIndex))
	}
	_, err = redisConnection.Del(keys...)
	if err != nil {
		logx.Errorf("[CreateMempoolTx] error with redis: %s", err.Error())
		return err
	}
	// write into mempool
	err = mempoolModel.CreateBatchedMempoolTxs([]*mempool.MempoolTx{nMempoolTx})
	if err != nil {
		errInfo := fmt.Sprintf("[CreateMempoolTx] %s", err.Error())
		logx.Error(errInfo)
		return errors.New(errInfo)
	}
	return nil
}

func CheckGasAccountIndex(txGasAccountIndex int64, sysConfigModel sysconfig.SysconfigModel) error {
	gasAccountIndexConfig, err := sysConfigModel.GetSysconfigByName(sysconfigName.GasAccountIndex)
	if err != nil {
		logx.Errorf("[GetSysconfigByName] err: %v", err)
		return err
	}
	gasAccountIndex, err := strconv.ParseInt(gasAccountIndexConfig.Value, 10, 64)
	if err != nil {
		logx.Errorf("[ParseInt] param:%v,err:%v", gasAccountIndexConfig.Value, err)
		return err
	}
	if gasAccountIndex != txGasAccountIndex {
		logx.Errorf("[ParseInt] param:%v, txGasAccountIndex:%v, err:%v", gasAccountIndex, txGasAccountIndex, err)
		return errcode.ErrInvalidGasAccountIndex
	}
	return nil
}
