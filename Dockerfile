FROM golang:1.15

#RUN mkdir -p /go/src/app
#WORKDIR /go/src/app
#COPY ./src/app /go/src/app
#EXPOSE 8080
#RUN go get "github.com/gorilla/mux"
#RUN go build -o app ./
#CMD ["./app"]

# Move to working directory /build
RUN mkdir -p /build
WORKDIR /build

RUN go get "github.com/gorilla/mux"
RUN go get "github.com/ghodss/yaml"
RUN go get "github.com/kelseyhightower/envconfig"
RUN go get "gopkg.in/yaml.v2"

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o main .

#Run Unit Tests
#RUN go test -v main_test.go

RUN mkdir -p /dist
# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/main .

RUN mkdir resources
RUN cp /build/resources/config.yml ./resources/config.yml

EXPOSE 8080

# Command to run when starting the container
CMD ["/dist/main"]