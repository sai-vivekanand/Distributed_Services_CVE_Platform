import logging
import json
from flask import Flask, request

def create_json_logger():
    logger = logging.getLogger("flask_app")
    logger.setLevel(logging.INFO)
    handler = logging.StreamHandler()
    formatter = logging.Formatter(json.dumps({
        "timestamp": "%(asctime)s",
        "level": "%(levelname)s",
        "message": "%(message)s",
        "module": "%(module)s",
        "method": "%(funcName)s",
        "line": "%(lineno)d"
    }))
    handler.setFormatter(formatter)
    logger.addHandler(handler)
    return logger

def create_app():
    app = Flask(__name__)

    # Disable the default Flask logging
    log = logging.getLogger('werkzeug')
    log.disabled = True
    app.logger.disabled = True

    # Set up the JSON logger
    logger = create_json_logger()

    # Register blueprints
    from .routes import main
    app.register_blueprint(main)

    # Add a before request hook to log incoming requests
    @app.before_request
    def log_request():
        logger.info(f"Received request: {request.method} {request.url} with data: {request.json}")

    # Add an after request hook to log the responses
    @app.after_request
    def log_response(response):
        logger.info(f"Response status: {response.status} with data: {response.get_json()}")
        return response

    # Attach the logger to the app
    app.logger = logger

    return app
