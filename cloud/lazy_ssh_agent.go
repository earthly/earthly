package cloud

import (
	"net"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// ErrNoSSHAgent occurs when no ssh auth agent exists
var ErrNoSSHAgent = errors.Errorf("no ssh auth agent socket")

type lazySSHAgent struct {
	sockPath string
	sshAgent agent.ExtendedAgent
}

func (lsa *lazySSHAgent) maybeInit() error {
	if lsa.sshAgent != nil {
		return nil
	}
	if lsa.sockPath == "" {
		return ErrNoSSHAgent
	}
	agentSock, err := net.Dial("unix", lsa.sockPath)
	if err != nil {
		return errors.Wrap(err, "failed to connect to ssh-agent")
	}

	lsa.sshAgent = agent.NewClient(agentSock)
	return nil
}

func (lsa *lazySSHAgent) List() ([]*agent.Key, error) {
	err := lsa.maybeInit()
	if err != nil {
		return nil, err
	}
	return lsa.sshAgent.List()
}

func (lsa *lazySSHAgent) Sign(key ssh.PublicKey, data []byte) (*ssh.Signature, error) {
	err := lsa.maybeInit()
	if err != nil {
		return nil, err
	}
	return lsa.sshAgent.Sign(key, data)
}

func (lsa *lazySSHAgent) Add(key agent.AddedKey) error {
	err := lsa.maybeInit()
	if err != nil {
		return err
	}
	return lsa.sshAgent.Add(key)
}

func (lsa *lazySSHAgent) Remove(key ssh.PublicKey) error {
	err := lsa.maybeInit()
	if err != nil {
		return err
	}
	return lsa.sshAgent.Remove(key)
}

func (lsa *lazySSHAgent) RemoveAll() error {
	err := lsa.maybeInit()
	if err != nil {
		return err
	}
	return lsa.sshAgent.RemoveAll()
}

func (lsa *lazySSHAgent) Lock(passphrase []byte) error {
	err := lsa.maybeInit()
	if err != nil {
		return err
	}
	return lsa.sshAgent.Lock(passphrase)
}

func (lsa *lazySSHAgent) Unlock(passphrase []byte) error {
	err := lsa.maybeInit()
	if err != nil {
		return err
	}
	return lsa.sshAgent.Unlock(passphrase)
}

func (lsa *lazySSHAgent) Signers() ([]ssh.Signer, error) {
	err := lsa.maybeInit()
	if err != nil {
		return nil, err
	}
	return lsa.sshAgent.Signers()
}

func (lsa *lazySSHAgent) SignWithFlags(key ssh.PublicKey, data []byte, flags agent.SignatureFlags) (*ssh.Signature, error) {
	err := lsa.maybeInit()
	if err != nil {
		return nil, err
	}
	return lsa.sshAgent.SignWithFlags(key, data, flags)
}

func (lsa *lazySSHAgent) Extension(extensionType string, contents []byte) ([]byte, error) {
	err := lsa.maybeInit()
	if err != nil {
		return nil, err
	}
	return lsa.sshAgent.Extension(extensionType, contents)
}
