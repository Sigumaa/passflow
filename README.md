# PassFlow

### require

- docker
- go 1.21

**.env**

(local development only)  
```env
REDIS_HOST="localhost"
REDIS_ADDR="6379"
REDIS_DB="0"
REDIS_POOL="1000"
REDIS_PWD=""
ECHO_ADDR="1323"
```

### Usage

```bash
make up
go run .
```

Postmanなどを使うと便利です

流れ

1. /info でユーザー情報を登録
2. /info/:id でユーザー情報が取得できるか確認
3. /pos でユーザーの位置情報を登録

再起動した場合Infoが消えるので1,2をやり直さないとダメ(現状永続化してない やります)

### API

- [GET] /ping

    - response
    
        ```json
        {
            "message": "pong"
        }
        ```

- [POST] /info

    - request
    
        ```json
        {
            "id": "shiyui",
            "name": "しゆい",
            "like": "Go",
            "dislike": "",
            "From": "Ibaraki"
        }
        ```
    
    - response
    
        ```json
        {
            "id": "shiyui",
            "name": "しゆい",
            "like": "Go",
            "dislike": "",
            "From": "Ibaraki"
        }
        ```

- [GET] /info/:id

    - response
    
        ```json
        {
            "id": "shiyui",
            "name": "しゆい",
            "like": "Go",
            "dislike": "",
            "From": "Ibaraki",
            "friends": [
                "sayoi",
                "kyre"
            ]
        }
        ```

- [POST] /pos

    **Infoに登録されているユーザーのみ、postを受け付けられます。**  
    **Infoは現状永続化されていないので、再起動するたびにInfoを登録しなおさないとだめ**  
    **Redisに保存してるid,lat,lonは再起動しても消えません。(redis-dataみたいなdirectoryが生成されるはず)**  

    - request
    
        ```json
        {
            "id": "shiyui",
            "lat": 1,
            "lng": 1
        }
        ```
    
    - response
    
        ```json
        {
            "cnt": 2,
            "users": [
                {
                    "id": "sayoi",
                    "lat": 1.1,
                    "lon": 1.1
                },
                {
                    "id": "kyre",
                    "lat": 1.2,
                    "lon": 1.2
                }
            ]
        }
        ```