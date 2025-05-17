#!/bin/bash
echo "🔍 Checking for print/debug statements..."
grep -rnw './' -e 'fmt.Print' --exclude-dir=vendor
grep -rnw './' -e 'log.Print' --exclude-dir=vendor
grep -rnw './' -e 'println(' --exclude-dir=vendor

echo -e "\n📝 Checking for TODO, FIXME, DEBUG comments..."
grep -rnw './' -e 'TODO' --exclude-dir=vendor
grep -rnw './' -e 'FIXME' --exclude-dir=vendor
grep -rnw './' -e 'DEBUG' --exclude-dir=vendor

echo -e "\n🗃️  Checking for large files (>1MB)..."
find . -type f -size +1M -exec ls -lh {} \; | awk '{ print $NF ": " $5 }'

echo -e "\n🗑️  Checking for untracked ignored files..."
git status --ignored -s

echo -e "\n📁 Checking for unneeded folders (testdata, .vscode, *.out)..."
find . -type d \( -name ".vscode" -o -name "testdata" \)
find . -name "*.out"
