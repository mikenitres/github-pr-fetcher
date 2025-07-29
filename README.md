# GitHub PR Fetcher

GitHub PR Fetcher is a command-line tool written in [Go](https://go.dev/) that allows users to fetch the five most recent merged pull requests from specified GitHub repositories on the `main` branch by default.

## Features

- Fetches the latest merged pull requests from multiple GitHub repositories.
- Supports authentication with GitHub API.
- Displays the results in your terminal.

## Build

To install the GitHub PR Fetcher, follow these steps:

1. Clone the repository:

   ```
   git clone git@github.com:mikenitres/github-pr-fetcher.git
   ```

2. Navigate to the project directory:

   ```
   cd github-pr-fetcher
   ```

3. Build the project:

   ```
   go build -o github-pr-fetcher ./cmd/main.go
   ```

4. Optionally, move the binary to a directory in your PATH for easier access. For e.g:

   ```
   mv github-pr-fetcher /usr/local/bin/
   ```

## Configuration

Before using the tool, you need to set up your GitHub API token. In [config.json](config.json), edit the variable `github_token` and place your own Github token.

## Usage

To use the GitHub PR Fetcher, run the following command:

```
github-pr-fetcher -r "<repo1>,<repo2> ..."
```

Replace `<repo1>`, `<repo2>`, etc., with the GitHub repositories you want to query. The repositories should be in the format `owner/repo`.

Example:

```
github-pr-fetcher -r "microsoft/vscode,microsoft/TypeScript"
```

This command will fetch and display the five most recent merged pull requests from the `main` branch of the specified repositories. To specify the branch, say `develop`, do so like this

```
github-pr-fetcher -r "microsoft/vscode,microsoft/TypeScript" -b "develop"
```



## License

[Do What The F*ck You Want To Public License](LICENSE)