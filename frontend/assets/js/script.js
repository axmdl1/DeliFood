// Get references to the elements
let navbar = document.querySelector('.navbar');
let searchForm = document.querySelector('.search-form');

// Menu button to toggle the navbar
document.querySelector('#menu-btn').onclick = () => {
    navbar.classList.toggle('active');   // Toggle navbar visibility
    searchForm.classList.remove('active');  // Ensure other forms are hidden
};

// Search button to toggle the search form
document.querySelector('#search-btn').onclick = () => {
    searchForm.classList.toggle('active');  // Toggle search form visibility
    navbar.classList.remove('active');   // Ensure navbar is hidden when search form is active
};

// Login button to redirect to login page
document.querySelector('#login-btn').onclick = () => {
    window.location.href = '/auth/login';  // Redirect to the login page
};

document.querySelector('#cart-btn').onclick = () => {
    window.location.href = '/cart/items';
}

// Remove active class when scrolling
window.onscroll = () => {
    navbar.classList.remove('active');  // Hide navbar on scroll
    searchForm.classList.remove('active');  // Hide search form on scroll
};

// Initialize Swiper for review slider
var swiper = new Swiper(".review-slider", {
    loop: true,
    spaceBetween: 30,
    centeredSlides: true,
    autoplay: {
        delay: 5500,
        disableOnInteraction: false,
    },
    pagination: {
        el: ".swiper-pagination",
    },
});/*

// Pagination for packages
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
                <a href="#" class="btn">Order now</a>
            </div>
        </div>
    `).join('');

    pageIndicator.textContent = `Page ${page}`;
};

// Handle pagination (previous and next buttons)
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


// Initial load of packages
fetchPackages(currentPage);*/
