## Terminal App for 1 Billion Row Challenge

# Description

Just building this terminal application to have multiple solutions with different approaches to solve the 1 Billion Row Challenge. The challenge is to read a file with 1 billion rows of city name and it's reported temperature and print the cities with their min/avg/max temperature.

# Features

Create Command: Quickly create files with a specified name and content. This command is optimized to handle large inputs and can be used to generate files with a substantial amount of data.

- Takes Around 4mins to create a file with 1 billion rows of data.

```bash
./out/1brc create --file=largefile.txt --rows=1000000000
```

Simple Process Command: Read and display the content of a specified file. This command is just to process the file in sequential and in very inefficient way

- Takes Around 2mins15secs to process a file with 1 billion rows of data.

```bash
./out/1brc simple-process --file=largefile.txt
```

Use MMap Command: Read and display the content of a specified file using memory mapping. This command is to process the file using memory mapping and is more efficient than the simple process command.

- Takes Around 2mins10secs to process a file with 1 billion rows of data.

```bash
./out/1brc use-basic-mmap --file=largefile.txt
```

Use Parallel Mmap Command: Read and display the content of a specified file using memory mapping and processing the data parallely.

- This command is to process the file using memory mapping and break into smaller chunks to process them parallely using go routines and is more efficient than the basic mmap process command.

- Takes Around 36secs to process a file with 1 billion rows of data.

```bash
./out/1brc use-parallel-mmap --file=largefile.txt
```
