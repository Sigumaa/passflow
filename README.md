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

- [GET] /user/:id

    - response
    
        ```json
        [
            {
                "id":"shiyui",
                "lat":1, 
                "lon":1
            },
            ...
        ]
        ```
    
    or

    - response
    
        ```json
        {
            "message": "No Users around you"
        }
        ```

- [POST] /user

    - request
    
        ```json
        {
            "id":"shiyui",
            "lat":1, 
            "lon":1
        }
        ```
    
    - response
    
        ```json
        {
            "id":"shiyui",
            "lat":1,
            "lon":1
        }
        ```


