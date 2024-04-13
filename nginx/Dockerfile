# Nginx setup
FROM nginx:latest

COPY --from=build /app/dist /app

COPY private/ /etc/nginx/ssl
COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 443

CMD ["nginx", "-g", "daemon off;"]