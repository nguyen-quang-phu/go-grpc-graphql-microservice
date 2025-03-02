default:
  just --list

alias s := server
alias t := tidy

server:
  air

tidy:
  go mod tidy

up:
  devenv up
