FROM golang:1.18-alpine 

RUN mkdir app 

COPY . /app 

WORKDIR /app 

RUN CGO_ENABLED=0 go build -o listenerApp main.go

CMD [ "/app/listenerApp" ]
