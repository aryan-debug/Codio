# Codio
A web app that can run multiple languages (Python, Golang, C++ for now)

# Demo

https://github.com/user-attachments/assets/2fa79259-9772-4ac7-a2b4-cde0f38a05fb

## Technologies Used
- SvelteKit
- Golang
- TypeScript
- Docker

## How it Works
1. The backend recieves the code and one of the supported languages.
2. It writes the code to a temporary file, which is then mounted on a Docker container
3. Depending on the language, the appropriate container runs and executes the code in the file
4. If the code times out, then it stops the container and removes it
5. Otherwise, send the output via a Go channel and remove the container


## How to Run It Locally
1. Build the docker images in `/languages` and tag them as `[language_name]_runner`
2. Set environment variables in `/frontend` and `/backend`
3. Make sure the Docker daemon is running
4. Run the backend with `go run main.go`
5. Run the frontend with `yarn run dev`

## Security Considerations

The website is not fully tested. The only security precaution I have taken is to not let the containers run for a long period of time to prevent crypto-mining.
I also don't really plan to make it public (unless you are a recruiter ;) then you have access to it).
