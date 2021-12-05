# DTAP - DNSTAP Message router
DTAP is dnstap message router.
Supported multi input and multi output.

## Install
```
go get -u github.com/mimuret/dtap
```

## example
see [example dir](https://github.com/mimuret/dtap/tree/master/example)

## Input config
### Unix Socket
Make unix domain socket for server software writting DNSTAP Frame.
Required parameter `Path` is unix domain socket path,
Optional parameter `User` is socket owner.
```
[[InputUnix]]
Path="/var/run/unbound/dnstap.sock"
User="unbound"
```

### TCP Socket
Make TCP socket for client writting DNSTAP Frame.

Required parameter `Address` is listen address,that default value is `"0.0.0.0"`,
Optional parameter `Port` is listen port, that default value is `10053`.
```
[[InputTCP]]
Address="0.0.0.0"
Port=10053
```

### File
Once read DNSTAP Frame from file.
Can read a compress file gz, bzip2 and xz.

```
[[InputFile]]
Path="/var/dnscap/tap.fstrm.gz"
```

### Tail
Tail read DNSTAP frame from files.
Supported glob format.

```
[[InputTail]]
Path=/var/dnstap/*.fstrm

```

## Output config
### Unix Socket
Write DNSTAP frame to unix domain socket.
If can't open socket, try reconnect interval 1s.
```
[[OutputUnix]]
Path="/var/run/unbound/dnstap.sock"
```

### TCP Socket
Write DNSTAP frame to tcp domain socket.
If can't open socket, try reconnect interval 1s.

```
[[OutputTCP]]
Host="otherhost.exmaple.jp"
Port=10053
```

### File
Write DNSTAP frame to file.
file path supported strftime format for file rotate.
```
[[OutputFile]]
Path = "/var/dnstap/dnstap-%Y%m%d-%H%M.fstrm"
```

### Fluent
Make flatting DNSTAP message,And it forawrd to fluend host.
If can't open socket, try reconnect interval 1s.

Example setting is [here](elasticsearch.md)

```
[[OutputFluent]]
Host = "fluent.example.jp"
Tag  = "dnstap.message"
```

### Kafka
Make flatting DNSTAP message,And it forawrd to kafka host.


```
[[OutputKafks]]
Hosts = ["kafka.example.jp:9092"]
Topic  = "dnstap_message"
```


### Nats
Make flatting DNSTAP message,And it forawrd to nats host.


```
[[OutputNats]]
Host = "nats://kafka.example.jp:5000"
Subject  = "dnstap"
User = "dnstap"
Password = "hogehoge

```
