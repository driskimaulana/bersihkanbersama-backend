#Compile stage
FROM golang:1.19.5-alpine AS compiler

# Add required packages
ENV MONGOURI="mongodb+srv://driskimaulana:maulanadriski77@cluster0.opvmrgn.mongodb.net/?retryWrites=true&w=majority"
ENV API_SECRET=asdfghjkl
ENV TOKEN_HOUR_LIFESPAN=1
WORKDIR /app
ADD go.mod ./
RUN go mod download
ADD . .
RUN go build -o goapp

# Run stage
FROM alpine:3.16
WORKDIR /usr/src/app
COPY --from=compiler /app/goapp .
CMD ["/usr/src/app/goapp"]