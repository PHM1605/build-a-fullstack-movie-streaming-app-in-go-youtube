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

