# Mini Twitter App

This is a clone of Twitter with a minimal set of features built using GO and distributed storage using RAFT. 

Technical specifications:

1. A Web Server to serve user requests using HTTP APIs. No state is stored in this service.

2. 3 backend services (namely authentication service, users service, tweets service) which talk to Web Server through GRPC in order to fulfill users requests. 

3. The backend services persist all data onto RAFT backed storage using etcd. 
   We are using etcd's implementation of RAFT where data is stored in a Key-Value map which is replicated and kept consistent using RAFT in the background.  For each microservice of our application, we use a Key to serialize data as JSON and store using the RAFT client.

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

Firstly, [Install Go](https://go.dev/doc/install) and setup Go PATH

Then run following script to build RAFT and Go modules for our project 

``` 
    cd cmd
    ./setup.sh
``` 

### Start Raft

This will cleanup old data and launch 3 member RAFT cluster

```bash
    ./run_raft.sh
```

### Run Tests

This will seed initial data and run Go Tests for all services.
In a new terminal, start from root of project and run

```bash
    cd cmd
    ./run_tests.sh
```

### Start the Web Server
In a new terminal, start from root of project and run

```bash
    cd ./web
    go run ./web.go
```

### Start the Backend Services
In a new terminal, start from root of project and run

```bash
    cd ./web
    go run ./server.go
```
> The application should be available at https://ide8000.anubis-lms.io/ or http://localhost:8080/ if running locally. 

> Note: To Stop RAFT once done 
> ```bash
>     goreman run stop-all
> ```