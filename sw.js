const APP_NAME = 'daily-wikipedia-img';
const CACHE_NAME = `${APP_NAME}-v1`;
const urlsToCache = [
  `/${APP_NAME}/`,
  `/${APP_NAME}/index.html`,
  `/${APP_NAME}/favicon.ico`	
];

// Perform install steps
self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME).then((cache) => {
      console.log('Cache opened');
      return cache.addAll(urlsToCache);
    })
  );
});

// Respond with cached resources
self.addEventListener('fetch', (event) => {
  event.respondWith(
    caches.match(event.request).then((request) => {
      return request || fetch(event.request).then((response) => {
        return caches.open(CACHE_NAME).then((cache) => {
          cache.put(event.request, response.clone());
          return response;
        });
      });
    })
  )
})

// Delete outdated caches
self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches.keys().then((keyList) => {
      // `keyList` contains all cache names under your username.github.io
      // filter out ones that has this app prefix to create white list
      let cacheWhitelist = keyList.filter((key) => {
        return key.indexOf(APP_NAME)
      })
      // add current cache name to white list
      cacheWhitelist.push(CACHE_NAME);

      return Promise.all(keyList.map((key, i) => {
        if (cacheWhitelist.indexOf(key) === -1) {
          console.log(`deleting cache: ${keyList[i]}`)
          return caches.delete(keyList[i])
        }
      }))
    })
  )
});

// Skip waiting handler
self.addEventListener('message', (event) => {
  if (event.data.action === 'skipWaiting') {
    self.skipWaiting();
  }
});
