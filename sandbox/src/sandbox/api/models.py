from pydantic import BaseModel, Field

class Dependency(BaseModel):
    name: str
    version_constraints: str

class CodeInput(BaseModel):
    code: str
    dependencies: list[Dependency] = Field(default_factory=list)

class CodeRunOutput(BaseModel):
    output: str
    error: str