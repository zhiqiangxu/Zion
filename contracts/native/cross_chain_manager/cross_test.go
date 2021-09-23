package cross_chain_manager

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"

	scom "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/common"

	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

func TestGetStorageSlot(t *testing.T) {
	txHash, _ := hex.DecodeString("422979a53e8fc8c5646d49d382d0f3eb2212c31d0036bca72d8317648e289161")
	var chainID uint64 = 7
	contract := utils.CrossChainManagerContractAddress
	chainIDBytes := utils.GetUint64Bytes(chainID)
	key := utils.ConcatKey(contract, []byte(scom.REQUEST), chainIDBytes, txHash)
	slot := state.Key2Slot(key[common.AddressLength:])
	t.Logf("\ntxHash : %x\nchainId : %x", txHash, chainIDBytes)
	t.Logf("\nkey : %x\nslot: %x", key, slot)
}
