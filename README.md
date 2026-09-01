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