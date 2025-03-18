FROM tiangolo/nginx-rtmp

RUN apt-get update && apt-get install -y ffmpeg

COPY nginx.conf /etc/nginx/nginx.conf

COPY index.html /usr/share/nginx/html/index.html

EXPOSE 80
