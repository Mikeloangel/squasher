# URL Shortener Service

This project is a simple URL shortener service written in Go. It provides an HTTP API to shorten URLs and retrieve the original URLs using the shortened versions.

## Table of Contents

- [Description](#description)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [API Endpoints](#api-endpoints)
- [Contributing](#contributing)
- [License](#license)

## Description

The URL Shortener Service allows users to shorten long URLs and retrieve the original URLs using the shortened versions. It uses an in-memory storage to manage the shortened URLs and their corresponding original URLs.

## Installation

To install and run the URL Shortener Service, follow these steps:

1. Clone the repository:
    ```sh
    git clone https://github.com/Mikeloangel/squasher.git
    cd squasher
    ```

2. Build the project:
    ```sh
    go build -o url-shortener cmd/shortener/main.go
    ```

3. Run the executable:
    ```sh
    ./url-shortener
    ```

## Usage

The URL Shortener Service can be configured using command-line flags or environment variables.

### Command-Line Flags

- `-b`: API host location to get redirect from (default: `http://localhost:8080`)
- `-a`: Sets server location and port in format `host:port` (default: `localhost:8080`)

Example:
```sh
./url-shortener -b http://example.com -a 0.0.0.0:8080
```

### Enviromental variables

- BASE_URL: API host location to get redirect from (default: http://localhost:8080)
- SERVER_ADDRESS: Server location and port in format host:port (default: localhost:8080)

Example:
```sh
export BASE_URL=http://example.com
export SERVER_ADDRESS=0.0.0.0:8080
./url-shortener
```

## Configuration
The configuration is managed using the config package. The State struct holds the application state, including configuration and storage.

## API Endpoints

### Create Short URL
- Endpoint: POST /
- Description: Creates a shortened URL.
- Request Body: The original URL to be shortened.
- Response: The shortened URL.

Example:
```sh
curl -X POST -d "https://www.example.com" http://localhost:8080/
```

### Get Original URL
- Endpoint: GET /{id}
- Description: Retrieves the original URL for the given shortened version.
- Response: Redirects to the original URL.

Example:
```sh
curl -L http://localhost:8080/{id}
```

## Contributing
Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License
This project is licensed under the MIT License. See the LICENSE file for details.