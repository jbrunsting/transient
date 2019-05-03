FROM node:8.16.0-alpine

ENV NODE_ENV development

RUN mkdir /src
ADD . /src
WORKDIR /src
RUN rm -rf node_modules

EXPOSE 8080

CMD npm install && npm run serve
