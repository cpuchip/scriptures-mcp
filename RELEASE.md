# Release Process

This document explains how to create releases for the scriptures-mcp project.

## Creating a New Release

Releases are triggered by pushing version tags to the repository. The GitHub Actions workflow will automatically build binaries for all supported platforms and create a GitHub release with all the assets.

### Steps:

1. **Ensure your code is ready for release**:
   ```bash
   # Make sure all tests pass
   go test ./...
   
   # Verify the build works
   go build -o scriptures-mcp .
   ```

2. **Create and push a version tag**:
   ```bash
   # Tag the current commit (use semantic versioning)
   git tag v1.0.0
   
   # Push the tag to GitHub
   git push origin v1.0.0
   ```

3. **Monitor the release process**:
   - Go to the Actions tab on GitHub to monitor the workflow
   - The workflow will:
     - Run tests
     - Build binaries for all platforms
     - Create a GitHub release
     - Upload all binaries as release assets

### Supported Platforms

The workflow builds binaries for:
- **Linux**: amd64, arm64, 386, arm
- **macOS**: amd64, arm64 
- **Windows**: amd64, arm64, 386

### Release Artifacts

Each release includes the following binaries:
- `scriptures-mcp-linux-amd64`
- `scriptures-mcp-linux-arm64`
- `scriptures-mcp-linux-386`
- `scriptures-mcp-linux-arm`
- `scriptures-mcp-darwin-amd64`
- `scriptures-mcp-darwin-arm64`
- `scriptures-mcp-windows-amd64.exe`
- `scriptures-mcp-windows-arm64.exe`
- `scriptures-mcp-windows-386.exe`

### Version Tagging Best Practices

Use semantic versioning (MAJOR.MINOR.PATCH):
- **MAJOR**: Incompatible API changes
- **MINOR**: New functionality that's backward compatible
- **PATCH**: Bug fixes that are backward compatible

Examples:
- `v1.0.0` - First stable release
- `v1.1.0` - New feature added
- `v1.1.1` - Bug fix
- `v2.0.0` - Breaking changes

### Troubleshooting

If a release fails:
1. Check the Actions tab for error details
2. Fix any issues in your code
3. Delete the tag if needed: `git tag -d v1.0.0 && git push origin :refs/tags/v1.0.0`
4. Create a new tag with the fix