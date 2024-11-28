# Cafeteller

Cafeteller is a review blog platform for coffee enthusiasts, powered by a robust API v2 built with Go and Gin. This repository contains the source code for Cafeteller's backend, providing endpoints for managing coffee reviews, and more.

## Features

- Comprehensive coffee reviews
- Review submission and moderation
- Get Reviews
- Get Reviews by ID

## Getting Started

### Prerequisites

- Go 1.18 or later
- Gin web framework
- Cloud Firestore

### Installation

1. Clone the repository:
    ```sh
    git clone hhttps://github.com/realkyr/cafeteller-api-v2.git
    ```
2. Navigate to the project directory:
    ```sh
    cd cafeteller
    ```
3. Install dependencies:
    ```sh
    go mod download
    ```
4. Set up your database and configure the connection in the `.env` file.
   ```dotenv
   GO_ENV=
   PROJECT_ID=
   SERVICE_ACCOUNT_KEY_PATH=
   ```

### Running the Application

1. Run the application:
    ```sh
    go run main.go
    ```

2. The API will be available at `http://localhost:8080`.

## Usage

Use the provided endpoints to interact with the Cafeteller API. Refer to the API documentation for detailed information on available endpoints and their usage.

## License

This repository is not under any specific open-source license, and the code is provided for viewing purposes only. While it is not legally binding, I highly recommend not using this code as a reference for your own projects. This code is intended to be part of my portfolio, showcasing my skills and experience.

## Contact

For any questions or inquiries, please contact me at [phuree.kanusont@gmail.com](mailto:phuree.kanusont@gmail.com).