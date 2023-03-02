package utils

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gcash/bchd/bchec"
	"github.com/gcash/bchd/chaincfg"
	"github.com/gcash/bchd/chaincfg/chainhash"
	"github.com/gcash/bchd/wire"
	"github.com/gcash/bchutil"
	"github.com/gcash/bchutil/base58"
)

func VerifySignature(address string, message string, signature string) (bool, error) {
	addr, err := bchutil.DecodeAddress(address, &chaincfg.MainNetParams)
	if err != nil {
		fmt.Println("Error decoding address:", err)
		return false, err
	}
	legacyAddr := base58.CheckEncode(addr.ScriptAddress(), chaincfg.MainNetParams.LegacyPubKeyHashAddrID)
	if _, ok := addr.(*bchutil.AddressPubKeyHash); !ok {
		return false, errors.New("address is not a p2pkh address")
	}
	// Decode the signature and address from base64 to bytes
	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		fmt.Println("Failed to decode signature:", err)
		return false, err
	}

	var buf bytes.Buffer
	wire.WriteVarString(&buf, 0, "Bitcoin Signed Message:\n")
	wire.WriteVarString(&buf, 0, message)
	expectedMessageHash := chainhash.DoubleHashB(buf.Bytes())
	pk, wasCompressed, err := bchec.RecoverCompact(bchec.S256(), sig,
		expectedMessageHash)
	if err != nil {
		return false, err
	}
	var serializedPK []byte
	if wasCompressed {
		serializedPK = pk.SerializeCompressed()
	} else {
		serializedPK = pk.SerializeUncompressed()
	}
	newAddress, err := bchutil.NewAddressPubKey(serializedPK, &chaincfg.MainNetParams)
	if err != nil {
		return false, err
	}
	return newAddress.EncodeAddress() == legacyAddr, nil
}
