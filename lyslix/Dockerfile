# build the executable
FROM golang:alpine AS build
ADD . /src
RUN cd /src && go build -o lyslix

# prepare the service
FROM scratch
WORKDIR /server
COPY --from=build /src/lyslix /server
EXPOSE 56700
ENV PORT 56700
ENTRYPOINT ['./lyslix']
