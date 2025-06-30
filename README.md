# 🔐 OPAQUE Client (Go) - Full Featured CLI

This repository contains a fully working **OPAQUE client** (`cmd/client/main.go`) that interacts seamlessly with the matching server (`cmd/server/main.go`). It supports secure **password registration** and **authentication** based on [OPAQUE](https://datatracker.ietf.org/doc/html/draft-irtf-cfrg-opaque).

> ✅ Compatible with RSA-512 configuration used in the server.

---

## 🧠 How It Works

### 🔐 Registration Mode (`-pwreg`)
Registers a new username-password securely using OPRF-based key exchange.

```bash
go run main.go -conn localhost:9999 -pwreg -username=alice -password=1234
```

📶 Flow:
1. `PwRegInit()` → generate `PwRegMsg1`
2. Client → Send `PwRegMsg1` to server
3. Server → Responds with `PwRegMsg2`
4. Client runs `PwReg2()` → produces `PwRegMsg3`
5. Client → Sends `PwRegMsg3`
6. Server → Confirms with `ok`

### 🔐 Authentication Mode (`-auth`)
Performs login + secure key exchange + encrypted communication.

```bash
go run main.go -conn localhost:9999 -auth -username=alice -password=1234
```

📶 Flow:
1. `AuthInit()` → generates `AuthMsg1`
2. Client → Send `AuthMsg1`
3. Server → Responds with `AuthMsg2`
4. Client runs `Auth2()` → gets shared secret and `AuthMsg3`
5. Client → Sends `AuthMsg3`
6. Server → Confirms with `ok`

### 🔐 Encrypted Exchange:
- Server → Sends encrypted: **"Hi client!"**
- Client → Decrypts, replies with: **"Hi server!"**

🔑 Shared secret from OPRF used for symmetric AES encryption/decryption.

---

## 📁 Folder Structure
```
cmd/
├── client/        # Client CLI code (this app)
└── server/        # Server counterpart for testing
```

---

## 💡 Tip
Make sure the server is running:
```bash
go run cmd/server/main.go
```

Then run the client with your desired mode.

---

## 📜 License
BSD-style license © 2018 Fredrik Kuivinen
