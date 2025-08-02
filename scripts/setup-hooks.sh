#!/bin/bash
# Setup git hooks

echo "Setting up git hooks..."

# Set git hooks path
git config core.hooksPath .githooks

echo "✅ Git hooks configured to use .githooks directory"

# Check if pre-commit is available
if command -v pre-commit > /dev/null; then
    # Install pre-commit hooks
    pre-commit install
    pre-commit install --hook-type pre-push
    echo "✅ Pre-commit hooks installed"
else
    echo "⚠️  pre-commit not found. Git hooks are configured but pre-commit framework not available."
    echo "   The .githooks directory hooks will still work!"
fi
echo ""
echo "Hooks installed:"
echo "- pre-commit: Runs on every commit"
echo "- pre-push: Runs before pushing to remote"
echo ""
echo "To run hooks manually:"
echo "  pre-commit run --all-files"