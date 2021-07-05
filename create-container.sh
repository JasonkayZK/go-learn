# Create Container
docker run -dit \
  --name go-v.17 \
  -v /root/workspace/go-v1.17-code:/code \
  --privileged \
  golang:1.17-rc /bin/bash

# Attach Container
docker exec -it go-v.17 bash
