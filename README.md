# Personal Website

A personal portfolio and blog website built with Go, HTMX, and Templ. Features a markdown-based blog system, reading list tracker, and creative showcase page.

## Tech Stack

**Backend:** Go 1.23, net/http, PostgreSQL (pgx), Templ
**Frontend:** HTMX, Tailwind CSS, DaisyUI
**Infrastructure:** Docker, Traefik (reverse proxy/TLS), Watchtower

## Project Structure

```
cmd/
├── run/          # Main server application
└── upload/       # Data upload utility (posts, books)

internal/
├── db/           # Database layer and schemas
├── ds/           # Custom data structures (Set, OrderedList, StrictDict)
├── handlers/     # HTTP request handlers
└── pub/          # Templates, assets, and styles
    ├── pages/    # Full page templates
    ├── partials/ # Reusable components
    ├── shared/   # Layout templates
    └── assets/   # Static files (CSS, images, fonts)

data/
├── blog_posts/   # Markdown blog post files
├── books.json    # Reading list data
└── resources.md  # Resources page content
```

## Development

**Prerequisites:** Go 1.23+, Node.js, PostgreSQL

```bash
# Install dependencies
npm install

# Run development server (Tailwind watch + Go server + browser-sync)
npm run dev

# Or run individual processes
npm run tailwind       # Watch Tailwind compilation
npm run server         # Watch Go server with nodemon
npm run templ          # Watch Templ template compilation
```

## Building

```bash
# Build binaries and Docker image
make build

# Run locally
make run
```

## Deployment

The project uses Docker Compose with Traefik for automatic TLS and load balancing:

```bash
docker compose up -d
```

Services:
- **server:** Go application (3 replicas)
- **traefik:** Reverse proxy with Let's Encrypt TLS
- **watchtower:** Automatic container updates
