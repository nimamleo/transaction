package infrastructure

import (
	"encoding/hex"

	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

func uint128ToString(id types.Uint128) string {
	return hex.EncodeToString(id[:])
}

func stringToUint128(s string) (types.Uint128, error) {
	bytes, err := hex.DecodeString(s)
	if err != nil {
		return types.Uint128{}, err
	}
	var id types.Uint128
	copy(id[:], bytes)
	return id, nil
}
