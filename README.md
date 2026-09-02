Coding along the course:
```sh
https://www.youtube.com/watch?v=jBf7of9JTV8
```
# Installation & setup
```sh
brew tap mongodb/brew
brew trust mongodb/brew
brew install mongodb-community
brew services start mongodb-community
brew install --cask mongodb-compass
```
Verify
```sh
brew services list 
mongosh --version
```
Open `Mongodb Compass`, enter connection string: `mongodb://127.0.0.1:27017`

In the place of our `go.mod` i.e. `MagicStreamMoviesServer` folder:
- Install Go web framework: `go get -u github.com/gin-gonic/gin`
- Install MongoDB Go driver: `go get go.mongodb.org/mongo-driver/v2/mongo`
- Install `dotenv`: `go get github.com/joho/godotenv`
- Install `Go Playground` (to validate Model - like `validate:required` tag): `go get github.com/go-playground/validator/v10`
- Install `bcrypt` to hass password: `go get golang.org/x/crypto/bcrypt`
- Install JWT library: `go get github.com/golang-jwt/jwt/v5`
- Install LangChain-Go: `go get github.com/tmc/langchaingo/llms/openai`
- Install CORS: `go get github.com/gin-contrib/cors`

## Install Frontend: 
Initialization with `vite`: `npm create vite@latest`\
Use `Bootstrap` for styling: `npm i react-bootstrap bootstrap`\
Include this line in `main.jsx`:  `import 'bootstrap/dist/css/bootstrap.min.css';`\
Install `axios`: `npm i axios`\
Install `react-router-dom`: `npm i react-router-dom`\
Install `react-player` for streaming: `npm i react-player@2.16.0`\
Install for icons: 
```sh
npm i @fortawesome/free-solid-svg-icons
npm i @fortawesome/react-fontawesome
```

# DevOps

## For MongoDB: `MongoDB Atlas`
Install tools for MongoDB devops: `brew install mongodb-database-tools`\
Check installation of `mongodump` with `mongodump --version`\
Open `Compass` and check local URL of our local DB: `localhost:27017`\
Command line typing: `mongodump --uri="mongodb://localhost:27017/magic-stream-movies" --out=./dump` to create a copy of our local DB.\
To push `dump` to cluster: `mongodbrestore --uri="<MONGODB_URI>"` 

## For Client: `Vercel`
Login with Github credentials.\
Choose `Client/magic-stream-client` as `root`.\
Type Build Command: `npm run build` => to `dist/`\

Set Environment Variable:
```sh
VITE_API_BASE_URL: https://build-a-fullstack-movie-streaming-app-in.onrender.com/
```

## For backend: `Render`
Go to Render website.\
Sign up/ Sign in an account with Github.\
Setup for deploying our Go backend:
- choose branch `main`
- type Root Directory `Server/MagicStreamMoviesServer`
- type Build Command: `go build -o app`
- type Start Command: `./app`
- add Environment variables
  - `DATABASE_NAME`: `magic-stream-movies`
  - `MONGODB_URI`: `<atlas-credentials.env>`
  - Other .env variables: `SECRET_KEY`, `SECRET_REFRESH_KEY`, `BASE_PROMPT_TEMPLATE`, `OPENAI_API_KEY`, `RECOMMENDED_MOVIE_LIMIT`
  - `ALLOWED_ORIGINS` as in Vercel platform: `https://build-a-fullstack-movie-streaming-a.vercel.app`
  
After deploy, copy and save the public URL: `https://build-a-fullstack-movie-streaming-app-in.onrender.com/`
  
# Done
Now our app is available at: `https://build-a-fullstack-movie-streaming-a.vercel.app`
