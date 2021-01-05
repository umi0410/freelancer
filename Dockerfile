FROM golang:1.14
COPY . /freelancer
WORKDIR /freelancer
RUN go get .
EXPOSE 1323
ENTRYPOINT ["freelancer"]