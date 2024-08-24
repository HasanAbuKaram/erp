# erp

remove all build cache:
docker builder prune --all


To clean up Docker cache files and reclaim disk space, you can use the following commands:

### 1. Remove Unused Images
To remove unused Docker images, including dangling ones (images that are not tagged and not used by any container):
```bash
docker image prune
```
To remove all unused images (not just dangling ones):
```bash
docker image prune -a
```

### 2. Remove Stopped Containers
To remove stopped containers:
```bash
docker container prune
```

### 3. Remove Unused Volumes
To remove unused volumes (volumes not referenced by any container):
```bash
docker volume prune
```

### 4. Remove Unused Networks
To remove unused networks (networks not used by any container):
```bash
docker network prune
```

### 5. Remove All Unused Data
To remove all unused data (images, containers, volumes, and networks):
```bash
docker system prune
```
To include unused images (not just dangling ones):
```bash
docker system prune -a
```

### 6. Remove Build Cache
To remove build cache used during image builds:
```bash
docker builder prune
```
To remove all build cache:
```bash
docker builder prune --all
```

### 7. Remove Specific Docker Resources
If you want to remove a specific resource, you can use:
- **Containers**: `docker rm <container_id>`
- **Images**: `docker rmi <image_id>`
- **Volumes**: `docker volume rm <volume_name>`
- **Networks**: `docker network rm <network_name>`

### Example Cleanup
To run a general cleanup and remove unused resources:
```bash
docker system prune -a --volumes
```
This will remove all unused images, containers, volumes, and networks, and prompt for confirmation before deleting.

Make sure to review what will be deleted before confirming the cleanup to avoid removing important resources.