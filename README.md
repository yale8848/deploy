# 通过SSH连接，简单的配置文件，让部署更简单

## 用法

```cmd

deploy -c config.json

```

config.json

配置执行顺序: `preCommands-->uploads-->commands-->verify`

```json
{
  "concurrency":true,
  "servers":[
    {
      "host":"ip1,ip2",
      "port":22,
      "user":"root",
      "password":"xxxx",
      "userPasswordPath": "/path/to/usr.password",
      "socks5UrlPath": "/path/to/socks.url",
      "privateKeyPath":"/path/to/key.pem",
      "preCommands":[
        "mkdir /home/server"
      ],
      "uploads":[
        {
          "local":["resource","start.sh","G:\\tmp\\mylog.txt"],
          "zipRegexp":["test/bb.txt","ccc$"],
          "remote":"/home/server"
        }
      ],
      "commands":[
       "sh /home/server/start.sh"
      ],
      "verify":{
           "http":"http",
           "delay":3,
           "gap":2,
           "path":":8080/api/appInfo",
           "count":3,
           "successStrFlag":"1.10"
      }
    }
  ]

}
```

配置介绍：

```
{
  "concurrency":true, //是否并发执行，默认false
  "servers":[
    {
      "host":"ip1,ip2",//服务器ip或域名，多个以逗号分隔
      "port":22,//默认22，可以不填
      "user":"root",//
      "password":"xxxx",//
      "userPasswordPath": "usr.password",//可选，有这个文件时可以不填上面的user和password字段；文件内容格式usr@password，
      "privateKeyPath":"/path/to/key.pem",//私钥路径，可选
      "socks5UrlPath": "/path/to/socks.url",//socks5代理，可选，文件内容如：socks5://usr:password@ip:port?tcpTimeout=60&tcpDeadline=0&udpDeadline=60
      "privateKeyPath":"/path/to/key.pem",//私钥地址，可选，如果不想输入密码可以填私钥地址
      "preCommands":[ //上传文件前执行服务器命令
        "mkdir /home/server"
      ],
      "uploads":[//上传文件配置
        {
          "local":["resource","start.sh","G:\\tmp\\mylog.txt"],//本地要上传的目录和文件列表，上传时会打包为一个zip文件；上传文件路径为执行deploy命令目录的相对路径或者绝对路径
          "zipRegexp":["test/bb.txt","ccc$"],//zip打包过滤正则，从local中过滤符合正则条件的文件
          "remote":"/home/server" //要上传的服务器路径
        }
      ],
      "commands":[ //上传后执行服务器命令
       "sh /home/server/start.sh"
      ],
      "verify":{ //上传完后给服务器接口发送http get请求来验证是否部署成功
           "http":"http", //http或https，默认http
           "delay":3,     //上传完文件后延迟多长时间发送请求,默认3秒
           "count":3,     //轮询次数，默认3次
           "gap":2,       //轮询间隔时间，默认2秒
           "path":":8080/api/appInfo", //接口path，会和上面的host组成完整url
           "successStrFlag":"1.10"     //验证返回数据是否包含字符串，以此来判定部署成功
      }
    }
  ]

}
```
## 下载

[windows-64](https://github.com/yale8848/deploy/blob/master/release/windows-64/deploy.exe?raw=true)

[linux-64](https://github.com/yale8848/deploy/blob/master/release/linux-64/deploy.exe?raw=true)

[darwin-64](https://github.com/yale8848/deploy/blob/master/release/darwin-64/deploy.exe?raw=true)


## 上传war包例子

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
          "local":["target\\javawebdeploy.war"],
          "remote":"/coder/tomcat/apache-tomcat-7.0.55/webapps"
        }
      ],
      "commands":[
        "sh /coder/tomcat/apache-tomcat-7.0.55/bin/shutdown.sh",
        "rm -rf /coder/tomcat/apache-tomcat-7.0.55/webapps/javawebdeploy",
        "sh /coder/tomcat/apache-tomcat-7.0.55/bin/startup.sh"
      ],
      "verify":{
          "path":":8080/api/info",
          "successStrFlag":"1.10"
       }
    }
  ]

}


```



