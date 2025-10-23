import os

from fastapi.param_functions import Header
from fastapi.exceptions import HTTPException

internal_key = os.getenv("SANDBOX_API_KEY")

async def check_api_key(x_api_key: str = Header(None)):
    if x_api_key == internal_key:
        return x_api_key
    else:
        raise HTTPException(status_code=401, detail="Invalid API key")
