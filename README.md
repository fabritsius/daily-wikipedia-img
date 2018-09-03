# Wikipedia picture of the day feed

### About this project

This project is my first attempt at golang. I tried to include every single concept and almost succeeded. I actually have supporting project which I decided to exclude temporarily before committing the code (I will add it back later). This repository can be helpful to those who also start their journey with Go. Such "useful" projects are still great for learning purposes.

### Usage

To run this server locally:

0. Have Go installed (remember about GOPATH, if it matters on your machine);
1. Clone this repo `git clone https://github.com/fabritsius/daily-wikipedia-img`
2. Change directory `cd daily-wikipedia-img`
3. Enter `go run main.go` to start the server
4. Visit [localhost:5000](https://localhost:5000) with your browser

### TODO

- [x] Add core features and create this repo
- [ ] Add all temporarily excluded features
- [ ] Add another path to the server for JSON requests
- [ ] Upload server to Heroku
- [ ] Create gh-pages wrapper for the project which uses JSON response
- [ ] Host the website with GitHub

### Thoughts about Go

Overall, I like it a lot and agree with every benefit people say when the talk about golang. The only thing that's frustrating is this mandatory "no unused variable" thing. Don't get me wrong, it is great in the production or right before the project is complete. But, for the most part it is a huge pain in somewhere. Simple example is when I comment out one line and now I am told to comment out this function and some unused variable (I will have to return everything back eventually) and this process can continue further. Maybe this behavior doesn't matter in a huge projects, but during learning process the do matter. Also, I heard a lot about lightning speed and how fast is Go. I have to say that after all this torment time saved is probably negative.

My solution is to have this rule only on build process. So, when I use `go run main.go` I see a warning, but my code compiles and when I use `go build main.go` process stops the way it currently does.
