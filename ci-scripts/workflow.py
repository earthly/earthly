#!/usr/bin/env python3
import sys
import os

DOCKERHUB_MIRROR_USERNAME = os.environ.get('DOCKERHUB_MIRROR_USERNAME')
DOCKERHUB_MIRROR_PASSWORD = os.environ.get('DOCKERHUB_MIRROR_PASSWORD')


def run(command_to_run: str, cmd_name: str) -> int:
    print(f'running {cmd_name}', flush=True)
    print(f'===== {cmd_name} output =====', flush=True)
    status = os.system(command_to_run)
    sys.stdout.flush()
    sys.stderr.flush()
    print(f'\n===== {cmd_name} output finished =====', flush=True)
    if status:
        print(f'failed with exit code {status} > ', flush=True)
    else:
        print(f'exit code 0', flush=True)
    return status


class FrontendCommon:
    @staticmethod
    def _login(binary):
        run(
            f'{binary} login registry-1.docker.io.mirror.corp.earthly.dev --username "{DOCKERHUB_MIRROR_USERNAME}" --password "{DOCKERHUB_MIRROR_PASSWORD}"',
            f'{binary} login')


class DockerWorkflowRunner(FrontendCommon):
    def login(self):
        FrontendCommon._login("docker")

    def ensure_single_install(self):
        pass
        # status = run("podman --version", "podman --version")
        # if status:
        #     # Assume podman is NOT installed and return
        #     print("podman may already be removed...", flush=True)
        #     return
        # # Uninstall podman, Assuming Ubuntu, ignore error because it may not be installed
        # run("apt-get purge podman -y", "apt-get purge podman -y")
        # status = run("podman --version", "podman --version")
        # if status:
        #     # podman uninstalled successfully
        #     return
        # run("which podman", "which podman")
        # raise RuntimeError("still detected Podman after purge command")


class PodmanWorkflowRunner(FrontendCommon):

    def ensure_single_install(self):
        status = run("docker --version", "docker --version")
        if status:
            # Assume Docker is NOT installed and return
            return
        # Uninstall docker completely, ignore errors because some stuff may not be installed
        for uninstall in ["docker-engine", "docker", "docker.io", "docker-ce", "docker-ce-cli"]:
            run(f"apt-get purge -y {uninstall}", f"apt-get purge -y {uninstall}")
            run(f"apt-get autoremove -y --purge {uninstall}", f"apt-get autoremove -y --purge {uninstall}")
        status = run("docker --version", "docker --version")
        if status:
            # docker uninstalled successfully
            return
        run("which docker", "which docker")
        run("rm -rf /usr/bin/docker", "rm -rf /usr/bin/docker")
        run("which docker", "which docker")
        run("credential-helper-docker", "which docker")
        raise RuntimeError(f"docker still detected after uninstall commands > {status}")

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
