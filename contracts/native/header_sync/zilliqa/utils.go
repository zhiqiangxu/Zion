/*
 * Copyright (C) 2021 The poly network Authors
 * This file is part of The poly network library.
 *
 * The  poly network  is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The  poly network  is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 * You should have received a copy of the GNU Lesser General Public License
 * along with The poly network .  If not, see <http://www.gnu.org/licenses/>.
 */

package zilliqa

import (
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"github.com/Zilliqa/gozilliqa-sdk/util"
	"github.com/ethereum/go-ethereum/contracts/native"
	scom "github.com/ethereum/go-ethereum/contracts/native/header_sync/common"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	cstates "github.com/polynetwork/poly/core/states"
)

const dsCommKey = "dsComm"

func IsHeaderExist(native *native.NativeContract, hash []byte, chainID uint64) (bool, error) {
	headerStore, err := native.GetCacheDB().Get(utils.ConcatKey(utils.HeaderSyncContractAddress,
		[]byte(scom.HEADER_INDEX), utils.GetUint64Bytes(chainID), hash))
	if err != nil {
		return false, fmt.Errorf("IsHeaderExist, get blockHashStore error: %v", err)
	}
	if headerStore == nil {
		return false, nil
	} else {
		return true, nil
	}
}

func GetTxHeaderByHeight(native *native.NativeContract, height, chainID uint64) (*core.TxBlock, error) {
	latestHeight, err := GetCurrentTxHeaderHeight(native, chainID)
	if err != nil {
		return nil, err
	}

	if height > latestHeight {
		return nil, fmt.Errorf("GetTxHeaderByHeight, height is too big")
	}

	headerStore, err := native.GetCacheDB().Get(utils.ConcatKey(utils.HeaderSyncContractAddress,
		[]byte(scom.MAIN_CHAIN), utils.GetUint64Bytes(chainID), utils.GetUint64Bytes(height)))

	if err != nil {
		return nil, fmt.Errorf("GetTxHeaderByHeight, get blockHashStore error: %v", err)
	}

	if headerStore == nil {
		return nil, fmt.Errorf("GetTxHeaderByHeight, can not find any header records")
	}
	hashBytes, err := cstates.GetValueFromRawStorageItem(headerStore)
	if err != nil {
		return nil, fmt.Errorf("GetHeaderByHeight, deserialize headerBytes from raw storage item err:%v", err)
	}
	return GetTxHeaderByHash(native, hashBytes, chainID)
}

func GetTxHeaderByHash(native *native.NativeContract, hash []byte, chainID uint64) (*core.TxBlock, error) {
	headerStore, err := native.GetCacheDB().Get(utils.ConcatKey(utils.HeaderSyncContractAddress,
		[]byte(scom.HEADER_INDEX), utils.GetUint64Bytes(chainID), hash))
	if err != nil {
		return nil, fmt.Errorf("GetTxHeaderByHash, get blockHashStore error: %v", err)
	}
	if headerStore == nil {
		return nil, fmt.Errorf("GetTxHeaderByHash, can not find any header records")
	}
	storeBytes, err := cstates.GetValueFromRawStorageItem(headerStore)
	if err != nil {
		return nil, fmt.Errorf("GetTxHeaderByHash, deserialize headerBytes from raw storage item err:%v", err)
	}
	var txBlock core.TxBlock
	if err := json.Unmarshal(storeBytes, &txBlock); err != nil {
		return nil, fmt.Errorf("GetTxHeaderByHash, deserialize header error: %v", err)
	}
	return &txBlock, nil
}

func GetCurrentTxHeader(native *native.NativeContract, chainId uint64) (*core.TxBlock, error) {
	height, err := GetCurrentTxHeaderHeight(native, chainId)
	if err != nil {
		return nil, err
	}

	txBlock, err := GetTxHeaderByHeight(native, height, chainId)
	if err != nil {
		return nil, err
	}

	return txBlock, nil
}

func AppendHeader2Main(native *native.NativeContract, height uint64, txHash []byte, chainID uint64) error {
	contract := utils.HeaderSyncContractAddress
	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(scom.MAIN_CHAIN), utils.GetUint64Bytes(chainID), utils.GetUint64Bytes(height)),
		cstates.GenRawStorageItem(txHash))
	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(scom.CURRENT_HEADER_HEIGHT),
		utils.GetUint64Bytes(chainID)), cstates.GenRawStorageItem(utils.GetUint64Bytes(height)))
	scom.NotifyPutHeader(native, chainID, height, util.EncodeHex(txHash))
	return nil
}

func GetCurrentTxHeaderHeight(native *native.NativeContract, chainID uint64) (uint64, error) {
	heightStore, err := native.GetCacheDB().Get(utils.ConcatKey(utils.HeaderSyncContractAddress,
		[]byte(scom.CURRENT_HEADER_HEIGHT), utils.GetUint64Bytes(chainID)))

	if err != nil {
		return 0, fmt.Errorf("GetCurrentTxBlockHeight error: %v", err)
	}

	if heightStore == nil {
		return 0, fmt.Errorf("GetCurrentTxBlockHeight, heightStore is nil")
	}

	heightBytes, err := cstates.GetValueFromRawStorageItem(heightStore)
	if err != nil {
		return 0, fmt.Errorf("GetCurrentTxBlockHeight, deserialize headerBytes from raw storage item err:%v", err)
	}
	return utils.GetBytesUint64(heightBytes), nil
}

func putTxBlockHeader(native *native.NativeContract, txBlock *core.TxBlock, chainID uint64) error {
	contract := utils.HeaderSyncContractAddress
	storeBytes, _ := json.Marshal(txBlock)
	hash := txBlock.BlockHash[:]
	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(scom.HEADER_INDEX), utils.GetUint64Bytes(chainID), hash),
		cstates.GenRawStorageItem(storeBytes))
	scom.NotifyPutHeader(native, chainID, txBlock.BlockHeader.BlockNum, util.EncodeHex(hash))
	return nil
}

func GetDsHeaderByHash(native *native.NativeContract, hash []byte, chainID uint64) (*core.DsBlock, error) {
	headerStore, err := native.GetCacheDB().Get(utils.ConcatKey(utils.HeaderSyncContractAddress,
		[]byte(scom.HEADER_INDEX), utils.GetUint64Bytes(chainID), hash))
	if err != nil {
		return nil, fmt.Errorf("GetDsHeaderByHash, get blockHashStore error: %v", err)
	}
	if headerStore == nil {
		return nil, fmt.Errorf("GetDsHeaderByHash, can not find any header records")
	}
	storeBytes, err := cstates.GetValueFromRawStorageItem(headerStore)
	if err != nil {
		return nil, fmt.Errorf("GetDsHeaderByHash, deserialize headerBytes from raw storage item err:%v", err)
	}
	var dsBlock core.DsBlock
	if err := json.Unmarshal(storeBytes, &dsBlock); err != nil {
		return nil, fmt.Errorf("GetDsHeaderByHash, deserialize header error: %v", err)
	}
	return &dsBlock, nil
}

func putDsBlockHeader(native *native.NativeContract, dsBlock *core.DsBlock, chainID uint64) error {
	contract := utils.HeaderSyncContractAddress
	storeBytes, _ := json.Marshal(dsBlock)
	hash := dsBlock.BlockHash[:]
	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(scom.HEADER_INDEX), utils.GetUint64Bytes(chainID), hash),
		cstates.GenRawStorageItem(storeBytes))
	scom.NotifyPutHeader(native, chainID, dsBlock.BlockHeader.BlockNum, util.EncodeHex(hash))
	return nil
}

func putGenesisBlockHeader(native *native.NativeContract, txBlockAndDsComm TxBlockAndDsComm, chainID uint64) error {
	blockHash := txBlockAndDsComm.TxBlock.BlockHash[:]
	blockNum := txBlockAndDsComm.TxBlock.BlockHeader.BlockNum
	dsBlockNum := txBlockAndDsComm.TxBlock.BlockHeader.DSBlockNum
	contract := utils.HeaderSyncContractAddress
	storeBytes, _ := json.Marshal(&txBlockAndDsComm.TxBlock)
	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(scom.GENESIS_HEADER), utils.GetUint64Bytes(chainID)),
		cstates.GenRawStorageItem(storeBytes))
	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(scom.HEADER_INDEX), utils.GetUint64Bytes(chainID), blockHash),
		cstates.GenRawStorageItem(storeBytes))
	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(scom.MAIN_CHAIN), utils.GetUint64Bytes(chainID), utils.GetUint64Bytes(blockNum)),
		cstates.GenRawStorageItem(blockHash))
	native.GetCacheDB().Put(utils.ConcatKey(contract, []byte(scom.CURRENT_HEADER_HEIGHT),
		utils.GetUint64Bytes(chainID)), cstates.GenRawStorageItem(utils.GetUint64Bytes(blockNum)))
	putDsComm(native, dsBlockNum, txBlockAndDsComm.DsComm, chainID)
	putDsBlockHeader(native, txBlockAndDsComm.DsBlock, chainID)
	scom.NotifyPutHeader(native, chainID, blockNum, util.EncodeHex(blockHash))
	return nil
}

func putDsComm(native *native.NativeContract, blockNum uint64, dsComm []core.PairOfNode, chainID uint64) {
	contract := utils.HeaderSyncContractAddress
	dsbytes, _ := json.Marshal(dsComm)
	native.GetCacheDB().Put(utils.ConcatKey(contract, utils.GetUint64Bytes(chainID), []byte(dsCommKey), utils.GetUint64Bytes(blockNum)), cstates.GenRawStorageItem(dsbytes))
	native.GetCacheDB().Delete(utils.ConcatKey(contract, utils.GetUint64Bytes(chainID), []byte(dsCommKey), utils.GetUint64Bytes(blockNum-1)))
}

func getDsComm(native *native.NativeContract, blockNum uint64, chainID uint64) ([]core.PairOfNode, error) {
	contract := utils.HeaderSyncContractAddress
	dsbytesStore, err := native.GetCacheDB().Get(utils.ConcatKey(contract, utils.GetUint64Bytes(chainID), []byte(dsCommKey), utils.GetUint64Bytes(blockNum)))
	if err != nil {
		return nil, err
	}
	dsbytes, err := cstates.GetValueFromRawStorageItem(dsbytesStore)
	if err != nil {
		return nil, err
	}
	var dsComm []core.PairOfNode
	err = json.Unmarshal(dsbytes, &dsComm)
	if err != nil {
		return nil, err
	}

	return dsComm, nil
}

func dsCommListFromArray(dscomm []core.PairOfNode) *list.List {
	dsComm := list.New()
	for _, ds := range dscomm {
		dsComm.PushBack(ds)
	}
	return dsComm
}

func dsCommArrayFromList(dscomm *list.List) []core.PairOfNode {
	var dsArray []core.PairOfNode
	head := dscomm.Front()
	for head != nil {
		pairOfNode := head.Value.(core.PairOfNode)
		dsArray = append(dsArray, pairOfNode)
		head = head.Next()
	}
	return dsArray
}
