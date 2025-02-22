server {
    listen 80;
    server_name localhost;

    # User Service
    location /api/users {
        proxy_pass http://user-service:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # Order Service
    location /api/orders {
        proxy_pass http://order-service:8081;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}