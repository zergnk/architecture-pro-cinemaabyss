FROM amazoncorretto:21-alpine-jdk AS builder
WORKDIR /app
COPY . .
RUN ./gradlew  clean  build

FROM amazoncorretto:21-alpine-jdk
RUN apk --no-cache add curl
WORKDIR /app
COPY --from=builder /app/build/libs/events-0.0.1-SNAPSHOT.jar app.jar
EXPOSE 8082
ENTRYPOINT ["java","-jar", "app.jar"]
