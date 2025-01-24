document.getElementById("loginForm").addEventListener("submit", async (event) => {
    event.preventDefault();

    const email = document.querySelector('input[name="email"]').value;
    const password = document.querySelector('input[name="password"]').value;

    try {
        const response = await fetch("/auth/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/x-www-form-urlencoded",
            },
            body: new URLSearchParams({ email, password }),
        });

        if (response.ok) {
            const data = await response.json();
            alert(data.message);

            // Store token in localStorage
            localStorage.setItem("token", data.token);

            // Redirect to /admin/panel
            window.location.href = "/admin/panel";
        } else {
            const errorData = await response.json();
            alert(`Error: ${errorData.message}`);
        }
    } catch (error) {
        console.error("Login failed:", error);
        alert("An error occurred during login.");
    }
});


document.getElementById("registerForm").addEventListener("submit", async (event) => {
    event.preventDefault();

    const email = document.querySelector('input[name="email"]').value;
    const user = document.querySelector('input[name="user"]').value;
    const password = document.querySelector('input[name="password"]').value;


    try {
        const response = await fetch("/auth/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/x-www-form-urlencoded",
            },
            body: new URLSearchParams({ email, password }),
        });

        if (response.ok) {
            const data = await response.json();
            alert(data.message);

            // Store token in localStorage
            localStorage.setItem("token", data.token);

            // Redirect to /admin/panel
            window.location.href = "/admin/panel";
        } else {
            const errorData = await response.json();
            alert(`Error: ${errorData.message}`);
        }
    } catch (error) {
        console.error("Login failed:", error);
        alert("An error occurred during login.");
    }
});
