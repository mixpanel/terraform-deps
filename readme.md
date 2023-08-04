# Terraform Dependencies Parser

This Go program/library is designed to parse Terraform files and build a dependency tree.  It recursively parses all files in a directory and its subdirectories, and returns a map of modules and what (local) modules they import, then dumps that as a mermaid graph.