# Start from Golang v1.13.4 base image to have access to go modules
FROM golang:1.17

# Create a working directory
WORKDIR /app

# we will be expecting to get API_PORT as arguments
ARG API_PORT

# Fetch dependencies on seperate layer as they are less likely to
# Change on every build and will thus be cached
COPY ./go.mod ./
COPY ./go.sum ./

RUN go mod download

# Copy source from the host to the working directory in container
COPY . .

RUN go build -o /api

# EXPOSE PORT
EXPOSE ${API_PORT}

CMD [ "/api" ]
