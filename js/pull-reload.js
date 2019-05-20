const mainContent = document.querySelector('main');
const reloadTextDiv = document.querySelector('#reload-indicator-text');

const reloadText = 'reload...'
const dYtrigger = 60;

let startY = 0;

const moveEventHandler = (event) => {
    if (event.touches.length === 1) {
        const dY = event.touches[0].clientY - startY;
        const minOpacity = 1 - (reloadText.length * 0.1);
        const opacity = between(1 - (dY / 100), minOpacity, 1);
        mainContent.style.opacity = opacity;
        const lettersCount = Math.floor((1 - opacity) * 10);
        reloadTextDiv.innerHTML = reloadText.slice(0, lettersCount);
    }
}

const between = (value, minLimit, maxLimit) => {
    return Math.max(minLimit, Math.min(value, maxLimit)); 
}

const touchStartHandler = (event) => {
    if (window.pageYOffset === 0) {
        if (event.touches.length === 1) {
            startY = event.touches[0].clientY;
            window.addEventListener('touchmove', moveEventHandler);
        }
    }
}

const touchEndHandler = (event) => {
    mainContent.style.opacity = 1;
    reloadTextDiv.innerHTML = '';
    window.removeEventListener('touchmove', moveEventHandler);

    if (window.pageYOffset === 0) {
        const dY = event.changedTouches[0].clientY - startY;
        if (dY > dYtrigger) {
            window.location.reload();
        }
    }

}

window.addEventListener('touchstart', touchStartHandler);
window.addEventListener('touchend', touchEndHandler);