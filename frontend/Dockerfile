FROM node:20-alpine as builder
WORKDIR /app

COPY ./package.json ./
COPY ./yarn.lock ./

RUN yarn install --frozen-lockfile

COPY . .

RUN yarn build

FROM nginx:1.25.1

COPY --from=builder /app/build .
COPY --from=builder /app/nginx/default.conf /etc/nginx/templates/default.conf.template

CMD ["nginx", "-g", "daemon off;"]