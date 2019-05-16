const APP_NAME = 'daily-wikipedia-img';
const CACHE_NAME = `${APP_NAME}-v1`;
const urlsToCache = [
	'/',
	'styles.css',
	'favicon.ico',
	'sw.js'
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
		caches.keys().then((cacheNames) => {
			return Promise.all(
				cacheNames.filter((cacheName) => {
					return cacheName.startsWith(APP_NAME) &&
						CACHE_NAME != cacheName;
				}).map((cacheName) => {
					return caches.delete(cacheName);
				})
			);
		})
	);
});

// Skip waiting handler
self.addEventListener('message', (event) => {
	if (event.data.action === 'skipWaiting') {
		self.skipWaiting();
	}
});