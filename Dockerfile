################################################################################
# Copyrights © 2021-2022 Fiserv, Inc. or its affiliates. 
# Fiserv is a trademark of Fiserv, Inc., 
# registered or used in the United States and foreign countries, 
# and may or may not be registered in your country.  
# All trademarks, service marks, 
# and trade names referenced in this 
# material are the property of their 
# respective owners. This work, including its contents 
# and programming, is confidential and its use 
# is strictly limited. This work is furnished only 
# for use by duly authorized licensees of Fiserv, Inc. 
# or its affiliates, and their designated agents 
# or employees responsible for installation or 
# operation of the products. Any other use, 
# duplication, or dissemination without the 
# prior written consent of Fiserv, Inc. 
# or its affiliates is strictly prohibited. 
# Except as specified by the agreement under 
# which the materials are furnished, Fiserv, Inc. 
# and its affiliates do not accept any liabilities 
# with respect to the information contained herein 
# and are not responsible for any direct, indirect, 
# special, consequential or exemplary damages 
# resulting from the use of this information. 
# No warranties, either express or implied, 
# are granted or extended by this work or 
# the delivery of this work
################################################################################

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