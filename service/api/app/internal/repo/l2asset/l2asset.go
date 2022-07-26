package l2asset

import (
	"context"
	"fmt"
	table "github.com/bnb-chain/zkbas/common/model/assetInfo"
	"github.com/bnb-chain/zkbas/pkg/multcache"
	"gorm.io/gorm"
	"strconv"
)

type l2asset struct {
	table string
	db    *gorm.DB
	cache multcache.MultCache
}

/*
	Func: GetL2AssetsList
	Params:
	Return: err error
	Description: create account table
*/
func (m *l2asset) GetL2AssetsList(ctx context.Context) ([]*table.AssetInfo, error) {
	f := func() (interface{}, error) {
		res := []*table.AssetInfo{}
		dbTx := m.db.Table(m.table).Find(&res)
		if dbTx.Error != nil {
			return nil, dbTx.Error
		}
		if dbTx.RowsAffected == 0 {
			return nil, ErrNotFound
		}
		return &res, nil
	}
	res := []*table.AssetInfo{}
	value, err := m.cache.GetWithSet(ctx, multcache.KeyGetL2AssetsList, &res, 10, f)
	if err != nil {
		return nil, err
	}
	res1, ok := value.(*[]*table.AssetInfo)
	if !ok {
		return nil, fmt.Errorf("[GetL2AssetsList] ErrConvertFail")
	}
	return *res1, nil
}

/*
	Func: GetL2AssetInfoBySymbol
	Params: symbol string
	Return: res *L2AssetInfo, err error
	Description: get l2 asset info by l2 symbol
*/
func (m *l2asset) GetL2AssetInfoBySymbol(ctx context.Context, symbol string) (*table.AssetInfo, error) {
	f := func() (interface{}, error) {
		res := table.AssetInfo{}
		dbTx := m.db.Table(m.table).Where("asset_symbol = ?", symbol).Find(&res)
		if dbTx.Error != nil {
			return nil, dbTx.Error
		}
		if dbTx.RowsAffected == 0 {
			return nil, ErrNotExistInSql
		}
		return &res, nil
	}
	res := table.AssetInfo{}
	value, err := m.cache.GetWithSet(ctx, multcache.KeyGetL2AssetInfoBySymbol+symbol, &res, 10, f)
	if err != nil {
		return nil, err
	}
	res1, ok := value.(*table.AssetInfo)
	if !ok {
		return nil, fmt.Errorf("[GetL2AssetInfoBySymbol] ErrConvertFail")
	}
	return res1, nil
}

/*
	Func: GetSimpleL2AssetInfoByAssetId
	Params: assetId uint32
	Return: L2AssetInfo, error
	Description: get layer-2 asset info by assetId
*/
func (m *l2asset) GetSimpleL2AssetInfoByAssetId(ctx context.Context, assetId uint32) (*table.AssetInfo, error) {
	f := func() (interface{}, error) {
		res := table.AssetInfo{}
		dbTx := m.db.Table(m.table).Where("asset_id = ?", assetId).Find(&res)
		if dbTx.Error != nil {
			return nil, dbTx.Error
		}
		if dbTx.RowsAffected == 0 {
			return nil, ErrNotFound
		}
		return &res, nil
	}
	res := table.AssetInfo{}
	value, err := m.cache.GetWithSet(ctx, multcache.KeyGetSimpleL2AssetInfoByAssetId+strconv.Itoa(int(assetId)), &res, 10, f)
	if err != nil {
		return nil, err
	}
	res1, ok := value.(*table.AssetInfo)
	if !ok {
		return nil, fmt.Errorf("[GetSimpleL2AssetInfoByAssetId] ErrConvertFail")
	}
	return res1, nil
}
