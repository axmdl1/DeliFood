document.getElementById('loginForm').addEventListener('submit', async function (e) {
    e.preventDefault(); // Prevent form submission

    const formData = new FormData(this);
    const response = await fetch('/auth/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            email: formData.get('email'),
            password: formData.get('password'),
        }),
    });

    const data = await response.json();
    if (response.ok) {
        // Store the token in localStorage
        localStorage.setItem('token', data.token);
        alert('Login successful!');
        // Redirect to the appropriate page
        if (data.role === 'admin') {
            window.location.href = '/admin/panel'; // Redirect to admin panel
        } else {
            window.location.href = '/'; // Redirect to main page
        }
    } else {
        alert('Login failed: ' + data.error);
    }
});
