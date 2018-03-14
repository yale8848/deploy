# Deploy files and execute commands tool

## Usage

```cmd

deploy -c config.json

```

config.json

execute order: `zipFiles-->preCommands-->uploads-->commands-->deletes`

```json

{

  "concurrency":true,
  "zipFiles":["dir","file"],
  "zipName":"zip.zip",
  "deletes":["zip.zip"],
  "servers":[
    {
      "host":"ip1,ip2",
      "port":22,
      "user":"root",
      "password":"xxxxx",
      "preCommands":[
              ""
       ],      
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

[windows-64](https://github.com/yale8848/deploy/blob/master/release/windows-64/deploy.exe?raw=true)

[linux-64](https://github.com/yale8848/deploy/blob/master/release/linux-64/deploy.exe?raw=true)

[darwin-64](https://github.com/yale8848/deploy/blob/master/release/darwin-64/deploy.exe?raw=true)


## Upload war file demo

```json

{

  "concurrency":true,
  "servers":[
    {
      "host":"ip1",
      "port":22,
      "user":"root",
      "password":"xxxxx",
      "preCommands":[
              ""
       ],
      "uploads":[
        {
          "local":"C:\\javawebdeploy.war",
          "remote":"/coder/tomcat/apache-tomcat-7.0.55/webapps"
        }
      ],
      "commands":[
        "sh /coder/tomcat/apache-tomcat-7.0.55/bin/shutdown.sh",
        "rm -rf /coder/tomcat/apache-tomcat-7.0.55/webapps/javawebdeploy",
        "sh /coder/tomcat/apache-tomcat-7.0.55/bin/startup.sh"
      ]

    }
  ]

}


```



