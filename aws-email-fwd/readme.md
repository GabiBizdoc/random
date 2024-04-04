# Email Forwarding

This is a simple Lambda function that forwards emails from S3 buckets using AWS Simple Notification Service (SNS)

## Setup

### Environment Variables

Make sure to set the following environment variables:

- `REGION`: The AWS region where the application will be deployed.
- `RECEIVER_EMAIL`: The email address where forwarded emails will be sent.
- `SENDER_EMAIL`: The email address from which forwarded emails will be sent.

## Note

This application assumes that a SES receipt rule is set to deliver incoming emails
to the designated S3 bucket, and an S3 event is configured to trigger this Lambda handler.

## Appendix

### Golang installation

Please ensure Go 1.x (where 'x' is the latest version) is installed as per the instructions on the official golang
website: https://golang.org/doc/install

A quickstart way would be to use Homebrew, chocolatey or your linux package manager.

#### Homebrew (Mac)

Issue the following command from the terminal:

```shell
brew install golang
```

If it's already installed, run the following command to ensure it's the latest version:

```shell
brew update
brew upgrade golang
```

#### Chocolatey (Windows)

Issue the following command from the powershell:

```shell
choco install golang
```

If it's already installed, run the following command to ensure it's the latest version:

```shell
choco upgrade golang
```
