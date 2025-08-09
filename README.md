# CV Hash Parser and Email Service

This project is a Go-based web service that receives a URL to a Curriculum Vitae (CV), calculates its SHA256 hash, generates a report, and sends it to a configured email address with the project's source code attached.

## Features

-   **HTTP Endpoint:** Exposes a simple `POST` endpoint to accept CV URLs.
-   **Hashing:** Calculates the SHA256 hash of the provided URL.
-   **Report Generation:** Creates a JSON report containing the URL, hash, a unique user ID, and a timestamp.
-   **Email Notification:** Sends the report as an email using SMTP.
-   **Source Code Attachment:** Attaches a zip file of the application's source code to the email.
-   **Containerized:** Fully configured to run with Docker and Docker Compose for easy setup and deployment.

## Getting Started

### Prerequisites

-   [Go](https://golang.org/doc/install) (for local development)
-   [Docker](https://docs.docker.com/get-docker/)
-   [Docker Compose](https://docs.docker.com/compose/install/)

### Installation

1.  Clone the repository:
    ```sh
    git clone <repository-url>
    cd <repository-directory>
    ```

2.  Create a `.env` file from the example:
    ```sh
    cp .env.example .env
    ```

3.  Update the `.env` file with your specific configuration, especially your SMTP server details. See the [Configuration](#configuration) section for more details.

## How to Run

### Using Docker (Recommended)

This is the easiest way to get the service running.

1.  Make sure you have completed the [Installation](#installation) steps.
2.  Build and run the service in detached mode:
    ```sh
    docker-compose up --build -d
    ```
3.  The service will be available at `http://localhost:8000`.

### Running Locally

If you prefer to run the service without Docker:

1.  Make sure you have completed the [Installation](#installation) steps.
2.  Install Go dependencies:
    ```sh
    go mod tidy
    ```
3.  The service expects a `source_code.zip` file to be present in the root directory for email attachments. Create it by running:
    ```sh
    zip -r source_code.zip . -x ".git/*" ".env" "source_code.zip"
    ```
4.  Run the application:
    ```sh
    go run cmd/main.go
    ```
5.  The service will be available at `http://localhost:8000`.

## Configuration

The service is configured using environment variables. These variables should be placed in a `.env` file in the root of the project.

| Variable         | Description                                     | Default                |
| ---------------- | ----------------------------------------------- | ---------------------- |
| `SERVER_PORT`    | The port on which the server will listen.       | `8000`                 |
| `SERVER_HOST`    | The host on which the server will run.          | `0.0.0.0`              |
| `SMTP_LOGIN`     | The username for your SMTP server.              | `example@gmail.com`    |
| `SMTP_HOST`      | The hostname of your SMTP server.               | `smtp.gmail.com`       |
| `SMTP_PORT`      | The port for your SMTP server.                  | `587`                  |
| `SMTP_PASSWORD`  | The password for your SMTP server.              | `yourpassword`         |
| `SMTP_RECIPIENT` | The email address that will receive the reports. | `test@gmail.com`       |

## API Usage

### Parse Hash and Send Email

-   **Endpoint:** `POST /parse-hash`
-   **Description:** Accepts a CV URL, processes it, and triggers an email with the report.
-   **Request Body:**

    ```json
    {
      "cv_url": "https://example.com/path/to/your/cv.pdf"
    }
    ```

-   **Success Response:**
    -   **Code:** `200 OK`
    -   **Body:** `Hash parsed and email sent successfully`

-   **Example `curl` Request:**

    ```sh
    curl -X POST http://localhost:8000/parse-hash \
    -H "Content-Type: application/json" \
    -d '{"cv_url": "https://example.com/path/to/your/cv.pdf"}'
    ```
