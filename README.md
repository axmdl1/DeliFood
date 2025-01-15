# DeliFood Web Application

DeliFood is a web application for managing menu items and facilitating communication between users and administrators. It supports rate-limiting, file uploads via the contact form, and dynamic content rendering.


![image](https://github.com/user-attachments/assets/321c3f89-602b-4dab-aa8c-aeb2fa1ced10)

---

## Features

- **Dynamic Menu Pages**: Display paginated menu items with filtering and sorting options.
- **Contact Form**: Users can send messages and attach files to the admin.
- **Rate Limiting**: Prevents abuse by limiting requests per second on all endpoints.
- **Graceful Shutdown**: Ensures the server shuts down cleanly while handling ongoing requests.
- **Admin Email Notifications**: Admin receives user-submitted emails with optional attachments.
- **Static Assets**: Includes a styled frontend with CSS and images.

---

## Prerequisites

Before running this project, ensure you have the following installed:

- **Go** (1.19 or higher)
- **SMTP Server Credentials** (e.g., Mail.ru, Gmail, etc.)
- **Internet Access** (to serve static assets and send emails)

---

## Setup and Installation

1. **Clone the Repository**:
   git clone https://github.com/axmdl1/DeliFood.git
   cd DeliFood
2. **Install Dependencies**: 
   go mod tidy
3. **Configure SMTP Credentials: Update the handlers.ContactUsHandler function with your SMTP email and password**:
   dialer := gomail.NewDialer("smtp.mail.ru", 587, "your-email@mail.ru", "your-password")
4. **Create Required Directories: Ensure the temp directory exists for file uploads**:
   mkdir temp
5. **Run the application**:
   go run main.go
6. **Access the Application: Open your browser and navigate to http://localhost:9078**


