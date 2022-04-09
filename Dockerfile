FROM golang

WORKDIR /usr/srcdock/app

COPY ./ ./
RUN go mod tidy
RUN rm -rf migration/
CMD [ "go", "run", "main.go" ]
