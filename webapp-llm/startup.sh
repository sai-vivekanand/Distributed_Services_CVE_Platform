#!/bin/sh

# Start Flask application
flask run --host=0.0.0.0 --port=5000 &

# Set environment variables for Streamlit
export STREAMLIT_SERVER_PORT=6000
export STREAMLIT_SERVER_HEADLESS=true
export STREAMLIT_SERVER_ADDRESS=0.0.0.0

# Start Streamlit application
streamlit run /app/app/ui.py
