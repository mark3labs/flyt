# Flyt Documentation

This directory contains the documentation website for Flyt, built with [Vocs](https://vocs.dev).

## Development

Install dependencies:

```bash
bun install
# or
npm install
```

Start the development server:

```bash
bun dev
# or
npm run dev
```

The documentation will be available at http://localhost:5173

## Building

Build the documentation for production:

```bash
bun build
# or
npm run build
```

Preview the production build:

```bash
bun preview
# or
npm run preview
```

## Structure

```
www/
├── docs/
│   ├── pages/           # Documentation pages (MDX)
│   │   ├── index.mdx    # Landing page
│   │   ├── getting-started/
│   │   ├── concepts/
│   │   ├── patterns/
│   │   ├── advanced/
│   │   └── examples/
│   └── public/          # Static assets
│       └── flyt-logo.png
├── vocs.config.ts       # Vocs configuration
├── package.json
└── README.md
```

## Adding Documentation

1. Create a new `.mdx` file in the appropriate directory under `docs/pages/`
2. Add the page to the sidebar in `vocs.config.ts`
3. Use MDX features for rich content (code blocks, components, etc.)

## Deployment

The documentation can be deployed to any static hosting service:

- Vercel
- Netlify
- GitHub Pages
- Cloudflare Pages

Build the docs and deploy the `dist` directory.