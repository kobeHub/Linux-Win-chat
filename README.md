# Linux and Win chat tool

A console chat tool between Linux and Windows systems. Written with golang and implement p2p communication via TLS .

## Installing

```shell
go get github.com/kobeHub/Linux-Win-chat
```

## Getting Started

### 1. Generating SSL keys

```shell
bash certs.sh
```

### 2. Build and run

+ Linux

  ```shell
  go build -o peer
  ./peer <user_name> <win_ip>
  ```

  

+ Windows

  ```shell
  go build -o peer.exe
  peer.exe <user_name> <linux_ip>
  ```

  



 

 