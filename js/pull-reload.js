const mainContent = document.querySelector('main');
const reloadDiv = document.querySelector('#reload-indicator');
const reloadTextDiv = document.querySelector('#reload-indicator-text');

const reloadText = 'reload '
const moons = ['🌕', '🌖', '🌗', '🌘', '🌑', '🌒', '🌓', '🌔']
const dYtrigger = 80;

let startY = 0;
let phasing = null;

const moveEventHandler = (event) => {
    if (event.touches.length === 1) {
        const dY = event.touches[0].clientY - startY;
        if (dY > 0) {
            mainContent.style.opacity = 0.7;
            reloadDiv.style.opacity = dY / dYtrigger;
            animateReloadBtn();
        }
    }
}

const animateReloadBtn = () => {
    if (!phasing) {
        const frames = moons.length - 1;
        let phase = 0
        phasing = setInterval(() => {
            reloadTextDiv.innerHTML = reloadText + moons[phase++];
            if (phase > frames) {
                phase = 0;
            }
        }, 150);
    }
}

const removeReloadBtnAnimation = () => {
    clearInterval(phasing);
    phasing = null;
}

const pageMinOffset = 20;

const touchStartHandler = (event) => {
    if (mainContent.parentElement.scrollTop <= pageMinOffset) {
        if (event.touches.length === 1) {
            startY = event.touches[0].clientY;
            window.addEventListener('touchmove', moveEventHandler);
        }
    }
}

const touchEndHandler = (event) => {
    mainContent.style.opacity = 1;
    window.removeEventListener('touchmove', moveEventHandler);
    removeReloadBtnAnimation();

    if (mainContent.parentElement.scrollTop <= pageMinOffset) {
        reloadDiv.style.opacity = 0.3;
        const dY = event.changedTouches[0].clientY - startY;
        if (dY > dYtrigger) {
            reloadDiv.className = navigator.onLine ? 'successful' : 'failed';
            setTimeout(() => {
                history.go(0);
            }, 100);
        }
    }

}

window.addEventListener('touchstart', touchStartHandler);
window.addEventListener('touchend', touchEndHandler);