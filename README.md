# Upload files and execute commands with ssh and sftp

## Hot use

```go
go get github.com/yale8848/sshclient
```


```go

sshclient -c config.json

```

config.json

```json

{

  "concurrency":true,
  "zipFiles":["dir","file"],
  "zipName":"zip.zip",
  "servers":[
    {
      "host":"ip1,ip2",
      "port":22,
      "user":"root",
      "password":"xxxxx",
      "uploads":[
        {
          "local":"G:\\tmp\\mylog.txt",
          "remote":"/home/soft"
        }
      ],
      "commands":[
        "date","uname","date"
      ]

    }
  ]

}
```
