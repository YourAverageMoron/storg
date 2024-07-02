# Distributed Filestore Using Raft

## Purpose of the project
- I started this project to beter understand Golang and distributed systems
- One main goal of the project is to not use third party packages within the production code, the one exception for this is stretchr/testify package for testing

## What is raft?
> Raft offers a generic way to distribute a state machine across a cluster of computing systems, ensuring that each node in the cluster agrees upon the same series of state transitions.
- For more information [see here](https://raft.github.io/)

## Where to start
### Version 1
- To get the full context of this project it is worth starting with V1 [see here](https://github.com/YourAverageMoron/storg/tree/7cc6f6c66be9f1262433b377108abfa11aa1c513)
- This involved creating a simple a distributed file server that provides a client that can replicate files by streaming them to other nodes.
- I created a transport layer using TCP [see here](https://github.com/YourAverageMoron/storg/blob/7cc6f6c66be9f1262433b377108abfa11aa1c513/p2p/tcp_transport.go)
- Used MD5 to encrypt replicated files [see here](https://github.com/YourAverageMoron/storg/blob/7cc6f6c66be9f1262433b377108abfa11aa1c513/crypto.go)
- And created a file server with Store and Get methods for interacting with data [see here](https://github.com/YourAverageMoron/storg/blob/7cc6f6c66be9f1262433b377108abfa11aa1c513/fileserver.go)

### Version 2
- **Please note that version 2 is still a work in progress**
- After version 1 I decided to start from scratch to improve and refactor sections, and integrate the Raft Consensus Algorithm from the ground up [see here](https://github.com/YourAverageMoron/storg/tree/main)
- Version 2 aimed at simplifying the transport layer and decoupling the encoding of messages from the transport layer by hand writing the functions to marshall and unmarshall binary [see here](https://github.com/YourAverageMoron/storg/blob/main/transport/tcp_message.go)
- Encoding and decoding is now handled in the Raft layer allowing for a single encoding instance to be used throught the whole application [see here](https://github.com/YourAverageMoron/storg/blob/89ffc3bce3f03be8c30785ba20e81a681ded5fea/raft/raft.go#L154)

## What I've learnt
- Better understanding of working with tcp
- Better understanding of distributed systems
- Improved my understanding of Golang, IO and channels
- Learned more about Raft and how distributed systems agree on consensus

