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
- **Flyt** - Workflow orchestration with batch flow processing:
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

- **Concurrent Processing**: Uses Flyt batch flows to process multiple images in parallel
- **Smart Transcription**: Claude generates both content and descriptive titles
- **Automatic Cleanup**: Processed images are deleted to save storage
- **Error Recovery**: Retries failed transcriptions up to 3 times
- **Format Support**: Handles PNG and JPEG images (skips unsupported formats)
- **SST Resource Linking**: Clean access to S3 and secrets without env vars
- **Progress Tracking**: Real-time progress updates with completion percentages

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
- **Batch Flow Processing**: Processes multiple images concurrently for better performance
- **Retry Logic**: Failed transcriptions retry up to 3 times
- **Fallback Handling**: Gracefully skips to cleanup on permanent failures
- **Format Validation**: Only processes PNG and JPEG images
- **Progress Monitoring**: Real-time batch progress with completion statistics

## Batch Flow Implementation

This example uses Flyt's batch flow feature to process multiple S3 images concurrently, providing better performance and progress tracking.

### Key Changes from Sequential Processing

1. **Concurrent Processing**: Multiple images are processed in parallel instead of sequentially
2. **Progress Tracking**: Real-time progress updates showing completed and failed items
3. **Better Error Handling**: Individual failures don't stop the entire batch
4. **Resource Isolation**: Each image gets its own flow instance with isolated SharedStore

### How Batch Flow Works

1. **Flow Factory**: Creates a new flow instance for each S3 record
2. **Batch Function**: Generates inputs (bucket, key) for each record
3. **Concurrent Execution**: Processes multiple images simultaneously
4. **Progress Tracking**: Monitors completion and failure rates

### Performance Benefits

- **Faster Processing**: Concurrent execution reduces total processing time
- **Better Resource Utilization**: Makes full use of Lambda's available resources
- **Scalability**: Handles large batches efficiently

### Monitoring Batch Progress

The logs now show real-time progress:
```
Batch progress: 33.3% (1/3) - Completed: 1, Failed: 0
Batch progress: 66.7% (2/3) - Completed: 2, Failed: 0
Batch progress: 100.0% (3/3) - Completed: 3, Failed: 0
```

### Error Handling

- Individual image failures don't stop the batch
- Final summary shows total completed vs failed
- Lambda returns error if any records failed