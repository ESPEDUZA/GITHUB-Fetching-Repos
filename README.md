# <?xml version="1.0" ?><svg height="1024" width="768" xmlns="http://www.w3.org/2000/svg"><path d="M128 64C57.344 64 0 121.34400000000005 0 192c0 47.219 25.906 88.062 64 110.281V721.75C25.906 743.938 0 784.75 0 832c0 70.625 57.344 128 128 128s128-57.375 128-128c0-47.25-25.844-88.062-64-110.25V302.28099999999995c38.156-22.219 64-63.062 64-110.281C256 121.34400000000005 198.656 64 128 64zM128 896c-35.312 0-64-28.625-64-64 0-35.312 28.688-64 64-64 35.406 0 64 28.688 64 64C192 867.375 163.406 896 128 896zM128 256c-35.312 0-64-28.594-64-64s28.688-64 64-64c35.406 0 64 28.594 64 64S163.406 256 128 256zM704 721.75V320c0-192.5-192-192-192-192h-64V0L256 192l192 192V256c0 0 26.688 0 64 0 56.438 0 64 64 64 64v401.75c-38.125 22.188-64 62.938-64 110.25 0 70.625 57.375 128 128 128s128-57.375 128-128C768 784.75 742.125 743.938 704 721.75zM640 896c-35.312 0-64-28.625-64-64 0-35.312 28.688-64 64-64 35.375 0 64 28.688 64 64C704 867.375 675.375 896 640 896z"/></svg> GITHUB-Fetching-Repos

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
