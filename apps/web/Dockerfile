FROM node:20 as builder

RUN npm install -g pnpm

WORKDIR /app
COPY package.json /app
COPY pnpm-lock.yaml /app

RUN pnpm install

COPY . /app
RUN pnpm build

FROM nginx:alpine

COPY --from=builder /app/deployment/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80
