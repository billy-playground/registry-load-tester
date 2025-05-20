# Registry Load Tester

This project is a load testing tool for a registry. It allows you to pre-bake workload and simulate concurrent instances to pull from a registry to evaluate its performance.

## Build

Before running the tool, build the binary:

```bash
make build
```

## Usage

Run the tool using the following command:

```bash
go run main.go <num_instances>[=<size>/<interval>] <registry_domain> <token_mode> [<registry_endpoint>]
```

### Auth command

`auth` command cab be used to run authentication-related workloads against a registry. Please refer to `rlt auth -h` for more details.

### Pull command

`pull` command can be used to run image pulling workloads against a registry. Please refer to `rlt pull -h` for more details.

## Notes

- Ensure the `assets/images` directory contains JSON files before running the tool. Refer to the [prepare tool instruction](prepare/README.md) for more details.
- The tool outputs performance metrics in CSV format to the stdout.
