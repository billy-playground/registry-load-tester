# Registry Load Tester

This project is a load testing tool for a registry. It allows you to pre-bake workload and simulate concurrent instances to pull from a registry to evaluate its performance.

## Usage

Run the tool using the following command:

```bash
go run main.go <num_instances>[=<size>/<interval>] <registry_domain> <token_mode> [<registry_endpoint>]
```

### Parameters

- **num_instances**: Number of instances to run. Optionally, specify batch size and interval.
- **registry_domain**: Domain of the registry.
- **token_mode**: Token mode. Options:
  - `none`: No token.
  - `anonymous`: Obtain an anonymous registry token to share between instances.
  - `token=<token>`: Use the specified identity token to exchange a registry toekn and share between instances.
- **registry_endpoint** (optional): Endpoint of the registry. Defaults to the registry domain.

### Example

```bash
#  runs 10 instances against `registry.example.com` without using any token.
go run main.go 10 registry.example.com none

# runs 100 instances against `registry.example.com`, starting 10 instances every 500 milliseconds using the specified token.
go run main.go 100=10/500ms registry.example.com token=$registry_token

# runs 50 instances against `registry.example.com` using shared anonymous access.
go run main.go 50 registry.example.com anonymous

# runs 20 instances against `registry.example.com` via a custom endpoint.
go run main.go 20 registry.example.com none cus.fe.example.com
```

## Notes

- Ensure the `assets/images` directory contains JSON files before running the tool. Refer to the [prepare tool instruction](prepare/README.md) for more details.
- The tool outputs performance metrics in CSV format to the stdout.
