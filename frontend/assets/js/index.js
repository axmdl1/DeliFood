let currentPage = 1;
const limit = 3;

async function fetchMenu(page) {
    const response = await fetch(`/menu?page=${page}&limit=${limit}`);
    if (response.status === 404) {
        alert("No more items.");
        return;
    }
    const data = await response.json();
    renderMenu(data);
}

function renderMenu(items) {
    const menuContainer = document.getElementById("menu-container");
    menuContainer.innerHTML = ""; // Clear existing content

    items.forEach(item => {
        const box = document.createElement("div");
        box.className = "box";
        box.innerHTML = `
            <div class="image">
                <img src="${item.image}" alt="">
                <h3> <i class="fas fa-utensils"></i> ${item.name} </h3>
            </div>
            <div class="content">
                <div class="price">${item.price}</div>
                <p>${item.desc}</p>
                <a href="#" class="btn">Order now</a>
            </div>
        `;
        menuContainer.appendChild(box);
    });
}

function nextPage() {
    currentPage++;
    fetchMenu(currentPage);
}

function prevPage() {
    if (currentPage > 1) {
        currentPage--;
        fetchMenu(currentPage);
    }
}

// Fetch the first page on load
fetchMenu(currentPage);
