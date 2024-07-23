package main

import (
	"fmt"
	"os"

	"github.com/gnana997/billion-rows-challenge/GenerateDataSet"
	parallelmmapimplementation "github.com/gnana997/billion-rows-challenge/ParallelMmapImplementation"
	simpleprocess "github.com/gnana997/billion-rows-challenge/SimpleProcess"
	mmapimplementation "github.com/gnana997/billion-rows-challenge/mmapImplementation"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{Use: "app"}

	var filePath string
	var rows int

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a file",
		Run: func(cmd *cobra.Command, args []string) {
			if rows == 0 {
				fmt.Println("Please provide number of rows to be generated  --rows")
				os.Exit(1)
			}
			if filePath == "" {
				fmt.Println("Please provide a file path using --file")
				os.Exit(1)
			}
			GenerateDataSet.GenerateDataSet(rows, filePath)
		},
	}
	createCmd.Flags().StringVarP(&filePath, "file", "f", "", "Name of the file to create")
	createCmd.Flags().IntVarP(&rows, "rows", "c", 0, "Content to write to the file")

	readCmd := &cobra.Command{
		Use:   "simple-process",
		Short: "Read a file and process it sequentially",
		Run: func(cmd *cobra.Command, args []string) {
			if filePath == "" {
				fmt.Println("Please provide a file name using --file")
				os.Exit(1)
			}
			simpleprocess.SimplProcessFunc(filePath)
		},
	}
	readCmd.Flags().StringVarP(&filePath, "file", "f", "", "Name of the file to read")

	mmapCmd := &cobra.Command{
		Use:   "use-basic-mmap",
		Short: "Read a file and load it into mmap in memory and process it sequentially",
		Run: func(cmd *cobra.Command, args []string) {
			if filePath == "" {
				fmt.Println("Please provide a file name using --file")
				os.Exit(1)
			}
			mmapimplementation.MmapImplementationFunc(filePath)
		},
	}
	mmapCmd.Flags().StringVarP(&filePath, "file", "f", "", "Name of the file to read")

	paralleMmapCmd := &cobra.Command{
		Use:   "use-parallel-mmap",
		Short: "Read a file and load it into mmap in memory and process it sequentially",
		Run: func(cmd *cobra.Command, args []string) {
			if filePath == "" {
				fmt.Println("Please provide a file name using --file")
				os.Exit(1)
			}
			parallelmmapimplementation.ParallelMmapImplementation(filePath)
		},
	}
	paralleMmapCmd.Flags().StringVarP(&filePath, "file", "f", "", "Name of the file to read")

	rootCmd.AddCommand(createCmd, readCmd, mmapCmd, paralleMmapCmd)
	rootCmd.Execute()
}
