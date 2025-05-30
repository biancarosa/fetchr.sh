#!/bin/bash

# Demo script for fetchr.sh changelog and release workflow
# This script demonstrates the conventional commit and release process

set -e

echo "🚀 fetchr.sh Release Demo"
echo "=========================="
echo ""

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "❌ Error: Not in a git repository"
    exit 1
fi

# Function to ask for user confirmation
confirm() {
    echo -n "$1 (y/N): "
    read -r response
    case "$response" in
        [yY][eE][sS]|[yY]) 
            return 0
            ;;
        *)
            return 1
            ;;
    esac
}

echo "This demo will:"
echo "1. Install changelog tools"
echo "2. Initialize git-cliff configuration"
echo "3. Make some example conventional commits"
echo "4. Generate a changelog"
echo "5. Create a demo release"
echo ""

if ! confirm "Do you want to continue?"; then
    echo "Demo cancelled."
    exit 0
fi

echo ""
echo "Step 1: Installing changelog tools..."
make install-changelog-tools

echo ""
echo "Step 2: Initializing git-cliff configuration..."
if [ ! -f "cliff.toml" ]; then
    make changelog-init
else
    echo "✅ cliff.toml already exists"
fi

echo ""
echo "Step 3: Let's make some conventional commits..."
echo ""

# Create a demo branch to avoid affecting main
current_branch=$(git branch --show-current)
demo_branch="demo/changelog-demo-$(date +%s)"

echo "Creating demo branch: $demo_branch"
git checkout -b "$demo_branch"

# Create some demo files and commits
echo "# Demo Feature" > demo-feature.md
git add demo-feature.md

echo "Making feature commit..."
make commit-feat msg="add demo feature documentation"

echo "# Demo API" > demo-api.md
git add demo-api.md

echo "Making API commit with scope..."
make commit-feat msg="implement demo API endpoints" scope="api"

echo "Updated demo feature" >> demo-feature.md
git add demo-feature.md

echo "Making fix commit..."
make commit-fix msg="resolve demo feature typo" scope="docs"

echo "# Performance improvements" > performance.md
git add performance.md

echo "Making performance commit..."
make commit-perf msg="optimize demo feature performance"

echo ""
echo "Step 4: Checking conventional commits..."
make check-conventional-commits

echo ""
echo "Step 5: Generating changelog..."
make changelog

echo ""
echo "📄 Generated changelog:"
echo "========================"
cat CHANGELOG.md | head -30
echo "... (truncated)"
echo ""

echo "Step 6: Creating a demo release..."
echo ""
echo "Available release types:"
echo "- patch:  Bug fixes (1.0.0 -> 1.0.1)"
echo "- minor:  New features (1.0.0 -> 1.1.0)"
echo "- major:  Breaking changes (1.0.0 -> 2.0.0)"
echo ""

if confirm "Create a patch release?"; then
    echo "Creating patch release..."
    make release-patch
    
    echo ""
    echo "✅ Release created successfully!"
    echo ""
    echo "To complete the release:"
    echo "1. Review the changes: git show"
    echo "2. Push the release: git push origin main --tags"
    echo "3. Create a GitHub/GitLab release from the tag"
fi

echo ""
echo "Step 7: Cleanup..."
echo "Switching back to $current_branch and cleaning up demo branch..."

git checkout "$current_branch"

if confirm "Delete the demo branch and files?"; then
    git branch -D "$demo_branch"
    git rm -f demo-feature.md demo-api.md performance.md
    git commit -m "chore: cleanup demo files"
    echo "✅ Cleanup completed"
else
    echo "Demo branch '$demo_branch' preserved for your review"
fi

echo ""
echo "🎉 Demo completed!"
echo ""
echo "Summary of what was demonstrated:"
echo "✅ Conventional commit creation with make commands"
echo "✅ Automatic changelog generation with git-cliff"
echo "✅ Semantic version bumping and release creation"
echo "✅ Commit validation and quality checks"
echo ""
echo "Next steps:"
echo "- Read docs/CHANGELOG_GUIDE.md for detailed information"
echo "- Customize cliff.toml for your project needs"
echo "- Integrate with your CI/CD pipeline"
echo "- Start using conventional commits in your workflow"
echo ""
echo "Happy committing! 🚀" 