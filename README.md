# Upload files and execute commands with ssh and sftp

## Execute use

```cmd
all_build.cmd
```

```cmd

deploy -c config.json

```

config.json

execute order: `zipFiles-->uploads-->commands-->deletes`

```json

{

  "concurrency":true,
  "zipFiles":["dir","file"],
  "zipName":"zip.zip",
  "deletes":["",""],
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


## Download

[windows-64](./release/windows-64/deploy.exe)

[linux-64](./release/linux-64/deploy)

[darwin-64](./release/darwin-64/deploy)



