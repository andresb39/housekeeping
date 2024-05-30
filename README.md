# Cloud Cost Optimization Tool

This tool is designed to help you identify and calculate the costs of unused or orphaned resources in AWS and Azure. By identifying these resources, you can optimize your cloud spending and save costs by removing them.

## Purpose

The primary purpose of this tool is to:

1. Identify unused or orphaned resources in AWS and Azure.
2. Calculate the cost associated with these resources.
3. Provide potential cost savings if these resources are removed.

## Supported Cloud Providers

- AWS
- Azure

## Usage

### Prerequisites

- Go 1.18 or later
- AWS and Azure SDKs for Go
- Proper AWS and Azure credentials configured

## Instalation

To install the tool, you can copy it to your `bin` directory:

Move the binary to your `bin` directory:

```sh
mv cmd/finops /usr/local/bin/
```

This will make the `finops` tool accessible from anywhere in your terminal.

## Running the Tool

The tool can be run from the command line. Below are the instructions for running it for both AWS and Azure.

### For AWS

To run the tool for AWS, use the following command:

```sh
finops --provider aws --shutdownMonths 6 --snapshotRetentionMonths 6
```

Parameters:

- `provider`: The cloud provider (`aws`).
- `shutdownMonths`: Number of months before considering instances as orphaned (optional, default is `6`).
- `snapshotRetentionMonths`: Number of months to retain snapshots (optional, default is `6`).
- `region`: Region to check for resources (optional, default is `us-east-1`).
- `deleteFlag`: Flag for delete orphaned resources (optional, default is `false`).

### For Azure

To run the tool for Azure, use the following command:

```sh
finops --provider azure --subscriptionID "your-subscription-id" --shutdownMonths 6 --snapshotRetentionMonths 6
```

Parameters:

- `provider`: The cloud provider (`azure`).
- `subscriptionID`: Your Azure subscription ID (required).
- `shutdownMonths`: Number of months before considering VMs as orphaned (optional, default is `6`).
- `snapshotRetentionMonths`: Number of months to retain snapshots (optional, default is `6`).
- `deleteFlag`: Flag for delete orphaned resources (optional, default is `false`).

### Example Commands

#### AWS

```sh
finops --provider aws
```

This command will use the default values for the cost parameters.

#### Azure

```sh
finops --provider azure --subscriptionID "your-subscription-id"
```

This command will use the default values for the cost parameters.

### Notes

Ensure your AWS and Azure credentials are correctly configured and accessible by the tool.
Review and adjust the cost parameters as needed to reflect your actual costs.

## Developer Contributions

### Setting Up the Development Environment

1. Clone the repository:

   ```sh
   git clone git@github.com:andresb39/housekeeping.git
   cd housekeeping
   ```

2. Ensure you have Go 1.22.3 or later installed. You can check your Go version with:

   ```sh
   go version
   ```

3. Install the necessary Go modules:

   ```sh
   go mod tidy
   ```

### Running the Tool Locally

To run the tool locally during development, you can use the following commands.

For AWS:

```sh
go run main.go --provider aws --shutdownMonths 6 --snapshotRetentionMonths 6
```

For Azure:

```sh
go run main.go --provider azure --subscriptionID "your-subscription-id" --shutdownMonths 6 --snapshotRetentionMonths 6
```

### Building the Binary

To build the binary for distribution:

```sh
GOARCH=amd64 go build -o cmd/finops main.go
```

### Submitting Contributions

1. Create a new branch for your feature or bugfix:

   ```sh
   git checkout -b my-feature-branch
   ```

2. Make your changes and commit them with a meaningful commit message:

   ```sh
   git add .
   git commit -m "Description of the feature or fix"
   ```

3. Push your branch to the remote repository:

   ```sh
   git push origin my-feature-branch
   ```

4. Open a pull request on GitHub, describing your changes and the purpose of the pull request.

### Coding Guidelines

- Follow the Go coding standards and best practices.
- Write clear, concise, and well-documented code.
- Ensure your code is tested and does not break existing functionality.

Thank you for contributing to the Cloud Cost Optimization Tool!
