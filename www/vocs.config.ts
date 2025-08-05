import { defineConfig } from 'vocs'

export default defineConfig({
  title: 'Flyt',
  description: 'A minimalist workflow framework for Go with zero dependencies',
  baseUrl: 'https://mark3labs.github.io/flyt',
  basePath: '/',
  logoUrl: '/flyt-logo.png',
  sidebar: [
    {
      text: 'Introduction',
      link: '/',
    },
    {
      text: 'Getting Started',
      items: [
        {
          text: 'Installation',
          link: '/getting-started/installation',
        },
        {
          text: 'Quick Start',
          link: '/getting-started/quick-start',
        },
        {
          text: 'Project Template',
          link: '/getting-started/template',
        },
      ],
    },
    {
      text: 'Core Concepts',
      items: [
        {
          text: 'Nodes',
          link: '/concepts/nodes',
        },
        {
          text: 'Actions',
          link: '/concepts/actions',
        },
        {
          text: 'Flows',
          link: '/concepts/flows',
        },
        {
          text: 'Shared Store',
          link: '/concepts/shared-store',
        },
      ],
    },
    {
      text: 'Patterns',
      items: [
        {
          text: 'Configuration via Closures',
          link: '/patterns/closures',
        },
        {
          text: 'Error Handling & Retries',
          link: '/patterns/error-handling',
        },
        {
          text: 'Fallback on Failure',
          link: '/patterns/fallback',
        },
        {
          text: 'Conditional Branching',
          link: '/patterns/branching',
        },
      ],
    },
    {
      text: 'Advanced',
      items: [
        {
          text: 'Custom Node Types',
          link: '/advanced/custom-nodes',
        },
        {
          text: 'Batch Processing',
          link: '/advanced/batch-processing',
        },
        {
          text: 'Batch Flows',
          link: '/advanced/batch-flows',
        },
        {
          text: 'Nested Flows',
          link: '/advanced/nested-flows',
        },
        {
          text: 'Flow as Node',
          link: '/advanced/flow-as-node',
        },
        {
          text: 'Worker Pool',
          link: '/advanced/worker-pool',
        },
        {
          text: 'Utilities',
          link: '/advanced/utilities',
        },
      ],
    },
    {
      text: 'Examples',
      items: [
        {
          text: 'Agent',
          link: 'https://github.com/mark3labs/flyt/tree/main/cookbook/agent',
        },
        {
          text: 'Chat',
          link: 'https://github.com/mark3labs/flyt/tree/main/cookbook/chat',
        },
        {
          text: 'LLM Streaming',
          link: 'https://github.com/mark3labs/flyt/tree/main/cookbook/llm-streaming',
        },
        {
          text: 'MCP',
          link: 'https://github.com/mark3labs/flyt/tree/main/cookbook/mcp',
        },
        {
          text: 'Summarize',
          link: 'https://github.com/mark3labs/flyt/tree/main/cookbook/summarize',
        },
        {
          text: 'Tracing',
          link: 'https://github.com/mark3labs/flyt/tree/main/cookbook/tracing',
        },
      ],
    },
    {
      text: 'Best Practices',
      link: '/best-practices',
    },
  ],
  socials: [
    {
      icon: 'github',
      link: 'https://github.com/mark3labs/flyt',
    },
  ],
})
