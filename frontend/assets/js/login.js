document.getElementById('loginForm').addEventListener('submit', async function (e) {
    e.preventDefault(); // Prevent default form submission

    const formData = new FormData(this);
    const response = await fetch('/auth/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',  // Ensure Content-Type is set to JSON
        },
        body: JSON.stringify({
            email: formData.get('email'),
            password: formData.get('password'),
        }),
    });

    const data = await response.json();
    if (response.ok) {
        alert('Login successful!');
        // Redirect to the protected admin page
        window.location.href = '/admin/panel';
    } else {
        alert('Login failed: ' + data.error);
    }
});
