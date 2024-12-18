# Use an official Python runtime as a parent image
FROM python:3.9-slim

# Set the working directory in the container
WORKDIR /app

# Install required Python packages
RUN pip install --no-cache-dir fastapi fastembed uvicorn

# Copy the current directory contents into the container at /app
COPY . .

# Make port 8000 available to the world outside this container
EXPOSE 8000

# Define environment variable
ENV PORT=8000

# Run embeddings.py when the container launches
CMD ["python", "embeddings.py"]

