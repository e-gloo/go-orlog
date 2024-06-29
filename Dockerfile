FROM golang:1.22.1-alpine

# Run as a non-privileged user
RUN addgroup app &&\
    adduser --ingroup app --disabled-password app
USER app

WORKDIR /app

# Install dependencies
RUN go install github.com/go-task/task/v3/cmd/task@latest
COPY --chown=app:app go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source files into application directory
COPY --chown=app:app . .

# Run app external dependencies check
RUN task check-all

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]
