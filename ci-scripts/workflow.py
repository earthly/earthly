#!/usr/bin/env python3
import sys
import pexpect
import os
import io

DOCKERHUB_MIRROR_USERNAME = os.environ.get('DOCKERHUB_MIRROR_USERNAME')
DOCKERHUB_MIRROR_PASSWORD = os.environ.get('DOCKERHUB_MIRROR_PASSWORD')


def run(command_to_run: str, cmd_name: str, **kwargs) -> (int, str):
    print(f'running {cmd_name}')
    output = io.StringIO()
    cmd = pexpect.spawn(command_to_run, **kwargs)
    cmd.logfile = output
    cmd.expect(pexpect.EOF)
    s = output.getvalue()
    status = cmd.wait()
    print(f'===== {cmd_name} output =====')
    print(f'===== {cmd_name} output =====\n', s, f'\n===== {cmd_name} output finished =====')
    if status:
        print(f'failed with exit code {status} > ')
    else:
        print(f'exit code 0')
    return status, s


class FrontendCommon:
    @staticmethod
    def _login(binary):
        run(
            f'{binary} login registry-1.docker.io.mirror.corp.earthly.dev --username "{DOCKERHUB_MIRROR_USERNAME}" --password "{DOCKERHUB_MIRROR_PASSWORD}"',
            f'{binary} login',
            encoding='utf-8')


class DockerWorkflowRunner(FrontendCommon):
    def login(self):
        FrontendCommon._login("docker")

    def ensure_single_install(self):
        status, output = run("podman --version", "podman --version")
        if status:
            # Assume podman is NOT installed and return
            return
        # Uninstall podman, Assuming Ubuntu, ignore error because it may not be installed
        run("sudo apt-get purge podman -y", "apt-get purge podman -y")
        status, output = run("podman --version", "podman --version")
        if status:
            # podman uninstalled successfully
            return
        raise RuntimeError("still detected Podman after purge command")


class PodmanWorkflowRunner(FrontendCommon):

    def ensure_single_install(self):
        status, output = run("docker --version", "docker --version")
        if status:
            # Assume Docker is NOT installed and return
            return
        # Uninstall docker completely, ignore errors because some stuff may not be installed
        for uninstall in ["docker-engine", "docker", "docker.io", "docker-ce", "docker-ce-cli"]:
            run(f"sudo apt-get autoremove purge -y {uninstall}", f"apt-get autoremove purge -y {uninstall}")
        status, output = run("docker --version", "docker --version")
        if status:
            # docker uninstalled successfully
            return
        raise RuntimeError(f"docker still detected after uninstall commands > {status} > {output}")

    def login(self):
        FrontendCommon._login("podman")


if __name__ == "__main__":
    if len(sys.argv) < 3:
        raise RuntimeError("workflow.py must be called like: workflow.py {binary} {command}")
    binary: str = sys.argv[1]
    command: str = sys.argv[2]
    if binary.lower() == "docker":
        runner = DockerWorkflowRunner()
    elif binary.lower() == "podman":
        runner = PodmanWorkflowRunner()
    else:
        raise RuntimeError(f"binary {binary} is invalid")

    commands = {
        "login": runner.login,
        "ensure_single_install": runner.ensure_single_install
    }

    if command not in commands:
        raise RuntimeError(f"invalid command given {command}, must be one of {commands.keys()}")
    commands[command]()
