# Overview

To meet DDoS protection requirements, add a Proof-of-Work (PoW) algorithm on the server side to make attacks more
expensive. The main task is to implement a PoW algorithm that maximally protects against DDoS attacks. The primary goal
of this document is to choose the PoW algorithm.

# Options

## HashCash

### Overview

Hashcash is a Proof-of-Work scheme built on a cryptographic hash function (originally SHA-1, often SHA-256 today).
A client iterates a nonce until the hash of a token has at least k leading zero bits (~2^k trials).
The server verifies with a single hash, shifting computation cost to the requester (anti-spam/DDoS).

### Pros

- Clear, compact spec and countless examples; verification is a single hash.
- Uses standard hash functions, with libraries and snippets in most languages.

### Cons

- ASIC and GPU friendly

## Argo2id

### Overview

Argon2id is a memory-hard function (hybrid of Argon2i/Argon2d) that requires meaningful RAM and reduces GPU/ASIC
advantages. In a PoW puzzle, the client searches a nonce so that Argon2id(challenge || nonce) yields a digest with k
leading zero bits; the server verifies with one evaluation. Keep m_cost/t_cost moderate and tune difficulty via k so
verification stays cheap.

### Pros

- requires meaningful RAM, shrinking GPU/ASIC advantages and limiting cheap massive parallelism.
- standardized ([RFC 9106](https://www.rfc-editor.org/rfc/rfc9106.html)), widely implemented across languages, and has
  proven stable in production over years.

### Cons

- more CPU time and RAM per attempt; can hurt mobiles/low-end devices and add latency.

## Cuckoo

### Overview

Cuckoo Cycle is a memory-bound Proof-of-Work where the task is to find a cycle of fixed length in a large sparse
bipartite graph deterministically derived from a challenge and nonce (often via SipHash). The cost is dominated by
random memory accesses during edge generation/trimming, while verification is fast by checking that the submitted edges
form the required cycle. Difficulty is tuned by the graph size parameter (edgebits) and/or an acceptance target.

### Pros

- Cuckoo Cycle forces heavy, memory-bound work on the solver, while the server verifies in fixed, minimal time.

### Cons

- No widely accepted spec, limited popularity, and few maintained implementations
- Insufficient real-world production use; resilience and operational robustness are not yet proven

## Merkle Tree Proof

MTP (Merkle Tree Proof) is a memory-hard PoW: the solver fills a large memory array via a deterministic procedure,
builds a Merkle tree over it, and submits the root plus a small set of authenticated samples as the proof.
The verifier checks those samples and Merkle paths instead of recomputing the whole memory, giving fast,
near-constant-time verification.

### Pros

- high memory cost for solvers limits cheap massive parallelism

### Cons

- higher implementation complexity, larger proofs, and fewer standardized, production-grade libraries.
- the ecosystem with implementations lacks maintained, well-documented libraries

# Conclusion

We will use Argon2id with
the [RFC 9106 “Recommended parameters, set #2”](https://www.rfc-editor.org/rfc/rfc9106.html#name-recommendations) (the
profile intended for memory-constrained applications). This choice is driven by its reputation, broad library
availability, and reduced GPU/ASIC advantage. Additionally, we employ adaptive difficulty: when load exceeds defined
thresholds, the required leading-zero bits increase in logarithmic, stepwise increments.

# Note on MTP (Merkle Tree Proof) Schemes

From a theoretical standpoint, MTP schemes are very attractive: they offer the memory-hard advantages we want (akin to
Argon2id) and enable fast, sample-based verification via Merkle proofs. However, public implementations are scarce, not
widely adopted, and lack a strong, battle-tested reputation with proven security. Absent ecosystem and library
constraints, choosing an MTP-based design would be a natural—and likely optimal—option.

# Note on Adaptive Complexity

The current description uses a simplified adaptation. In practice, this can be upgraded to an EMA-smoothed, spike-aware
controller (e.g., two-rate EWMA with α_up > α_down) that raises difficulty quickly on sudden spikes and decreases it
slowly. This prevents low-difficulty windows between bursts and hardens the system against intermittent DDoS patterns.

