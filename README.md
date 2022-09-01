# createcsv

The createcsv can create csv of a specified size or number of record.  
You can easily create large csv file.

## Usage

```
$ createcsv -s 10gb -o 10gb.csv
```

```
$ createcsv -n 1000 -o 1000record.csv
```
 
The arguments are as follows.

```
Usage:
  createcsv [flags]

Flags
  -n, --num int         Number of records.
  -s, --size string     Size. Can be specified in KB, MB, or GB. (ex. 1gb)
  -c, --col int         Number of columns. (default 100)
  -o, --output string   Output file path.
  -h, --help            Help.
```

## Install

You can download the binary from the following.

* https://github.com/onozaty/createcsv/releases/latest

## License

MIT

## Author

[onozaty](https://github.com/onozaty)
