# Pre-baked Image Description JSON Generator

This script generates JSON files containing metadata about container images, which can be used to facilitate load testing.

## Prerequisites

- Ensure you have `oras` installed. You can install it by following the instructions [here](https://oras.land/docs/installation).

## Usage

1. Clone the repository and navigate to the `prepare` directory:

    ```sh
    git clone <repository-url>
    cd prepare
    ```

1. Ensure you have a `repositories.json` file in the prepare directory. This file should contain the repositories you want to process. Example format:

    ```json
    {
        "repositories": [
            "repo1",
            "repo2",
            "repo3"
        ]
    }
    ```

The repositories checked-in and generated via `oras repo list mcr.azk8s.cn` on 2025-3-13.

1. Run the `generate_json.sh` script to generate the JSON files:

    ```sh
    ./generate_json.sh
    ```

    The script will:
    - Read the repositories from `repositories.json`.
    - Fetch the tags for each repository.
    - Fetch the manifest and layers for each tag.
    - Generate a JSON file for each image containing the size, manifest, and blob information.

4. The generated JSON files will be saved in the `json_files` directory.

## How this helps

With the generated JSON files, the required resource on the test client side will be considerably reduced for load testing.

- The image manifest and layers can be pulled in parallel, reducing the wait time for the manifest to be pulled. This significantly improves the overall number of outgoing requests per second.
- Adapted to the UX of `oras blob` and `oras manifest` commands. One can easily ditch the downloaded content to save memory and IO resources.
