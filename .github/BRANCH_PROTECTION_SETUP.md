# Branch Protection Setup Instructions

To complete the CI/CD setup and prevent merging PRs with failing tests, please follow these steps to configure branch protection for the `main` and `development` branches:

## Steps to Configure Branch Protection

1. Go to your repository on GitHub: https://github.com/BohdanStarunskyi/plantie
2. Navigate to **Settings** → **Branches**
3. Click **Add rule** or **Add branch protection rule**

### For the `main` branch:
1. Set **Branch name pattern** to: `main`
2. Check the following options:
   - ✅ **Require a pull request before merging**
   - ✅ **Require status checks to pass before merging**
   - ✅ **Require branches to be up to date before merging**
   - ✅ **Require conversation resolution before merging**
   - In the **Status checks** section, add: `test` (this is the job name from our CI workflow)
3. Optionally enable:
   - ✅ **Restrict pushes that create files larger than 100MB**
   - ✅ **Require linear history** (if you prefer a clean history)
4. Click **Create**

### For the `development` branch:
1. Repeat the same process but set **Branch name pattern** to: `development`
2. Use the same settings as for the `main` branch

## What this accomplishes:

- ✅ PRs to `main` and `development` must pass all tests before merging
- ✅ Direct pushes to protected branches are blocked (except for administrators)
- ✅ All discussions must be resolved before merging
- ✅ Branches must be up to date with the base branch before merging
- ✅ The CI/CD pipeline must complete successfully (the `test` job must pass)

## Testing the setup:

After configuring branch protection, you can test it by:
1. Creating a PR that introduces failing tests
2. Verifying that the PR cannot be merged until tests pass
3. Fixing the tests and confirming the PR can then be merged

The CI workflow will automatically run on every PR to these protected branches and prevent merging if tests fail.