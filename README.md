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

### API

- [GET] /ping

    - response
    
        ```json
        {
            "message": "pong"
        }
        ```

