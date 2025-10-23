import redis.asyncio as redis
import os

from fastapi.applications import FastAPI
from fastapi_limiter import FastAPILimiter
from contextlib import asynccontextmanager

@asynccontextmanager
async def lifespan(_: FastAPI):
    redis_connection = redis.from_url(os.getenv("REDIS_URL"))
    await FastAPILimiter.init(redis_connection)
    yield
    await FastAPILimiter.close()


