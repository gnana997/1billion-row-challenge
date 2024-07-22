## Terminal App for 1 Billion Row Challenge

# Description

Just building this terminal application to have multiple solutions with different approaches to solve the 1 Billion Row Challenge. The challenge is to read a file with 1 billion rows of city name and it's reported temperature and print the cities with their min/avg/max temperature.

# Features

Create Command: Quickly create files with a specified name and content. This command is optimized to handle large inputs and can be used to generate files with a substantial amount of data.

```bash
./1brc create --file=largefile.txt --rows=1000000000
```

Simple Process Command: Read and display the content of a specified file. This command is just to process the file in sequential and in very inefficient way

```bash
./1brc simple-process --file=largefile.txt
```
