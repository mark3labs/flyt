# SST Integration with Flyt

This example demonstrates how to integrate Flyt workflows with SST v3 to build a serverless image text extraction service.

## Overview

This application processes images uploaded to S3 by:
1. Triggering a Lambda function on image upload
2. Using Flyt's node-based workflow to orchestrate the process
3. Extracting text from images using Anthropic's Claude Vision API
4. Saving the extracted text to a separate S3 bucket

## Architecture

The application uses:
- **SST v3** for infrastructure and resource management
- **Flyt** for workflow orchestration with three nodes:
  - `DownloadNode`: Downloads the image from S3
  - `ExtractTextNode`: Calls Claude Vision API to extract text
  - `SaveTextNode`: Saves the extracted text to S3
- **AWS Lambda** for serverless compute
- **S3** for storage (source images and extracted text)
- **Anthropic Claude** for AI-powered text extraction

## Prerequisites

- Node.js 18+ and npm
- Go 1.21+
- AWS CLI installed and configured
- AWS account with appropriate permissions
- Anthropic API key

## Setup

1. Clone the repository and navigate to this example:
```bash
cd cookbook/sst-integration
```

2. Install dependencies:
```bash
npm install
go mod tidy
```

3. Initialize SST (if not already initialized):
```bash
npx sst init
```

4. Set your Anthropic API key as a secret:
```bash
npx sst secret set AnthropicApiKey sk-ant-api_xxxxxxxxxxxxx
```

5. Deploy the application:
```bash
npx sst deploy
```

The deployment will output the bucket names:
```
✓ Complete
   sourceBucket: sst-integration-sourceimages-xxxxx
   destBucket: sst-integration-extractedtext-xxxxx
```

## Usage

### Upload an Image

Use AWS CLI to upload an image to the source bucket:

```bash
# Upload a single image
aws s3 cp /path/to/your/image.jpg s3://sst-integration-sourceimages-xxxxx/

# Upload multiple images
aws s3 cp /path/to/images/ s3://sst-integration-sourceimages-xxxxx/ --recursive --exclude "*" --include "*.jpg" --include "*.png" --include "*.jpeg" --include "*.webp"
```

### Check Results

List the extracted text files:
```bash
aws s3 ls s3://sst-integration-extractedtext-xxxxx/
```

Download and view the extracted text:
```bash
aws s3 cp s3://sst-integration-extractedtext-xxxxx/image.txt ./
cat image.txt
```

### Monitor Execution

View Lambda logs to monitor the workflow:
```bash
npx sst logs --function ImageOCR
```

Or use AWS CLI:
```bash
aws logs tail /aws/lambda/sst-integration-ImageOCR --follow
```

## Managing Secrets

### View current secrets:
```bash
npx sst secret list
```

### Update the API key:
```bash
npx sst secret set AnthropicApiKey sk-ant-api_new_key_here
```

### Remove a secret:
```bash
npx sst secret remove AnthropicApiKey
```

## Cleanup

To remove all resources:
```bash
npx sst remove
```

This will delete:
- Both S3 buckets and their contents
- The Lambda function
- All associated IAM roles and policies
- The stored secret

## Key Features

- **Resource Linking**: Uses SST v3's resource linking to access S3 buckets and secrets without environment variables
- **Type-Safe Configuration**: Leverages SST's type-safe configuration
- **Node-Based Workflow**: Demonstrates Flyt's composable node architecture
- **Error Handling**: Each node can handle errors independently
- **Scalable**: Automatically scales with Lambda

## Project Structure

```
cookbook/sst-integration/
├── src/
│   ├── main.go      # Lambda handler with Flyt workflow
│   ├── go.mod       # Go dependencies
│   └── go.sum       # Go dependency checksums
├── sst.config.ts    # SST infrastructure configuration
├── package.json     # Node.js dependencies
└── README.md        # This file
```

## How It Works

1. **S3 Event Trigger**: When an image is uploaded to the source bucket, S3 sends an event to the Lambda function
2. **Flyt Workflow**: The Lambda handler creates a Flyt flow with three connected nodes
3. **Download**: The first node downloads the image from S3 using the AWS SDK
4. **Extract Text**: The second node sends the image to Claude Vision API for text extraction
5. **Save Result**: The final node saves the extracted text to the destination S3 bucket

## Monitoring

View Lambda logs in CloudWatch to monitor the workflow execution and debug any issues.