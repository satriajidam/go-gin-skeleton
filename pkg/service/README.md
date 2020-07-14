# service directory

This is the directory where all business logic related codes are stored.

There are 3 types of directory here and they encapsulate their own package:

- `domain`: this is where all the business domain entities and abstractions are placed.
- `client`: this is where all clients for communicating with external dependencies are placed.
- `<entity_name>`: this is where a specific business entity implementation is placed (e.g. `customer`, `provider`, etc.)

## References

- [Uncle Bob's Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [bxcodec/go-clean-arch](https://github.com/bxcodec/go-clean-arch)
