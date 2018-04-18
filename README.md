# Tinchain
A tiny blockchain

## Feature
- Consensus
    1. POW
    2. POW + DPOS
    3. Algorand maybe?
- Transaction
    - obey to EVM gas-using rules but not actually use gas
    - fix gas limit of contract (virtually)

## Tech stack
- Libp2p
- use Ethereum VM
    - support Solidity
    - will support jsvm in the future
- IPFS maybe?