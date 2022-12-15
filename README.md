# Mini Twitter App

This is a clone of Twitter with a minimal set of features built using GO and distributed storage using RAFT. 

Technical specifications:

1. A Web Server to serve user requests using HTTP APIs. No state is stored in this service.

2. 3 backend services (namely authentication service, users service, tweets service) which talk to Web Server through GRPC in order to fulfill users requests. 

3. The backend services persist all data onto RAFT backed storage using etcd.

---

> ### Currently supported features:
> 1. Create an account using Username, Name, and Password
> 2. Login with Username and Password
> 3. Follow and Unfollow other users
> 4. Post Tweets
> 5. View own Tweets
> 6. View TImeline with Tweets from people you follow
> 7. Multiple users login using sessions


---

## Instructions To Run App

The steps need to be followed in order.

### Setup

Firstly, [Install Go](https://go.dev/doc/install)

Then, install Goreman 
``` 
go install github.com/mattn/goreman@latest 
``` 

Then start from root of project and ensure go dependencies are satisfied 
``` 
    cd ./web
    go mod download 
``` 

### Start Raft

Start from root of project
```bash
    cd ./raft/src/go.etcd.io/etcd/contrib/raftexample

    go build -o raftexample

    goreman start
```

### Run Tests
Start from root of project. 

```bash
    cd ./web/users
    go test 

    cd ./web/tweets
    go test 

    cd ./web/authentication
    go test 

```
### Start the Web Server
Start from root of project. 

```bash
    go run ./web/web.go
```

### Start the Backend Services
Start from root of project.  

```bash
    go run ./web/server.go
```

> The application should be available at https://ide8000.anubis-lms.io/ or http://localhost:8080/ if running locally. 


