# Technical Design: Single-Line Version Extraction

## Implementation
In `.github/workflows/release.yml`:
```yaml
VERSION=$(grep -E '^\s*Version\s*=' cmd/root.go | head -n 1 | sed -E 's/.*"([^"]+)".*/\1/')
```
This guarantees a clean single-line semantic version string (`1.4.0`) that embeds properly into Go ldflags without syntax errors.
