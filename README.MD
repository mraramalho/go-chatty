```
    ____            ____ _           _   _
   / ___| ___      / ___| |__   __ _| |_| |_ _   _
  | |  _ / _ \____| |   | '_ \ / _' | __| __| | | |
  | |_| | (_) |___| |___| | | | (_| | |_| |_| |_| |
   \____|\___/     \____|_| |_|\__,_|\__|\__|\__, |
                                             |___/
   A simple TCP chat application
   ---------------------------
```

# 📢 Go-Chatty - A Secure TCP Chat

Go-Chatty is a TCP-based chat application focused on **privacy**, **security**, and **autonomy**. It enables direct communication between clients without the need for intermediary servers that could store or monitor messages, ensuring a confidential and secure information exchange environment.

## 🚀 Main Features

✅ **End-to-end encryption (AES-256)** to protect messages exchanged between participants. (at development stage)
✅ **Authentication via Secret Key**, ensuring only authorized users can access the chat.  
✅ **No message storage**, the server only relays messages without keeping logs.  
✅ **Open-source** and **self-hosted**, allowing users to run their own chat servers.  
✅ **Dynamic server and port configuration** via environment variables, making the application flexible.  
✅ **Runnable via Docker**, facilitating quick deployment in any environment.  
✅ **Supports multiple connections**, enabling efficient group conversations.  
✅ **Customizable username**, for a more personalized experience.  
✅ **Interactive mode**, allowing real-time message sending and receiving.

---

## 🛠️ How to Run Go-Chatty

### 📦 Using Docker

To run the server:

```bash
docker run -d --name Go-Chatty-server -e SECRET_KEY="my-key" -p 8080:8080 Go-Chatty server
```

To run the client:

```bash
docker run -it --rm Go-Chatty client --server-ip "SERVER_IP" --server-port 8080
```

### 🏗️ Running Manually (Go)

**1. Clone the repository:**

```bash
git clone https://github.com/your-user/go-chatty.git && cd go-chatty
```

**2. Build the code:**

```bash
go build -o go-chatty
```

**3. Start the server:**

```bash
./Go-Chatty server --port 8080 --secret "my-key"
```

**4. Connect a client:**

```bash
./Go-Chatty client --server-ip "SERVER_IP" --server-port 8080
```

---

## 🔒 Security & Privacy

Go-Chatty stands out by prioritizing user privacy. Here’s how:

🔐 **AES-256 Encryption:** All messages sent are encrypted using a key derived from the Secret Key, preventing any intermediary from intercepting and decrypting messages.  
🔑 **Authentication via Secret Key:** Only users with the correct key can enter the chat, increasing control over participation.  
🚫 **No Storage:** Messages are not logged on the server, ensuring nothing is permanently saved.  
🌐 **Self-Hosting:** Users can run their own servers, eliminating reliance on external providers and ensuring full control over infrastructure.  
🔄 **Key Rotation:** The system allows periodic security key changes to prevent unauthorized access.  
🕵 **Metadata Concealment:** No additional message data is stored, ensuring maximum privacy.  
⚠ **Secure Mode:** It is recommended to run the chat within a private network or over a VPN for enhanced security.

---

## 💡 Benefits of Go-Chatty Over Other Chats

📡 **No Centralized Server**: Unlike WhatsApp, Telegram, and Messenger, where messages pass through servers that may store logs, Go-Chatty operates independently and retains no data.  
🎭 **Complete Anonymity**: No personal information is required to join the conversation, unlike conventional services that require emails, phone numbers, or registrations.  
🛠 **Customizable**: Users can run the chat on their own machines or servers, adjusting configurations to meet specific security requirements.  
🌍 **Platform Independent**: Can be run on any system that supports Go or Docker.  
🎚 **Full Control**: Allows dynamic configuration adjustments, such as session expiration time, encryption methods, and user management.

---

## 📜 License

This project is **open-source** and distributed under the MIT license, allowing free use, modification, and distribution.

---

## 📬 Contributions

Contributions are always welcome! If you want to suggest improvements, report bugs, or add new features, feel free to open **issues** or submit **pull requests**. 😃

For discussions and feedback, contact us via [GitHub Issues](https://github.com/mraramlho/Go-Chatty/issues).

---

## 🚀 Future Roadmap

🔹 Implementation of separate chat rooms within the same server.  
🔹 Support for offline messages with an automatic reconnection buffer.  
🔹 Addition of WebRTC support for voice chats.  
🔹 Enhancements to the CLI interface for improved usability.

Stay tuned for updates! 🔥
