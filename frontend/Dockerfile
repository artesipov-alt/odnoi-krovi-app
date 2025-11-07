# frontend/Dockerfile
FROM caddy:2-alpine
WORKDIR /srv
# Build directory is mounted via volume in docker-compose.yml
# This copy provides initial files, volume allows updates without rebuild
# COPY ./build /srv
COPY Caddyfile /etc/caddy/Caddyfile
