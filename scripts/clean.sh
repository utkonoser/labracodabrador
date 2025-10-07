#!/bin/bash
set -e

echo "==================================="
echo "Cleaning Labracodabrador Data"
echo "==================================="

echo ""
echo "Stopping any running instances..."
pkill -f labracodabrador || true
echo "✓ Processes stopped"

echo ""
echo "Removing data directories..."
rm -rf data/
rm -rf logs/
echo "✓ Data directories removed"

echo ""
echo "Removing build artifacts..."
rm -rf bin/
echo "✓ Build artifacts removed"

echo ""
echo "==================================="
echo "Cleanup complete!"
echo "==================================="
echo ""
echo "To rebuild and restart:"
echo "  make build && make run"
echo ""

