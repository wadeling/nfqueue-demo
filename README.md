# nfqueue-demo

A demo for nfqueue detect network for containers

## usage

### first start 2 contaienrs

```
./script/start_test_container.sh
```

### inspect docker pid

```
docker inspect -f '{{.State.Pid}}' wade-test-1
docker inspect -f '{{.State.Pid}}' wade-test-2
```

### change src
modify main.go, change pid to above contaier pid

### build src

```
./build.sh
```

### run demo

```
./run.sh
```

### test

#### install iptables for container

```
docker exec -it wade-test-1 bash
apt install iptables
iptables -t raw -I OUTPUT -p tcp --syn -j NFQUEUE --queue-num=1 --queue-bypass
```

#### test curl

```
docker exec -it wade-test-1 bash
curl -v http://www.baidu.com
```

### see result

in terminal that running nfqueue-demo,you will see

```
2023/01/03 17:31:58 start
2023/01/03 17:31:58 do container wade-test-2,pid 1565,que num 2
2023/01/03 17:31:58 new nf queue 2 end
2023/01/03 17:31:58 do container wade-test-1,pid 1494,que num 1
2023/01/03 17:31:58 new nf queue 1 end
2023/01/03 17:33:06 A new tcp connection will be established: 172.17.0.3:41640 -> 14.215.177.39:80
```

## attention

we need "CAP_NET_ADMIN" to run nfqueue-demo. You can see in run.sh

```
setcap 'cap_net_admin=+ep' ./nfqueue-demo
```

but if you run nfqueue-demo in a vm which mount source code directory to host directory, you may encounter an error like:

```
Failed to set capabilities on file `./nfqueue-demo' (Operation not supported)
```

You should copy the binary to /usr/local/bin or /usr/bin, the reason [reason](https://stackoverflow.com/questions/29099797/raw-capture-capabilities-cap-net-raw-cap-net-admin-not-working-outside-usr-b)


