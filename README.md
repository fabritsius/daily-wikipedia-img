# Wikipedia's picture of the day feed

### About this project

This project is my first attempt at Golang. I tried to include every single concept (and kind of failed). This repository can be helpful to those who also start their journey with Go.

I also made a [version](https://github.com/fabritsius/daily-wikipedia-img/tree/app-engine-ver) to be hosted with Google's App Engine. You can click [the URI](https://fabritsius.github.io/daily-wikipedia-img/) on the top of the project page or use the link below. I also made a simple [HTML wrapper](https://github.com/fabritsius/daily-wikipedia-img/tree/gh-pages) to hide App Engine's URI (so it looks like the webpage is hosted by GitHub).

This application works offline after first reload and it can be installed on your device. If you choose to do so, please use [this version hosted by App Engine](https://daily-pict.appspot.com/) (without my HTML wrapper). Both URIs work, [this one](https://daily-pict.appspot.com/) just works better and you are bothering only one server this way.

#### [Open the App](https://fabritsius.github.io/daily-wikipedia-img/)

### Usage

To run this server locally:

0. Have Go installed (remember about GOPATH, if it matters on your machine)
1. Clone this repo `git clone https://github.com/fabritsius/daily-wikipedia-img`
2. Change directory `cd daily-wikipedia-img`
3. Enter `go run main.go` to start the server
4. Visit [localhost:5000](https://localhost:5000) with your browser

To deploy it to App Engine:

0. Have gcloud installed and run `gcloud components update`
1. Create a new project `gcloud projects create [YOUR_PROJECT_NAME]`
2. Create a new app `gcloud app create --project=[YOUR_PROJECT_NAME]`
3. Install gcloud components for Go `gcloud components install app-engine-go`
4. Clone this repo `git clone https://github.com/fabritsius/daily-wikipedia-img`
5. Change branch `git checkout app-engine-ver`
6. Change directory `cd daily-wikipedia-img`
7. Deploy `gcloud app deploy`
8. Open the app `gcloud app browse`

### TODO

- [x] Add core features and create this repo
- [x] Upload server to App Engine
- [x] Create gh-pages wrapper for the project
- [x] Host the website with GitHub
- [x] Make the design more pleasing to look at
- [x] Add a simple offline experience
- [x] Make the website installable (convert to PWA)
- [ ] Add more features (like image pop-ups)