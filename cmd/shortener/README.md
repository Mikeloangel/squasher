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

The URL Shortener Service allows users to shorten long URLs and retrieve the original URLs using the shortened versions. It can use different storagers to manage the shortened URLs and their corresponding original URLs.

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

### Settings priority
The most priority has Enviromental variables
Next is command-line flags
If none provided uses default settings

### Command-Line Flags

- `-a`: Sets server location and port in format `host:port` (default: `localhost:8080`)
- `-b`: API host location to get redirect from (default: `http://localhost:8080`)
- `-f`: Activates in file db storage for given file (default: `/tmp/short-url-db.json`) 
- `-d`: Activates SQL db storage (on pgs driver) dsn string is expected (default empty string: ``)

Example:
```sh
./url-shortener -b http://example.com -a 0.0.0.0:8080
```

### Enviromental variables

- `BASE_URL`: API host location to get redirect from (default: `http://localhost:8080`)
- `SERVER_ADDRESS`: Server location and port in format host:port (default: `localhost:8080`)
- `FILE_STORAGE_PATH`: Activates in file db storage for given file (default: `/tmp/short-url-db.json`) 
- `DATABASE_DSN`: Activates SQL db storage (on pgs driver) dsn string is expected (default empty string: ``)

### Storagers
Priority to load storager:
if `-d` or `DATABASE_DSN` string is not empty will try to use DB storager
if `-f` or `FILE_STORAGE_PATH` string is not empty will try to use file storager
Otherwise will activate in-memory storager

Example:
```sh
export BASE_URL=http://example.com
export SERVER_ADDRESS=0.0.0.0:8080
./url-shortener
```

## Configuration
The configuration is managed using the config package. The State struct holds the application state, including configuration and storage.

## API Endpoints
All api endpoints support gzip headers

### Create Short URL
- Endpoint: POST /
- Description: Creates a shortened URL.
- Request Body: The original URL to be shortened.
- Response: The shortened URL.
- Response status: 201 if created, 409 if already existed

Example:
```sh
curl -X POST -d "https://www.example.com" http://localhost:8080/
```

### Create Short URL with JSON format
- Endpoint: POST /api/shortens
- Description: Creates a shortened URL.
- Request Body: The original URL to be shortened in json format.
Example:
```json
{
    "url":"http://ya.ru/"
}
```
- Response: The shortened URL in json format.
Example:
```json
{
    "result": "http://localhost:8080/59c5b1a6"
}
```
- Response status: 201 if created, 409 if already existed

Example:
```sh
curl -X POST -d "https://www.example.com" http://localhost:8080/
```

### Create Short in batch
- Endpoint: POST /api/shorten/batch
- Description: Creates a shortened URLs in one batch.
- Request Body: JSON array
Example:
```json
[
    {
        "correlation_id": "{{$randomUUID}}",
        "original_url": "http://ya.ru/hello-go/451"
    }
]
```
- Response: JSON array with shortened URL and correlation_id from request
Example:
```json
[
    {
        "correlation_id": "f909fcef-1bf6-4174-8c53-daa05469ae37",
        "short_url": "http://localhost:8080/33da5a2c"
    }
]
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