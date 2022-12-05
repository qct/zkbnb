package tree

import (
	"crypto/sha256"
	"fmt"
	bsmt "github.com/bnb-chain/zkbnb-smt"
	"github.com/bnb-chain/zkbnb-smt/database/memory"
	"github.com/zeromicro/go-zero/core/logx"
	"hash"
	"testing"
)

func Test_fsdf(t *testing.T) {
	assetCacheSize := 1
	accountNums := int64(10)
	blockHeight := int64(1)

	hasher := bsmt.NewHasherPool(func() hash.Hash { return sha256.New() })

	accountAssetTrees := NewLazyTreeCache(assetCacheSize, accountNums-1, blockHeight, func(index, block int64) bsmt.SparseMerkleTree {
		tree, err := bsmt.NewBASSparseMerkleTree(hasher,
			memory.NewMemoryDB(), AssetTreeHeight, NilAccountAssetNodeHash,
			bsmt.GCThreshold(1024*10))
		if err != nil {
			logx.Errorf("unable to create new tree by assets: %s", err.Error())
			panic(err.Error())
		}
		return tree
	})

	a11 := accountAssetTrees.Get(0)
	a12 := accountAssetTrees.Get(0)
	a21 := accountAssetTrees.Get(1)
	a22 := accountAssetTrees.Get(1)
	fmt.Println(a11)
	fmt.Println(a12)
	fmt.Println(a21)
	fmt.Println(a22)
}
