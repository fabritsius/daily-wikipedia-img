const navText = document.querySelector('nav h1');
const content = document.querySelector('div.body-filler');

document.addEventListener('click', (event) => {
    if (event.target === navText) {
        content.scrollTo({
            top: 0,
            behavior: 'smooth'
        });
    }
});