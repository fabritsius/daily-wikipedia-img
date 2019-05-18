# Wikipedia's picture of the day feed

### About this version of the project

Code on this branch was altered to be hosted with Google's App Engine. There is a link on the top of this page or below. I also made a simple [HTML wrapper](https://github.com/fabritsius/daily-wikipedia-img/tree/gh-pages) to hide App Engine's URI (so it looks like the webpage is hosted by GitHub).

This application works offline after first reload and it can be installed to your device. If you choose to do so, please use [this version hosted by App Engine](https://daily-wiki-24-7.appspot.com/) (without my HTML wrapper). Both URIs do work, [this one](https://daily-wiki-24-7.appspot.com/) just works better and you are bothering only one server this way.

#### [Open the App](https://fabritsius.github.io/daily-wikipedia-img/)

### Usage

To run this server locally:

0. Have Go with appengine library installed locally;
1. Clone this repo `git clone https://github.com/fabritsius/daily-wikipedia-img`
2. Change directory `cd daily-wikipedia-img`
3. Enter `dev_appserver.py .` to start the server
4. Visit [localhost:8080](https://localhost:8080) with your browser

### TODO

- [x] Add core features and create this repo
- [x] Upload server to App Engine
- [x] Create gh-pages wrapper for the project
- [x] Host the website with GitHub
- [x] Make the design more pleasing to look at
- [x] Add a simple offline experience
- [x] Make the website installable (convert to PWA)
- [ ] Add more features (like image pop-ups)