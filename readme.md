# Personal Webpage

A modern, responsive personal webpage built with Go, Echo, Templ, GoMail and TailwindCSS.

## Tech Stack

- [Go](https://go.dev/) - Backend server
- [Echo](https://echo.labstack.com/) - Web framework
- [Templ](https://templ.guide/) - HTML templating
- [TailwindCSS](https://tailwindcss.com/) - Styling
- [Air](https://github.com/cosmtrek/air) - Live reloading
- [Task](https://taskfile.dev/) - Task runner
- [Docker](https://www.docker.com/) - Containerization

## Development

Start the development server with hot-reload:

```sh
task dev
```

Build the project:

```sh
task build
```

Build Docker container:

```sh
task container container_name=<name> container_tag=<tag>
```

## Project Structure

```
.
├── assets/         # Static assets and SVG icons
├── components/     # Reusable UI components
├── style/         # TailwindCSS config and styles
├── views/         # Page templates
├── main.go        # Entry point
└── Dockerfile     # Container definition
```

## Deployment

Automated deployment via GitHub Actions:

- Builds Docker image on push to master
- Pushes to Docker Hub
- Deploys to VPS via SSH

## Docker

Built with a multi-stage build:

- Go builder stage compiles binary
- Final stage uses distroless base image
- Exposes port 8080

## License

This project is open source and available under the MIT license.

## Want my badge? No problems:

add this file to your static assets:
`assets/js/augustg-dev-badge.js`

```
<augustg-dev-badge></augustg-dev-badge>
<script src="/assets/js/augustg-dev-badge.js"></script>
```

Badge is using this font:

https://github.com/0xViva/webpage/blob/master/assets/fonts/JetBrainsMonoNerdFontMono-Regular.ttf



