package pow

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"math"
	"math/bits"

	"golang.org/x/crypto/argon2"

	"github.com/svalasovich/pow-tcp-server/internal/log"
)

const (
	uint8Size             = 1
	randomDataSize        = 32
	parallelism    uint8  = 5
	keyLen         uint32 = 32
	timeCost              = 3
	memoryKiB             = 64 * 1024 // 64 MiB
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
	complexity := data[0]
	hash := argon2.IDKey(data, nonce, timeCost, memoryKiB, parallelism, keyLen)

	return verifyComplexity(complexity, hash)
}

func (e *Engine) Solve(ctx context.Context, data []byte) ([]byte, error) {
	complexity := data[0]

	for i := uint64(0); i < math.MaxUint64; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		nonce := binary.BigEndian.AppendUint64(nil, i)
		hash := argon2.IDKey(data, nonce, timeCost, memoryKiB, parallelism, keyLen)
		if verifyComplexity(complexity, hash) {
			return nonce, nil
		}
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

func verifyComplexity(complexity uint8, hash []byte) bool {
	zeroBits := 0
	for _, x := range hash {
		if x != 0 {
			zeroBits += bits.LeadingZeros8(x)
			break
		}
		zeroBits += 8
	}

	return int(complexity) <= zeroBits
}
