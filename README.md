# Terminal CV

a cv that you ssh into written in golang

## build image

```
docker build -t terminal-cv -f docker/Dockerfile .
```

## run image

```
docker run -it --rm \
  -p 2222:1337 \
  -v ssh_tui_keys:/app/keys \
  terminal-cv
```
