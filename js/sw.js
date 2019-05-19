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
		fromCache(event.request)
	);
	event.waitUntil(
		update(event.request).then((isUpdated) => {
			if (isUpdated) reload();
		})
	)
});

// Get data from the cache or fetch otherwise
const fromCache = (request) => {
	return caches.open(CACHE_NAME).then((cache) => {
		return cache.match(request).then((response) => {
			return response || fetch(request);
		});
	})
}

// Fetch data and update the cache
const update = (request) => {
	return fetch(request).then((response) => {
		return caches.open(CACHE_NAME).then((cache) => {
			let isUpdated = false;
			// check if new index page is different from what is in the cache
			return cache.match('/').then((indexCache) => {
				if (response.url === indexCache.url) {
					return Promise.all([
						response.clone().text(),
						indexCache.text()
					]);
				}
				return Promise.resolve(false);
			}).then((indexVersions) => {
				if (indexVersions) {
					const [newIndex, oldIndex] = indexVersions;
					if (newIndex !== oldIndex) {
						isUpdated = true;
					}
				}
				return cache.put(request, response);				
			}).then(() => {
				return Promise.resolve(isUpdated);
			});
		});
	});
}


// Sent a message to reload the window
const reload = () => {
	self.clients.matchAll().then((clients) => {
		clients.forEach((client) => client.postMessage('reload-window'));
    });
}

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