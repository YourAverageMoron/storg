# Content-addressible storage in GO :)

## What I've learnt
- Working more with golang and concurrency in go
    - First bug project in go
- Better understanding of working with tcp
- Better understanding of distributed systems


## Next steps after this 
- https://fly.io/dist-sys/ - This looks like it could build on some of the stuff in this prject might be good to look into
- Look into using RAFT consensus - https://raft.github.io/
- Better marshalling of tcp payloads
    - Payloads should have
        - version
        - request method/type (get/store etc)
        - payload size
        - payload body
    - Have 2 functions:
        - marshal
        - unmarshal
    - Peer descovery
- Public/Private key encryption
- Delete files from remote servers
- Dockerize it
    - Create a docker compose file to spin up some test servers
- Look more into file transfer protocol FTP
- can it implement stuff like rsync
