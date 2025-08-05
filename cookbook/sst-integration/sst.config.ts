/// <reference path="./.sst/platform/config.d.ts" />

export default $config({
  app(input) {
    return {
      name: "sst-integration",
      removal: "remove",
      home: "aws",
    };
  },
  async run() {
    // Create buckets
    const sourceBucket = new sst.aws.Bucket("SourceImages");
    const destBucket = new sst.aws.Bucket("ExtractedText");
    
    // Create secret for API key
    const anthropicKey = new sst.Secret("AnthropicApiKey");

    // Create Lambda function with linked resources
    const ocrFunction = new sst.aws.Function("ImageOCR", {
      runtime: "go",
      handler: "./src",
      link: [sourceBucket, destBucket, anthropicKey],
      timeout: "2 minutes",
    });

    // Add S3 trigger
    sourceBucket.notify({
      notifications: [
        {
          name: "ImageProcessor",
          function: ocrFunction.arn,
          events: ["s3:ObjectCreated:*"],
        },
      ],
    });
    
    return {
      sourceBucket: sourceBucket.name,
      destBucket: destBucket.name,
    };
  },
});