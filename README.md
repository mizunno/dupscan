# Duplicate File Finder

This is a simple Go program that finds duplicate files in a given directory. It uses the MD5 hash of the files to determine if they are duplicates. Before applying the MD5 hash, the program checks if the files have the same size, as files with different sizes cannot be duplicates.

- Scan files in the specified directory and identifies potential duplicates based on file size.
- Verify duplicates by comparing MD5 hashes of files.
- Utilize a worker pool of Go routines for improved performance.
- Generates a report with the list of duplicate files.

## Usage

### Flags

- `--dir`: Directory to scan (default: current directory `"."`).
- `--minsize`: Minimum file size to consider (in bytes, default: `0`).
- `--workers`: Number of worker goroutines to process files in parallel (default: the number of available CPU cores).

```bash
go run main.go --dir /path/to/scan --minsize 1024 --workers 4
```

This example scans the directory /path/to/scan, only considering files larger than 1024 bytes, and uses 4 workers.

### Output

The program will generate a report in the format duplicates-YYYYMMDDHHMMSS.txt, where each line contains a comma-separated list of duplicate files.

## Test

You can use the bash script `create_random_files.sh` to generate random files for testing.

```bash
#!/bin/bash

DIR="./files"
NUM_FILES=1000
FILE_SIZE=1024  # Size in bytes

# Create directory if it doesn't exist
mkdir -p "$DIR"

# Create files
for i in $(seq 1 $NUM_FILES); do
  dd if=/dev/urandom of="$DIR/file_$i.dat" bs=$FILE_SIZE count=1 status=none
  echo "Created file_$i.dat"
done
```


Run the script to generate 1000 files with 1024 bytes each.

```bash
./create_random_files.sh
```
