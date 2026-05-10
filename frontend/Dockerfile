FROM node:22-alpine

WORKDIR /app

ARG BACKEND_URL
ENV BACKEND_URL=${BACKEND_URL}

COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

CMD sh -c "[ \"$(ls -A /app/dist/)\" ] && { rm -rf /static/*; cp -r /app/dist/. /static/; } || :"
