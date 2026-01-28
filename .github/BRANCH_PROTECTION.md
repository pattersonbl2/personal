# Branch Protection Rules

## Policy: Never Merge Directly to Main

This repository enforces a **no direct merge to main** policy. All changes must go through Pull Requests.

### Workflow

1. **Create a feature branch** from `main`:
   ```bash
   git checkout main
   git pull origin main
   git checkout -b your-feature-branch
   ```

2. **Make your changes** and commit:
   ```bash
   git add .
   git commit -m "Your commit message"
   ```

3. **Push your branch**:
   ```bash
   git push -u origin your-feature-branch
   ```

4. **Create a Pull Request** on GitHub

5. **Review and merge** the PR (this will merge into main)

### Why?

- **Code Review**: Ensures all changes are reviewed before merging
- **CI/CD Safety**: Allows automated checks to run before merging
- **History**: Better git history and easier rollbacks
- **Collaboration**: Makes it easier to discuss changes before they're merged

### GitHub Branch Protection Setup

**IMPORTANT**: Configure branch protection in GitHub repository settings:

1. Go to: **Settings → Branches → Branch protection rules**
2. Click **Add rule**
3. Branch name pattern: `main`
4. Enable the following:
   - ✅ **Require a pull request before merging**
   - ✅ **Require approvals** (recommended: 1 approval)
   - ✅ **Dismiss stale pull request approvals when new commits are pushed**
   - ✅ **Require status checks to pass before merging** (optional, if you have CI/CD)
   - ✅ **Require branches to be up to date before merging** (if using status checks)
   - ✅ **Do not allow bypassing the above settings** (prevents admin overrides)

### Local Safeguard

A pre-push git hook is included (`.git/hooks/pre-push`) to prevent accidental direct pushes to main locally. This is a local safeguard only - GitHub branch protection is the primary enforcement mechanism.
