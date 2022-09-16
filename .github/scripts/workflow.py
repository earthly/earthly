import sys
import pexpect
import os
import io

DOCKERHUB_MIRROR_USERNAME = os.environ.get('DOCKERHUB_MIRROR_USERNAME')
DOCKERHUB_MIRROR_PASSWORD = os.environ.get('DOCKERHUB_MIRROR_PASSWORD')

class DockerWorkflowRunner:
    def login(self):
        output = io.StringIO()
        print('running docker login')
        cmd = pexpect.spawn(f'docker login registry-1.docker.io.mirror.corp.earthly.dev --username "{DOCKERHUB_MIRROR_USERNAME}" --password "{ DOCKERHUB_MIRROR_PASSWORD }"', encoding='utf-8')
        cmd.logfile_read = output
        status = cmd.wait()
        print('docker finished')
        if status:
            print(f'docker login failed with exit code {status} > ')
            print('===== docker login output =====')
            s = ''.join(ch for ch in output.getvalue() if ch.isprintable() or ch == '\n')
            print(s)
            print('===== docker login output finished =====')

    def ensureSingleInstallation(self):
        # todo: make sure only docker is installed?
        pass

class PodmanWorkflowRunner:

    def ensureSingleInstallation(self):
        # todo: make sure only podman is installed?
        pass

    def login(self):
        output = io.StringIO()
        print('running podman login')
        cmd = pexpect.spawn(f'podman login registry-1.docker.io.mirror.corp.earthly.dev --username "{DOCKERHUB_MIRROR_USERNAME}" --password "{ DOCKERHUB_MIRROR_PASSWORD }"', encoding='utf-8')
        cmd.logfile_read = output
        status = cmd.wait()
        print('podman finished')
        if status:
            print(f'podman login failed with exit code {status} > ')
            print('===== podman login output =====')
            s = ''.join(ch for ch in output.getvalue() if ch.isprintable() or ch == '\n')
            print(s)
            print('===== podman login output finished =====')



if __name__ == "__main__":
    if len(sys.argv) < 3:
        raise RuntimeError("workflow.py must be called like: workflow.py {binary} {command}")
    binary: str = sys.argv[1]
    command: str = sys.argv[2]
    runner = None
    if binary.lower() == "docker":
        runner = DockerWorkflowRunner()
    elif binary.lower() == "podman":
        runner = PodmanWorkflowRunner()
    else:
        raise RuntimeError(f"binary {binary} is invalid")

    commands = {
        "login": runner.login,
        "ensureSingleInstallation": runner.ensureSingleInstallation
    }

    commandToRun = commands[command]
    if commandToRun is None:
        raise RuntimeError(f'got invalid command {commandToRun}')
