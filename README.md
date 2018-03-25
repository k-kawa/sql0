# sql0
Sql0 is the simplest commandline tool to check if the given query results in 0.

*This project is still under PoC*

# Install

```sh
go get github.com/k-kawa/sql0/cmd/sql0
```

# Usage

## Prepare Configuration file

Create a configuration file like this.

```sh
cat <<EOF > config.json
{
    "ProjectID": "your_gcp_project",
    "Query": "SELECT * FROM `your_gcp_project.your_dataset.your_table`"
}
EOF

```

## Run sql0 with -c option

```sh
sql0 -c config.json
```

`sql0` executes the query written in the configuration file and fails if the result is not empty.
When the command fails, it reports the following information in JSON format.

| parameter | type | description |
|-----|-----|-----|
| Success | bool | always `false` |
| RowCount | bool | number of the rows hit |
| Time | string | current datetime in ISO8061 format |
| Errors | []string | error information |


