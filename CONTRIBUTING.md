## Contribute Workflow

1. Create a Pull Request (PR) with your changes
2. Ensure all tests and linters pass
3. Approve and merge the PR into the main branch
4. Create and push a new version tag:  
```bash 
git tag v2.0  
git push origin v2.0
```
📦 This will automatically trigger the GitHub Actions workflow to build and release the package.