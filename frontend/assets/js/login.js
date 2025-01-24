document.getElementById("loginForm").addEventListener("submit", function (e) {
    e.preventDefault();

    const email = document.querySelector('input[name="email"]').value;
    const password = document.querySelector('input[name="password"]').value;

    fetch('/auth/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: new URLSearchParams({
            email: email,
            password: password,
        }),
    })
        .then(response => response.json())
        .then(data => {
            if (data.token) {
                localStorage.setItem('token', data.token);  // Save token to localStorage (for example)
                if (data.role === "admin") {
                    window.location.href = "/admin/panel"; // Redirect to admin panel if admin
                } else {
                    window.location.href = "/"; // Redirect to main page if user
                }
            } else {
                alert('Login failed. Please check your credentials.');
            }
        })
        .catch(error => {
            console.error('Error:', error);
        });
});
