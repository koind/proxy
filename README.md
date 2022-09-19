# proxy
Simple service for proxying HTTP requests to third-party services.

## How to use

To run the required command in terminal:
```bash
make run
```
---

To check the methods and view the documentation, follow the link:
```bash
http://localhost:8080/swagger/index.html
```
---

To stop the service, run the following command in the terminal:
```bash
make down
```

## Request samples

cURL samples can be run from command line

### /

```sh
curl -X POST -H "Content-Type: application/json" -d '{ "method": "GET", "url": "http://google.com", "headers": { "Authentication": "Basic bG9naW46cGFzc3dvcmQ=&quot;," } }' http://localhost:8080
```

### /ping

```sh
curl -X GET 'http://localhost:8080/ping'
```

### /history

```sh
curl -X GET 'http://localhost:8080/history'
```
