package pow

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"math"
	"math/bits"

	"github.com/AidosKuneen/cuckoo"
	"golang.org/x/crypto/blake2b"

	"github.com/svalasovich/pow-tcp-server/internal/log"
)

const (
	uint64Size      = 8
	uint8Size       = 1
	uint32Size      = 4
	randomDataSize  = 32
	startComplexity = 170 // experimentally determined that PoW typically terminates with a difficulty of ~170
)

var ErrUnsolved = errors.New("could not solve")

type (
	Engine struct {
		logger *log.Logger
	}
)

func NewEngine() *Engine {
	return &Engine{
		logger: log.NewComponentLogger("pow.engine"),
	}
}

func (e *Engine) Verify(data []byte, nonce []byte) bool {
	if len(data) != randomDataSize+uint8Size || len(nonce) != uint64Size+cuckoo.ProofSize*uint32Size {
		return false
	}

	salt, proofs := deserializeNonce(nonce)

	return verifyAdditionalComplexity(data[0], proofs) && cuckoo.Verify(hashData(append(salt, data...)), proofs) == nil
}

func (e *Engine) Solve(ctx context.Context, data []byte) ([]byte, error) {
	algorithm := cuckoo.NewCuckoo()
	dataWithSalt := make([]byte, uint64Size+len(data))
	copy(dataWithSalt[uint64Size:], data)
	for i := uint64(0); i < math.MaxUint64; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		binary.BigEndian.PutUint64(dataWithSalt, i)
		proofs, solved := algorithm.PoW(hashData(dataWithSalt))
		if !solved || !verifyAdditionalComplexity(data[0], proofs) {
			continue
		}

		return serializeNonce(dataWithSalt[:uint64Size], proofs), nil
	}

	return nil, ErrUnsolved
}

func (e *Engine) GenerateData(complexity uint8) []byte {
	result := make([]byte, uint8Size+randomDataSize)
	result[0] = complexity

	if _, err := rand.Read(result[uint8Size:]); err != nil {
		// rand.Read never returns an error unless the OS has problems; in this case, the application has weak security.
		e.logger.Fatal("data generation failed", "error", err)
	}

	return result
}

func verifyAdditionalComplexity(complexity uint8, proofs []uint32) bool {
	solutionComplexity := 0
	for _, proof := range proofs {
		solutionComplexity += bits.LeadingZeros32(proof)
	}

	return solutionComplexity >= int(complexity)+startComplexity
}

func serializeNonce(salt []byte, proofs []uint32) []byte {
	result := make([]byte, len(proofs)*uint32Size+uint64Size)
	copy(result[:uint64Size], salt)

	resultProofs := result[uint64Size:]
	for _, proof := range proofs {
		binary.BigEndian.PutUint32(resultProofs[:uint32Size], proof)
		resultProofs = resultProofs[uint32Size:]
	}

	return result
}

func deserializeNonce(nonce []byte) ([]byte, []uint32) {
	salt := nonce[:uint64Size]
	nonce = nonce[uint64Size:]

	proofsSize := len(nonce) / uint32Size
	proofs := make([]uint32, proofsSize)
	for i := 0; i < proofsSize; i++ {
		proofs[i] = binary.BigEndian.Uint32(nonce[:uint32Size])
		nonce = nonce[uint32Size:]
	}

	return salt, proofs
}

func hashData(data []byte) []byte {
	result := blake2b.Sum256(data)

	return result[:]
}
