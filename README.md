# MockThis CLI

MockThis is a command-line interface (CLI) tool for managing mock API endpoints. It allows users to create, update, delete, and list mock endpoints for testing and development purposes.

## Installation

### macOS

To install MockThis on macOS using Homebrew, follow these steps:

1. First, tap the repository:
   ```
   brew tap nicobistolfi/carbon
   ```

2. Then, install MockThis:
   ```
   brew install mockthis
   ```

## Key Features

- User authentication (login and registration)
- Create new mock endpoints
- Update existing endpoints
- Delete endpoints
- List all created endpoints
- Get details of specific endpoints

## Usage

The general syntax for using MockThis is:

```
mockthis [command] [flags]
```

### Available Commands

- `create`: Create a new mock endpoint
- `update`: Modify an existing endpoint
- `delete`: Remove an endpoint
- `list`: Display all created endpoints
- `get`: Retrieve details of a specific endpoint
- `login`: Authenticate user
- `register`: Create a new user account
- `completion`: Generate shell autocompletion scripts

For detailed information on each command, use:

```
mockthis [command] --help
```

## Registering a new user

To register a new user account, use the register command. You can provide your email as an argument or enter it when prompted.

```
mockthis register
```

## Logging in

To log in to MockThis, use the login command. You can either provide your email as an argument or enter it when prompted.

```
mockthis login {email}
```
If you don't provide an email, you will be prompted to enter it.


## Creating a new endpoint

To create a new mock endpoint, use the create command. You can provide the endpoint details as arguments or enter them when prompted.

```
mockthis create --method GET --status 200 --body '{"message": "Hello, World!"}'
```
Comand line output will include the endpoint details, such as the method, status, body, and path.

```
Endpoint created successfully!

Endpoint ID: 1234567890
Endpoint Method: GET
Endpoint Status: 200
Endpoint Body: {"message": "Hello, World!"}

Endpoint URL: https://api.mockthis.io/m/abad312-das-3231nmd
```

