package stakescan

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// ChecksBase is DecimalChecks.CHECHS_STORAGE_LOCATION
// (= keccak256(abi.encode(uint256(keccak256("decimal.storage.DecimalChecks")) - 1)) & ~0xff).
var ChecksBase = common.HexToHash("0x6135605bbaca49c9a4eeea9432a78c88af8cb013228ed2e2e5c5c2143292b200")

// ChecksStorage field offsets (slot = ChecksBase + offset). Field 0 is _nonces
// (mapping(address => uint256)), not needed for check reconstruction.
const (
	fChecks       = 1 // mapping(bytes32 checkHash => Check)
	fCheckDetails = 2 // mapping(bytes32 checkDetailsHash => CheckDetails)
)

var checkStatusNames = []string{"None", "Redeemed", "Revoked"}
var checkTypeNames = []string{"DEL", "Token"}

// Check is a reconstructed entry of the _checks mapping joined with its CheckDetails.
type Check struct {
	CheckHash        common.Hash    `json:"checkHash"`
	CheckDetailsHash common.Hash    `json:"checkDetailsHash"`
	Signer           common.Address `json:"signer"`
	Status           uint8          `json:"status"`
	StatusName       string         `json:"statusName"`
	// Joined from _checkDetails[checkDetailsHash]:
	TypeChecks   uint8          `json:"typeChecks"`
	TypeName     string         `json:"typeName"`
	Amount       *big.Int       `json:"amount"`
	DueBlock     *big.Int       `json:"dueBlock"`
	Token        common.Address `json:"token"`
	Creator      common.Address `json:"creator"`
	DetailsFound bool           `json:"detailsFound"`
	BaseSlot     common.Hash    `json:"baseSlot"`
}

// CheckHashOf = keccak256(abi.encodePacked(checkDetailsHash, signer)) (see _createCheck).
func CheckHashOf(checkDetails common.Hash, signer common.Address) common.Hash {
	buf := make([]byte, 0, 32+20)
	buf = append(buf, checkDetails.Bytes()...)
	buf = append(buf, signer.Bytes()...)
	return common.BytesToHash(crypto.Keccak256(buf))
}

// ReconstructChecks scans storage for verified _checks entries. Each Check occupies two
// slots: checkDetails(bytes32) at base, and signer(addr)+status(enum @byte 20) packed at
// base+1. checkHash = keccak256(checkDetails ‖ signer) is recomputed from the struct and
// verified against keccak256(checkHash ‖ checksField); the CheckDetails are then joined in
// from _checkDetails[checkDetails].
func ReconstructChecks(s Storage, base common.Hash) []Check {
	checksField := slotAdd(base, fChecks)
	detailsField := slotAdd(base, fCheckDetails)

	var out []Check
	for k := range s {
		w1 := s[slotAdd(k, 1)]
		if !isSignerStatusWord(w1) {
			continue
		}
		checkDetails := s[k]
		signer := common.BytesToAddress(w1.Bytes())
		checkHash := CheckHashOf(checkDetails, signer)
		if mappingSlotBytes32(checkHash, checksField) != k {
			continue
		}
		status := uint8(fieldUint(w1, 20, 1).Uint64())

		c := Check{
			CheckHash:        checkHash,
			CheckDetailsHash: checkDetails,
			Signer:           signer,
			Status:           status,
			StatusName:       enumName(checkStatusNames, status),
			BaseSlot:         k,
		}
		db := mappingSlotBytes32(checkDetails, detailsField)
		c.TypeChecks = uint8(fieldUint(s[db], 0, 1).Uint64())
		c.TypeName = enumName(checkTypeNames, c.TypeChecks)
		c.Amount = s.uintAt(slotAdd(db, 1))
		c.DueBlock = s.uintAt(slotAdd(db, 2))
		c.Token = s.addrAt(slotAdd(db, 3))
		c.Creator = s.addrAt(slotAdd(db, 4))
		c.DetailsFound = c.Creator != (common.Address{}) || c.Amount.Sign() != 0
		out = append(out, c)
	}
	return out
}

// isSignerStatusWord reports whether a word looks like the packed Check {signer, status}
// slot: bytes at offsets 21..31 zero, status (offset 20) in 0..2, signer (low 20) non-zero.
func isSignerStatusWord(w common.Hash) bool {
	for i := 0; i < 11; i++ { // big-endian indices 0..10 == offsets 21..31
		if w[i] != 0 {
			return false
		}
	}
	if w[11] > 2 { // big-endian index 11 == offset 20 (status)
		return false
	}
	for i := 12; i < 32; i++ { // signer (low 20 bytes)
		if w[i] != 0 {
			return true
		}
	}
	return false
}

func enumName(names []string, i uint8) string {
	if int(i) < len(names) {
		return names[i]
	}
	return fmt.Sprintf("Unknown(%d)", i)
}
