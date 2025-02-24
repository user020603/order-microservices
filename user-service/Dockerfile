# Use official Maven image to build the project
FROM maven:3.8.4-openjdk-17 AS build
WORKDIR /app

# Copy the project files
COPY user /app

# Build the application
RUN mvn clean package -DskipTests

# Use OpenJDK as the runtime environment
FROM openjdk:17-jdk-slim
WORKDIR /app

# Copy the built JAR file from the build stage
COPY --from=build /app/target/*.jar app.jar

# Expose the application port
EXPOSE 8080

# Command to run the application
ENTRYPOINT ["java", "-jar", "app.jar"]
