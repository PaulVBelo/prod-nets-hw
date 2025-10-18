DNS-Server для A записей

Самый простой способ запустить - собрать бинарник, после чего запустить его с sudo:

Из корня этого проекта (dns-server):

```
go build -o dns-server ./cmd/main.go
sudo ./dns-server
```

Для конфига используется config.yaml из корня этого проекта.