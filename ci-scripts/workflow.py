#!/usr/bin/env python3
import sys
import pexpect
import os
import io

DOCKERHUB_MIRROR_USERNAME = os.environ.get('DOCKERHUB_MIRROR_USERNAME')
DOCKERHUB_MIRROR_PASSWORD = os.environ.get('DOCKERHUB_MIRROR_PASSWORD')


def run(command_to_run: str, cmd_name: str, **kwargs) -> (int, str):
    print(f'running {cmd_name}', flush=True)
    output = io.StringIO()
    cmd = pexpect.spawn(command_to_run, **kwargs)
    cmd.logfile_read = output
    status = cmd.wait()
    print(f'{cmd_name} finished', flush=True)
    s = ''.join(ch for ch in output.getvalue() if ch.isprintable() or ch == '\n')
    if status:
        print(f'{cmd_name} failed with exit code {status} > ', flush=True)
    else:
        print(f'{cmd_name} exit code 0', flush=True)
    print(f'===== {cmd_name} login output =====', flush=True)
    print(s)
    print(f'===== {cmd_name} login output finished =====', flush=True)
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
        # Uninstall podman, Assuming Ubuntu
        status, output = run("apt-get purge podman -y", "apt-get purge podman -y")
        if status == 0:
            status, output = run("podman --version", "podman --version")
            if status:
                # podman uninstalled successfully
                return
            raise RuntimeError("still detected Podman after purge command")

        raise RuntimeError(f"failed to uninstall Podman > {status} > {output}")


class PodmanWorkflowRunner(FrontendCommon):

    def ensure_single_install(self):
        status, output = run("docker --version", "docker --version")
        if status:
            # Assume Docker is NOT installed and return
            return
        # Uninstall docker completely
        status, output = run("apt-get purge -y docker-engine docker docker.io docker-ce docker-ce-cli",
                             "apt-get purge -y docker-engine docker docker.io docker-ce docker-ce-cli")
        if status:
            raise RuntimeError(f"failed to uninstall docker first step > {status} > {output}")
        status, output = run("apt-get autoremove -y --purge docker-engine docker docker.io docker-ce",
                             "apt-get autoremove -y --purge docker-engine docker docker.io docker-ce")
        if status:
            raise RuntimeError(f"failed to uninstall docker second step > {status} > {output}")
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
        "ensureSingleInstallation": runner.ensure_single_install
    }

    command_to_run = commands[command]
    if command_to_run is None:
        raise RuntimeError(f'got invalid command {command_to_run}')
