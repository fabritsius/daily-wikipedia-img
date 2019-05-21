const mainContent = document.querySelector('main');
const reloadDiv = document.querySelector('#reload-indicator');
const reloadTextDiv = document.querySelector('#reload-indicator-text');

const reloadText = 'reload...'
const dYtrigger = 80;

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

const pageMinOffset = 20;

const touchStartHandler = (event) => {
    if (window.pageYOffset <= pageMinOffset) {
        if (event.touches.length === 1) {
            startY = event.touches[0].clientY;
            window.addEventListener('touchmove', moveEventHandler);
        }
    }
}

const touchEndHandler = (event) => {
    mainContent.style.opacity = 1;
    window.removeEventListener('touchmove', moveEventHandler);

    if (window.pageYOffset <= pageMinOffset) {
        const dY = event.changedTouches[0].clientY - startY;
        if (dY > dYtrigger) {
            reloadDiv.className = 'successful';
            setTimeout(() => {
                history.go(0);
            }, 100);
        }
    }

}

window.addEventListener('touchstart', touchStartHandler);
window.addEventListener('touchend', touchEndHandler);