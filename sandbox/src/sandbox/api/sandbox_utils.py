import os
import base64

from secrets import token_urlsafe
from open_sandboxes.uv_config.config_pyproject import PyprojectConfig, PyprojectDependency
from open_sandboxes.sandbox import Sandbox
from .models import Dependency

sandbox_user = os.getenv("SANDBOX_USER", "")
sandbox_password = os.getenv("SANDBOX_PASSPHRASE", "")
sandbox_ssh_port = int(os.getenv("SANDBOX_SSH_PORT", "22"))
sandbox_host = os.getenv("SANDBOX_HOST", "")
sandbox_key_file = os.getenv("SANDBOX_KEY_FILE", "")
sandbox_key = os.getenv("SANDBOX_PRIVATE_KEY", "")

def sandbox_key_to_file() -> None:
    key = base64.b64decode(sandbox_key)
    os.makedirs(os.path.dirname(sandbox_key_file), exist_ok=True)
    with open(sandbox_key_file, "wb") as f:
        f.write(key)

def dependencies_to_pyproject_config(dependencies: list[Dependency]) -> PyprojectConfig:
    deps: list[PyprojectDependency] = []
    for dependency in dependencies:
        deps.append({"name": dependency.name, "version_constraints": dependency.version_constraints})
    return PyprojectConfig(
        dependencies=deps,
    )

def setup_sandbox(dependencies: list[Dependency]) -> Sandbox:
    sandbox_key_to_file()
    config = dependencies_to_pyproject_config(dependencies)
    return Sandbox.from_connection_args(
        host=sandbox_host,
        name="sandbox-" + token_urlsafe(16),
        port=sandbox_ssh_port,
        username=sandbox_user,
        passphrase=sandbox_password,
        key_file=sandbox_key_file,
        config=config
    )
