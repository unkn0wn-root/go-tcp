# Simple TCP Server

Made two implementations:
1. A standard implementation using Go's net package
2. A raw implementation using system calls

## Usage

### Running the Standard Server

```bash
make run-standard
```

### Running the Raw TCP Server

```bash
make run-custom
```

### Testing the Servers

You can test either server using netcat:

```bash
nc localhost 8080
```

Type any message and press enter. The server will echo back your message.
Type "quit" to close the connection.
