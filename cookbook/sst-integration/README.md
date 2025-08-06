# SST Integration with Flyt

This example demonstrates how to integrate Flyt workflows with SST v3 to build a serverless handwritten note transcription service.

## Overview

Upload handwritten notes as images to S3 and automatically:
1. Extract and transcribe text using Claude Vision
2. Generate descriptive titles for the notes
3. Save transcriptions as markdown files
4. Clean up processed images

## Architecture

- **SST v3** - Infrastructure as code and resource management
- **Flyt** - Workflow orchestration with four nodes:
  - `DownloadNode` - Downloads image and detects MIME type
  - `ExtractTextNode` - Transcribes using Claude Vision (with retry)
  - `SaveTextNode` - Saves markdown with generated title
  - `CleanupNode` - Deletes processed image
- **AWS Lambda** - Serverless compute
- **S3** - Image storage and markdown output
- **Claude Sonnet 4** - AI-powered transcription

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

List the transcribed markdown files:
```bash
aws s3 ls s3://sst-integration-extractedtext-xxxxx/
```

Download and view a transcription:
```bash
# Files are named based on the content, e.g., "Meeting_Notes_Q4_Planning.md"
aws s3 cp s3://sst-integration-extractedtext-xxxxx/Your_Note_Title.md ./
cat Your_Note_Title.md
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

- **Smart Transcription**: Claude generates both content and descriptive titles
- **Automatic Cleanup**: Processed images are deleted to save storage
- **Error Recovery**: Retries failed transcriptions up to 3 times
- **Format Support**: Handles PNG and JPEG images (skips unsupported formats)
- **SST Resource Linking**: Clean access to S3 and secrets without env vars

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

1. **Upload** an image to the source bucket
2. **Download** node fetches the image and detects its type
3. **Extract** node sends to Claude with this prompt:
   - Transcribe the handwritten note with markdown formatting
   - Generate a descriptive title
   - Return as JSON with "content" and "title" fields
4. **Save** node creates a markdown file named after the title
5. **Cleanup** node deletes the original image

The workflow includes:
- **Retry Logic**: Failed transcriptions retry up to 3 times
- **Fallback Handling**: Gracefully skips to cleanup on permanent failures
- **Format Validation**: Only processes PNG and JPEG images