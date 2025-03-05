set positional-arguments
default:
  just --list

alias s := server
alias t := tidy

server:
  air

tidy:
  go mod tidy

[no-cd]
@proto-gen bar:
  pwd
  protoc --go_out=./pb --go-grpc_out=./pb $1.proto

up:
  devenv up
