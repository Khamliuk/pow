# Simple POW using hashcash algorithm

This is example of using [POW (Proof of work)](https://en.wikipedia.org/wiki/Proof_of_work) with [hashcash](https://en.wikipedia.org/wiki/Hashcash) algorithm function

## For run locally please execute those commands:
```
    docker build -t pow-server -f server.Dockerfile .
    docker build -t pow-client -f client.Dockerfile .
    docker-compose up 
```

## Why hashcash algorithm was chosen:

It's my first experience in "blockchain and all related to that" so 
after reading about POW, hashcash and 
other [functions](https://en.wikipedia.org/wiki/Proof_of_work#List_of_proof-of-work_functions) 
I decided to use hashcash because of:
 - It's simplest and most clear for me, I found a lot if docs/guides and even simple [lib](https://github.com/PoW-HC/hashcash)
 - Compare to other functions like [Merkle tree](https://en.wikipedia.org/wiki/Merkle_tree) or [Guided tour puzzle protocol](https://en.wikipedia.org/wiki/Guided_tour_puzzle_protocol)
    - Merkle tree need more calculations on server side, and it's grow with number of leaves and depth
    - For using Guided tour puzzle protocol client should regularly request server about next parts of guide, that complicates logic of contract
 