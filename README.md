Github client testing parameters:
=================================

* work flow detail: https://developer.github.com/v3/guides/basics-of-authentication/#registering-your-app
* Client ID: a4c5ffb112e47562979e
* Client Secret: d0c47635f20d5447ec67090b32f6e73ca0cce50d
* name: marcos-first-oauth2-app
* home url: http://localhost:14002/idx
* call back url: http://localhost:14002/oauth/callback
* local index page: http://localhost:14002/idx

NOTE:
=================================

1. github将用户的登录信息记录在本地cookie中，如果用户已经登录，则直接进行发放授权码的跳转，给出授权码，让客户程序换取token。
1. 如果本地用户没有登录，github内部跳转登录页面，登录成功后发放授权码给客户程序。
1. 如果用户已经登录了，但没有授权；或者，授权的第三方服务行为异常（例如过高频率的请求授权码等）则跳转内部请求接受授权页面。
1. 对用同一个clientID，可能同时存在多个有效token（均能访问API）。 它们派生与不同的授权码。
1. 截止版本 91b06bec 目标可执行文件大小为 15MB（GOOS=linux GOARCH=amd64）。
1. 启动应用需提供环境变量，并且要预先在db server上建立OSIN_DB_DATABASE所指定的库。
	* OSIN_DB_USERNAME（required）
	* OSIN_DB_PASSWORD（optional，default empty）
	* OSIN_DB_HOST（required）
	* OSIN_DB_PORT（optional， default `3306`）
	* OSIN_DB_DATABASE（required）
	* OSIN_TABLE_PREFIX（optional， default `osin_`）
	* GIN_MODE（optional，default `release`)


About Docker Using:
=================================

#for install mysql container instance:


#for start an existing mysql container instance:
docker start mysql-db

docker container ls -a
GOOS=linux GOARCH=amd64 go build .
docker build -t go-app:latest .
docker run -e GIN_MODE=debug -it --link mysql-db:backend  --rm -p 14000:14000 go-app:latest /bin/sh