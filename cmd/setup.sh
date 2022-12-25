cd ../raft/src/go.etcd.io/etcd/contrib/raftexample
go install github.com/mattn/goreman@latest
go build -o raftexample

cd ../../../../../../web
go mod tidy
go mod download

echo "Setup completed!"