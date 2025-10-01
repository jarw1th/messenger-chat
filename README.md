# Messenger Chat API

A real-time chat server built in Go with WebSocket support, channels, and private messaging. Stores messages in PostgreSQL and supports fetching message history.

---

## Features

- Real-time messaging using WebSockets.
- Public channels and private messages.
- Message history stored in PostgreSQL.
- Messages include sender and receiver names.
- Supports multiple channels with unique names.
- Simple, scalable architecture with hub-client model.

---

## Tech Stack

- **Language:** Go
- **Database:** PostgreSQL
- **WebSocket Library:** Gorilla WebSocket
- **Hosting:** Any Go-compatible server

---

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL

### Installation

1. Clone the repository:

```bash
git clone https://github.com/your-username/messenger-chat.git
cd messenger-chat
```

2.	Set up PostgreSQL and create the database:

```bash
CREATE DATABASE chatdb;
```

3.	Create tables and insert initial data:

```bash
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL
);

CREATE TABLE channels (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    channel_id INT REFERENCES channels(id) ON DELETE CASCADE,
    sender_id INT REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE private_messages (
    id SERIAL PRIMARY KEY,
    sender_id INT REFERENCES users(id) ON DELETE CASCADE,
    receiver_id INT REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

INSERT INTO users (username) VALUES ('Alice'), ('Bob') ON CONFLICT (username) DO NOTHING;
INSERT INTO channels (name) VALUES ('General'), ('Random') ON CONFLICT (name) DO NOTHING;
```

4.	Install Go dependencies:

```bash
go mod tidy
```

5.	Run the server:

```bash
go run main.go
```

Server will start on http://localhost:8080.

---

## WebSocket API

### Connect

```bash
ws://localhost:8080/ws?username=<username>&channel_name=<channel_name>
```
•	username: required, default guest
•	channel_name: optional, default General

### Sending Messages

- **Channel message** (broadcast to all users in the channel):

```json
{
  "content": "Hello everyone"
}
```

- **Private message** (send to a specific user):

```json
{
  "content": "Hi Bob",
  "receiver_id": 2
}
```

### Receiving Messages

Server sends messages as JSON:

```json
{
  "id": 1,
  "channel_id": 1,
  "sender_id": 1,
  "sender_name": "Alice",
  "receiver_id": 2,
  "receiver_name": "Bob",
  "content": "Hello Bob",
  "created_at": "2025-10-01T22:50:10Z"
}
```

•	id: unique message ID
•	channel_id: ID of the channel (0 if private)
•	sender_id: ID of the sender
•	sender_name: username of the sender
•	receiver_id: ID of the receiver (0 if channel message)
•	receiver_name: username of the receiver (empty if channel message)
•	content: text of the message
•	created_at: timestamp of when the message was sent

---

## Client Usage Example

Run the mock client for a channel:

```bash
go run mockTest.go Alice General
```

Send a private message to user with ID 2:

```bash
go run mockTest.go Alice General 2
```

All messages sent to a channel are broadcast to every connected client in that channel.

---

## Notes

- Server automatically creates channels if they do not exist.
- History of the last 50 messages is sent to a client on connect.
- Users cannot see messages from other channels.
- Hub-client architecture allows scaling to multiple channels and clients.

---

## Contributing

1. Fork the repository.  
2. Create a feature branch (`git checkout -b feature-name`).  
3. Commit your changes (`git commit -m 'Add feature'`).  
4. Push to the branch (`git push origin feature-name`).  
5. Open a Pull Request.

---

## License

MIT License © Ruslan Parastaev