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
 *
 */

package txVerification

import (
	"errors"
	"github.com/bnb-chain/zkbas-crypto/ffmath"
	"github.com/bnb-chain/zkbas-crypto/wasm/legend/legendTxTypes"
	"github.com/bnb-chain/zkbas/common/commonAsset"
	"github.com/bnb-chain/zkbas/common/commonConstant"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
	"math/big"
)

/*
	VerifyTransferNftTx:
	accounts order is:
	- FromAccount
		- Assets
			- AssetGas
	- GasAccount
		- Assets
			- AssetGas
*/
func VerifyTransferNftTxInfo(
	accountInfoMap map[int64]*AccountInfo,
	nftInfo *NftInfo,
	txInfo *TransferNftTxInfo,
) (txDetails []*MempoolTxDetail, err error) {
	// verify params
	if accountInfoMap[txInfo.FromAccountIndex] == nil ||
		accountInfoMap[txInfo.GasAccountIndex] == nil ||
		accountInfoMap[txInfo.FromAccountIndex].AssetInfo == nil ||
		accountInfoMap[txInfo.FromAccountIndex].AssetInfo[txInfo.GasFeeAssetId] == nil ||
		nftInfo == nil ||
		nftInfo.OwnerAccountIndex != txInfo.FromAccountIndex ||
		nftInfo.NftIndex != txInfo.NftIndex ||
		txInfo.GasFeeAssetAmount.Cmp(ZeroBigInt) < 0 {
		logx.Errorf("[VerifyTransferNftTxInfo] invalid params")
		return nil, errors.New("[VerifyTransferNftTxInfo] invalid params")
	}
	// verify nonce
	if txInfo.Nonce != accountInfoMap[txInfo.FromAccountIndex].Nonce {
		log.Println("[VerifyTransferNftTxInfo] invalid nonce")
		return nil, errors.New("[VerifyTransferNftTxInfo] invalid nonce")
	}
	// set tx info
	var (
		assetDeltaMap = make(map[int64]map[int64]*big.Int)
		newNftInfo    *NftInfo
	)
	// init delta map
	assetDeltaMap[txInfo.FromAccountIndex] = make(map[int64]*big.Int)
	if assetDeltaMap[txInfo.GasAccountIndex] == nil {
		assetDeltaMap[txInfo.GasAccountIndex] = make(map[int64]*big.Int)
	}
	// from account asset Gas
	assetDeltaMap[txInfo.FromAccountIndex][txInfo.GasFeeAssetId] = ffmath.Neg(txInfo.GasFeeAssetAmount)
	// to account nft info
	newNftInfo = &NftInfo{
		NftIndex:            nftInfo.NftIndex,
		CreatorAccountIndex: nftInfo.CreatorAccountIndex,
		OwnerAccountIndex:   txInfo.ToAccountIndex,
		NftContentHash:      nftInfo.NftContentHash,
		NftL1TokenId:        nftInfo.NftL1TokenId,
		NftL1Address:        nftInfo.NftL1Address,
		CreatorTreasuryRate: nftInfo.CreatorTreasuryRate,
		CollectionId:        nftInfo.CollectionId,
	}
	// gas account asset Gas
	if assetDeltaMap[txInfo.GasAccountIndex][txInfo.GasFeeAssetId] == nil {
		assetDeltaMap[txInfo.GasAccountIndex][txInfo.GasFeeAssetId] = txInfo.GasFeeAssetAmount
	} else {
		assetDeltaMap[txInfo.GasAccountIndex][txInfo.GasFeeAssetId] = ffmath.Add(
			assetDeltaMap[txInfo.GasAccountIndex][txInfo.GasFeeAssetId],
			txInfo.GasFeeAssetAmount,
		)
	}
	// check balance
	if accountInfoMap[txInfo.FromAccountIndex].AssetInfo[txInfo.GasFeeAssetId].Balance.Cmp(txInfo.GasFeeAssetAmount) < 0 {
		logx.Errorf("[VerifyMintNftTxInfo] you don't have enough balance of asset Gas")
		return nil, errors.New("[VerifyMintNftTxInfo] you don't have enough balance of asset Gas")
	}
	// compute hash
	hFunc := mimc.NewMiMC()
	msgHash, err := legendTxTypes.ComputeTransferNftMsgHash(txInfo, hFunc)
	if err != nil {
		logx.Errorf("[VerifyMintNftTxInfo] unable to compute hash: %s", err.Error())
		return nil, err
	}
	// verify signature
	hFunc.Reset()
	pk, err := ParsePkStr(accountInfoMap[txInfo.FromAccountIndex].PublicKey)
	if err != nil {
		return nil, err
	}
	isValid, err := pk.Verify(txInfo.Sig, msgHash, hFunc)
	if err != nil {
		log.Println("[VerifyTransferNftTxInfo] unable to verify signature:", err)
		return nil, err
	}
	if !isValid {
		log.Println("[VerifyTransferNftTxInfo] invalid signature")
		return nil, errors.New("[VerifyTransferNftTxInfo] invalid signature")
	}
	// compute tx details
	// from account asset gas
	order := int64(0)
	accountOrder := int64(0)
	txDetails = append(txDetails, &MempoolTxDetail{
		AssetId:      txInfo.GasFeeAssetId,
		AssetType:    GeneralAssetType,
		AccountIndex: txInfo.FromAccountIndex,
		AccountName:  accountInfoMap[txInfo.FromAccountIndex].AccountName,
		BalanceDelta: commonAsset.ConstructAccountAsset(
			txInfo.GasFeeAssetId, ffmath.Neg(txInfo.GasFeeAssetAmount), ZeroBigInt, ZeroBigInt).String(),
		Order:        order,
		AccountOrder: accountOrder,
	})
	// to account empty delta
	order++
	accountOrder++
	txDetails = append(txDetails, &MempoolTxDetail{
		AssetId:      txInfo.GasFeeAssetId,
		AssetType:    GeneralAssetType,
		AccountIndex: txInfo.ToAccountIndex,
		AccountName:  accountInfoMap[txInfo.ToAccountIndex].AccountName,
		BalanceDelta: commonAsset.ConstructAccountAsset(
			txInfo.GasFeeAssetId, ZeroBigInt, ZeroBigInt, ZeroBigInt).String(),
		Order:        order,
		AccountOrder: accountOrder,
	})
	// to account nft delta
	order++
	txDetails = append(txDetails, &MempoolTxDetail{
		AssetId:      txInfo.NftIndex,
		AssetType:    NftAssetType,
		AccountIndex: txInfo.ToAccountIndex,
		AccountName:  accountInfoMap[txInfo.ToAccountIndex].AccountName,
		BalanceDelta: newNftInfo.String(),
		Order:        order,
		AccountOrder: commonConstant.NilAccountOrder,
	})
	// gas account asset gas
	order++
	accountOrder++
	txDetails = append(txDetails, &MempoolTxDetail{
		AssetId:      txInfo.GasFeeAssetId,
		AssetType:    GeneralAssetType,
		AccountIndex: txInfo.GasAccountIndex,
		AccountName:  accountInfoMap[txInfo.GasAccountIndex].AccountName,
		BalanceDelta: commonAsset.ConstructAccountAsset(
			txInfo.GasFeeAssetId, txInfo.GasFeeAssetAmount, ZeroBigInt, ZeroBigInt).String(),
		Order:        order,
		AccountOrder: accountOrder,
	})
	return txDetails, nil
}
