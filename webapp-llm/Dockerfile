# Use an official Python runtime as a parent image
FROM python:3.10-slim

# Set the working directory in the container
WORKDIR /app

# Copy the current directory contents into the container
COPY . .

# Install the required packages
RUN pip install --no-cache-dir -r requirements.txt

# Make port 5000 and 6000 available to the world outside this container
EXPOSE 5000
EXPOSE 6000

# Define environment variable
ENV FLASK_APP=run.py

# Ensure the start script is executable
RUN chmod +x startup.sh

# Run the application
CMD ["/bin/sh", "startup.sh"]