package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	flag "github.com/spf13/pflag"
)

var (
	Version = "dev"
	Commit  = "dev"
)

func main() {

	var sizeStr string
	var recordNum int64
	var columnNum int
	var outputPath string
	var help bool

	flag.Int64VarP(&recordNum, "num", "n", 0, "Number of records.")
	flag.StringVarP(&sizeStr, "size", "s", "", "Size. Can be specified in KB, MB, or GB. (ex. 1gb)")
	flag.IntVarP(&columnNum, "col", "c", 100, "Number of columns.")
	flag.StringVarP(&outputPath, "output", "o", "", "Output file path.")
	flag.BoolVarP(&help, "help", "h", false, "Help.")
	flag.Parse()
	flag.CommandLine.SortFlags = false
	flag.Usage = func() {
		fmt.Printf("createcsv v%s (%s)\n\n", Version, Commit)
		fmt.Fprint(os.Stderr, "Usage:\n  createcsv [flags]\n\nFlags\n")
		flag.PrintDefaults()
	}

	if help {
		flag.Usage()
		os.Exit(0)
	}

	if (recordNum == 0 && sizeStr == "") || outputPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	start := time.Now()

	totalRecordNum, totalSize, err := run(recordNum, sizeStr, columnNum, outputPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(
		"Created.\nTime(seconds): %s / Number of records: %s / Number of size(byte): %s",
		humanize.Commaf(float64(time.Since(start).Milliseconds())/1000),
		humanize.Comma(totalRecordNum),
		humanize.Comma(totalSize))

}

func run(recordNum int64, sizeStr string, columnNum int, outputPath string) (int64, int64, error) {

	size, err := parseSize(sizeStr)
	if err != nil {
		return 0, 0, err
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	totalSize := int64(0)
	totalRecordNum := int64(0)

	// ヘッダ出力
	columns := []string{}
	for i := 1; i <= columnNum; i++ {
		// col1,col2,col3... といった形で
		columns = append(columns, fmt.Sprintf("col%d", i))
	}

	writeSize, err := writeRecord(writer, columns)
	if err != nil {
		return 0, 0, err
	}

	totalSize += int64(writeSize)

	for (size != 0 && totalSize < size) || (recordNum != 0 && totalRecordNum < recordNum) {
		writeSize, err := writeRecord(writer, randomValues(columnNum))
		if err != nil {
			return 0, 0, err
		}

		totalSize += int64(writeSize)
		totalRecordNum++
	}

	return totalRecordNum, totalSize, writer.Flush()
}

func writeRecord(writer *bufio.Writer, items []string) (int, error) {

	record := strings.Join(items, ",") + "\n"

	return writer.WriteString(record)
}

func randomValue() string {
	// 同じ条件で作成した際に同じ内容になるように、乱数のシードは変えずに実行
	return strconv.Itoa(int(rand.Int31()))
}

func randomValues(num int) []string {
	values := []string{}
	for i := 0; i < num; i++ {
		values = append(values, randomValue())
	}
	return values
}

func parseSize(sizeStr string) (int64, error) {

	if sizeStr == "" {
		return 0, nil
	}

	sizeStr = strings.ToUpper(sizeStr)

	if strings.HasSuffix(sizeStr, "GB") {
		baseSize, err := parseInt64(sizeStr[:len(sizeStr)-len("GB")])
		if err != nil {
			return -1, err
		}

		return baseSize * 1024 * 1024 * 1024, nil
	}

	if strings.HasSuffix(sizeStr, "MB") {
		baseSize, err := parseInt64(sizeStr[:len(sizeStr)-len("MB")])
		if err != nil {
			return -1, err
		}

		return baseSize * 1024 * 1024, nil
	}

	if strings.HasSuffix(sizeStr, "KB") {
		baseSize, err := parseInt64(sizeStr[:len(sizeStr)-len("KB")])
		if err != nil {
			return -1, err
		}

		return baseSize * 1024, nil
	}

	return parseInt64(sizeStr)
}

func parseInt64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}
