## Run locally

```
docker-compose up
```

## Debug remotely

The docker files located in `docker/debug` are designed to create a container with Delve installed and executing the application binary.

```
cd docker/debug
docker-compose up
```

If using VS Code, use the `Attach Remote` config located in `.vscode/launch.json` to attach the IDE debugger to the container to enable breakpoints.

## Migrations (WIP)

```
docker-compose run server tern migrate --migrations migrations --config migrations/tern.conf
```
