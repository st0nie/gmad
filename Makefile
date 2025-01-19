TARGET = gmad
SRC = main.go

all: $(TARGET)

$(TARGET): $(SRC)
	GOOS=linux GOARCH=amd64 go build -o $(TARGET) $(SRC)

clean:
	rm -f $(TARGET)