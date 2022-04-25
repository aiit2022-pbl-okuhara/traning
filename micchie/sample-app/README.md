# sample-app

## 最初に実施すること
一番最初は `sample-app` データベースが存在しないのでエラーになります。そのため、初回だけコンテナに入って `CREATE DATABASE` をします。

```
$ docker=compose up --build
```

コンテナの中に入る。
```
$ docker ps -a
$ docker exec -it {Container ID} bin/bash
```

`admin` ユーザで `Postgres` に接続する。
```
$ psql -U admin
```

次のクエリを実行する。
```
CREATE DATABASE "sample-app" OWNER = admin TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'ja_JP.UTF-8' LC_CTYPE = 'ja_JP.UTF-8';
```

## How to launch the application
```
$ docker-compose up --build
```

## How to connect the database
After launching your application with docker-compose, connect to the `sample-app-db` container with the `admin` user.
```
$ docker exec -it sampe-app-db psql -U admin
```