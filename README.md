# PassFlow

### require

- docker
- go 1.21

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


