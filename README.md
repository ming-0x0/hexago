# Hexago Project

A modular Go application demonstrating domain-driven design patterns and clean architecture principles.

## Features
- Customer domain management
- Service type management
- Shared utilities (email, transactions, repositories)
- Mock implementations for testing

## Installation
```bash
git clone https://github.com/ming-0x0/hexago.git
cd hexago
go mod download
```

## Building
```bash
go build ./...
```

## Running Tests
```bash
go test ./...
```

## Project Structure
```
internal/
├── customer/               # Customer domain
│   ├── adapter/            # Adapters for external systems
│   ├── domain/             # Core domain logic
│   │   ├── customer/       # Customer entity
│   │   └── service_type/   # Service type entity
│   └── port/               # Port interfaces
└── shared/                 # Shared utilities
    ├── dbmocker/           # Database mocking utilities
    ├── domain/             # Shared domain objects
    │   └── email/          # Email value object
    ├── errors/             # Custom error handling
    ├── repository/         # Repository pattern
    └── transaction/        # Transaction management
```

## Contributing
1. Fork the repository
2. Create a new feature branch
3. Commit your changes
4. Push to your fork
5. Submit a pull request

## License
[MIT](https://choosealicense.com/licenses/mit/)
