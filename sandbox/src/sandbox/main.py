import uvicorn

from sandbox.api import app

def main():
    uvicorn.run(app, host="0.0.0.0", port=9999)