# Go Code Styleguide

1. **Formatting**: Always format code using `gofmt` and `goimports`.
2. **Testing**: Maintain >80% test coverage for all new Go packages (`go test -coverprofile=coverage.out ./...`).
3. **Linting**: Ensure code passes `go vet ./...` without warnings or errors.
4. **Documentation**: Document all exported functions, types, constants, and packages.
5. **Errors**: Handle errors explicitly; do not swallow errors silently.
