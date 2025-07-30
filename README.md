# VampireZ Queue Monitor Bot 🦇

A lightweight Discord bot designed to assist with monitoring and managing **VampireZ** queues, with additional utilities for your server.

Automatically tracks and notifies when a VampireZ lobby is filling up. Customize queue thresholds, ping roles, and channels to get alerts.

## 🤖 Available Commands

- `register` - Registers your Hypixel API key
- `start` - Starts the VampireZ Queue Monitoring (Requires Hypixel API key)
- `stop` - Stops the VampireZ Queue Monitoring (Requires Hypixel API key)

## 🛠️ Installation Guide
Follow these steps to install and run the Discord bot locally.

### ✅ Prerequisites

- Go 1.20 or higher
- A Discord bot token

### 📦 1. Clone the Repository

```
git clone https://github.com/yourusername/your-bot-repo.git
```

```
cd your-bot-repo
```


### ⚙️ 2. Set Up Environment Variables

```
make setup
```

Then, edit the `.env` file to add your actual Discord bot token.

### 🚀 3. Run the Bot

```
make run
```

### 🧹 4. Clean Up

```
make clean
```

