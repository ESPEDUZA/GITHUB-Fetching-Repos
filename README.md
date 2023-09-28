# 🚀 CC-GO

## 📝 Description
This project allows you to clone all GitHub repositories of a user or an organization and archive them into a zip file. The archive can then be downloaded via an HTTP API. If a GitHub token is provided, the application will also clone private repositories. The repositories are fetched, and the latest active branch is pulled.

## 🛠 Prerequisites
- [Docker](https://www.docker.com/get-started)
- [GitHub Token](https://github.com/settings/tokens) (optional, required for cloning private repositories)

## 🚀 Setup & Run

### 1. **Clone the Repository**
```sh
git clone https://github.com/ESPEDUZA/CC-GO
cd CC-GO
```

### 2. Configure Environment Variables

Copy the .env.dist file to a new file named .env and update the environment variables as needed.
```sh
cp .env.dist .env
```
Edit the .env file and set the values for the GitHub user and token (if available).

### 3. Build and Run with Docker
```sh
docker build -t cc-go .
docker run -p 8080:8080 cc-go
```

## 🌐 Usage

Once the application is running, you can download the archived repositories by navigating to:
```url
http://localhost:8080/download
```
Replace 8080 with the port number you have configured if it's different.

## 🧑‍💻 Author

Eliott GERMAIN
