let navbar = document.querySelector('.navbar')

document.querySelector('#menu-btn').onclick = () =>{
    navbar.classList.toggle('active');
    loginForm.classList.remove('active');
    searchForm.classList.remove('active');
}
 

let searchForm = document.querySelector('.search-form')

document.querySelector('#search-btn').onclick = () =>{
    searchForm.classList.toggle('active');
    navbar.classList.remove('active');
    loginForm.classList.remove('active');
}

window.onscroll = () =>{
    navbar.classList.remove('active');
    loginForm.classList.remove('active');
    searchForm.classList.remove('active');
}
 
var swiper = new Swiper(".review-slider", {
    loop:true,
    spaceBetween: 30,
    centeredSlides: true,
    autoplay: {
        delay: 5500,
        disableOnInteraction: false,
    },
    pagination: {
        el: ".swiper-pagination",
    },
});

const packagesContainer = document.getElementById('packages-container');
const prevBtn = document.getElementById('prev-btn');
const nextBtn = document.getElementById('next-btn');
const pageIndicator = document.getElementById('page-indicator');

let currentPage = 1;

const fetchPackages = async (page) => {
    const response = await fetch(`/api/packages?page=${page}&limit=3`);
    const data = await response.json();

    packagesContainer.innerHTML = data.packages.map(pkg => `
            <div class="box">
                <div class="image">
                    <img src="${pkg.image}" alt="${pkg.name}">
                    <h3><i class="fas fa-utensils"></i> ${pkg.name} </h3>
                </div>
                <div class="content">
                    <div class="price">${pkg.price} <span>${pkg.oldPrice}</span></div>
                    <p>${pkg.description}</p>
                    <a href="#" class="btn"> Order now</a>
                </div>
            </div>
        `).join('');

    pageIndicator.textContent = `Page ${page}`;
};

prevBtn.addEventListener('click', () => {
    if (currentPage > 1) {
        currentPage--;
        fetchPackages(currentPage);
    }
});

nextBtn.addEventListener('click', () => {
    currentPage++;
    fetchPackages(currentPage);
});

// Initial load
fetchPackages(currentPage);