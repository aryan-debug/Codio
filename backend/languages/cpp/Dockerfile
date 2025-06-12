FROM ubuntu:latest

WORKDIR /app

RUN apt-get update && apt-get install -y \
  g++ \
  libcpprest-dev \
  libboost-all-dev \
  libssl-dev \
  cmake

CMD g++ -o test test.cpp && ./test
