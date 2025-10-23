from fastapi.applications import FastAPI
from fastapi.responses import ORJSONResponse
from fastapi.param_functions import Depends
from fastapi_limiter.depends import RateLimiter
from fastapi.exceptions import HTTPException

from .api_key_utils import check_api_key
from .limiter import lifespan
from .models import CodeInput, CodeRunOutput
from .sandbox_utils import setup_sandbox

app = FastAPI(default_response_class=ORJSONResponse, lifespan=lifespan)

@app.post("/run", dependencies=[Depends(RateLimiter(times=10, seconds=60))])
async def run_code_in_sandbox(input: CodeInput, x_api_key: str = Depends(check_api_key)) -> CodeRunOutput:
    code_to_run = input.code
    dependencies = input.dependencies
    s = setup_sandbox(dependencies=dependencies)
    try:
        result = s.run_code(code=code_to_run)
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Error while provisioning the sandbox and running the code: {e}")
    else:
        return CodeRunOutput(
            output=result["output"],
            error=result["error"]
        )