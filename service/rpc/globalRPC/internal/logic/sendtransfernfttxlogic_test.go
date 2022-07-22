/*
 * Copyright © 2021 Zecrey Protocol
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
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/zecrey-labs/zecrey-legend/common/commonAsset"
	"github.com/zecrey-labs/zecrey-legend/pkg/zerror"
	"github.com/zecrey-labs/zecrey-legend/service/rpc/globalRPC/globalRPCProto"
	"github.com/zecrey-labs/zecrey-legend/service/rpc/globalRPC/internal/repo/commglobalmap"
	"github.com/zecrey-labs/zecrey-legend/service/rpc/globalRPC/internal/repo/failtx"
)

func TestSendTransferNftTxLogic_SendTransferNftTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCommglobalmap := commglobalmap.NewMockCommglobalmap(ctrl)
	mockFailtx := failtx.NewMockModel(ctrl)
	l := &SendTransferNftTxLogic{
		commglobalmap: mockCommglobalmap,
		failtx:        mockFailtx,
	}

	mockFailtx.EXPECT().CreateFailTx(gomock.Any()).Return(zerror.New(-1, "error")).AnyTimes()

	// error case
	mockCommglobalmap.EXPECT().GetLatestAccountInfo(gomock.Any(), gomock.Any()).Return(nil, zerror.New(-1, "error")).MaxTimes(1)
	req := &globalRPCProto.ReqSendTxByRawInfo{TxInfo: ""}
	_, err := l.SendTransferNftTx(req)
	assert.NotNil(t, err)

	// normal case
	mockCommglobalmap.EXPECT().GetLatestAccountInfo(gomock.Any(), gomock.Any()).Return(&commonAsset.AccountInfo{}, nil).AnyTimes()
	req = &globalRPCProto.ReqSendTxByRawInfo{TxInfo: ""}
	_, err = l.SendTransferNftTx(req)
	assert.Nil(t, err)
}