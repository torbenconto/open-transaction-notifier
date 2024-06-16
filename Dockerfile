# Use the official Rust image from the Docker Hub
FROM rust:1.79.0

# Create a new directory to run our app
WORKDIR /usr/src/app

# Copy the current directory contents into the container at /usr/src/app
COPY . .

# Build the application
RUN if [ "$DEV" = "true" ] ; then cargo build ; else cargo build --release ; fi

# Run the application
CMD if [ "$DEV" = "true" ] ; then cargo run ; else cargo run --release ; fi