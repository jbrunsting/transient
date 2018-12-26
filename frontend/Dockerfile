FROM node:8.0.0-alpine

ENV NODE_ENV development

RUN mkdir /src
ADD . /src
WORKDIR /src
RUN rm package-lock.json
RUN rm -r node_modules

EXPOSE 8080

CMD npm install && npm run serve
