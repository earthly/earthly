#!/usr/bin/env python3
import sys
import pexpect
import os
import io

DOCKERHUB_MIRROR_USERNAME = os.environ.get('DOCKERHUB_MIRROR_USERNAME')
DOCKERHUB_MIRROR_PASSWORD = os.environ.get('DOCKERHUB_MIRROR_PASSWORD')

class FrontendCommon:
    def _login(self, binary):
        output = io.StringIO()
        print(f'running {binary} login')
        cmd = pexpect.spawn(
            f'{binary} login registry-1.docker.io.mirror.corp.earthly.dev --username "{DOCKERHUB_MIRROR_USERNAME}" --password "{DOCKERHUB_MIRROR_PASSWORD}"',
            encoding='utf-8')
        cmd.logfile_read = output
        status = cmd.wait()
        print(f'{binary} finished')
        if status:
            print(f'{binary} login failed with exit code {status} > ')
            print(f'===== {binary} login output =====')
            s = ''.join(ch for ch in output.getvalue() if ch.isprintable() or ch == '\n')
            print(s)
            print(f'===== {binary} login output finished =====')

class DockerWorkflowRunner(FrontendCommon):
    def login(self):
        super(DockerWorkflowRunner, self)._login("docker")

    def ensureSingleInstallation(self):
        # todo: make sure only docker is installed?
        pass

class PodmanWorkflowRunner(FrontendCommon):

    def ensureSingleInstallation(self):
        # todo: make sure only podman is installed?
        pass

    def login(self):
        super(PodmanWorkflowRunner, self)._login("podman")



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
