package eip712

import (
	"github.com/FlareZone/melon-backend/common/hexutil"
	"github.com/FlareZone/melon-backend/common/signature"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

/**

function permit(
	address owner,
	address spender,
	uint256 value,
	uint256 deadline,
	uint8 v,
	bytes32 r,
	bytes32 s
)
bytes32 private constant PERMIT_TYPEHASH =
	keccak256("Permit(address owner,address spender,uint256 value,uint256 nonce,uint256 deadline)");


bytes32 structHash = keccak256(abi.encode(PERMIT_TYPEHASH, owner, spender, value, _useNonce(owner), deadline));

       bytes32 hash = _hashTypedDataV4(structHash);

       address signer = ECDSA.recover(hash, v, r, s);
*/

func TestNewPermitToken(t *testing.T) {
	var (
		owner     = "0x3415c8EfF9f4066FA7c4e89526735C866B606081"
		spender   = "0x5AD07396DCe3a6aA66eeAEA2f2207D81e35Ba48a"
		amount    = new(big.Int).SetUint64(100000000)
		nonceData = new(big.Int).SetUint64(1)
		timestamp = new(big.Int).SetInt64(1706775163)
	)

	permitToken := NewPermitToken("Melon", "1", math.NewHexOrDecimal256(int64(1337)), "0xdEdDA3fD0ACd5860a362c260D659ec09148B7757")
	permitToken.
		WithMessage(owner,
			spender,
			amount,
			nonceData,
			timestamp)
	sign, err := permitToken.Sign(
		signature.GetPrivateKeyFromPrivateHex(hexutil.StringToHex("0xa8344e3d88588d3d1cd587f72b94130566d3b4bb0ce02d71192c372df8c24b76")))
	assert.NoError(t, err)
	t.Log("sign: ", sign)
}
