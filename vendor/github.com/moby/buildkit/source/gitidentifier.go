package source

import (
	"net/url"
	"path"
	"strings"

	srctypes "github.com/moby/buildkit/source/types"
	"github.com/moby/buildkit/util/sshutil"
	"github.com/pkg/errors"
)

type GitIdentifier struct {
	Remote           string
	Ref              string
	Subdir           string
	KeepGitDir       bool
	AuthTokenSecret  string
	AuthHeaderSecret string
	MountSSHSock     string
	KnownSSHHosts    string
}

func NewGitIdentifier(remoteURL string) (*GitIdentifier, error) {
	repo := GitIdentifier{}

	if !isGitTransport(remoteURL) {
		remoteURL = "https://" + remoteURL
	}

	var fragment string
	if sshutil.IsImplicitSSHTransport(remoteURL) {
		// implicit ssh urls such as "git@.." are not actually a URL, so cannot be parsed as URL
		parts := strings.SplitN(remoteURL, "#", 2)

		repo.Remote = parts[0]
		if len(parts) == 2 {
			fragment = parts[1]
		}
		repo.Ref, repo.Subdir = getRefAndSubdir(fragment)
	} else {
		u, err := url.Parse(remoteURL)
		if err != nil {
			return nil, err
		}

		repo.Ref, repo.Subdir = getRefAndSubdir(u.Fragment)
		u.Fragment = ""
		repo.Remote = u.String()
	}
	return &repo, nil
}

func (i *GitIdentifier) ID() string {
	return srctypes.GitScheme
}

// isGitTransport returns true if the provided str is a git transport by inspecting
// the prefix of the string for known protocols used in git.
func isGitTransport(str string) bool {
	return strings.HasPrefix(str, "http://") || strings.HasPrefix(str, "https://") || strings.HasPrefix(str, "git://") || strings.HasPrefix(str, "ssh://") || sshutil.IsImplicitSSHTransport(str)
}

// validateFullRemoteURL requires fullURL to resolve to the same host and
// repository path as originalRemote (an already-normalized GitIdentifier.Remote),
// returning fullURL unchanged if so. This prevents the git.fullurl attribute
// from retargeting a build step's actual clone URL to a different repository
// than the one its source identifier - and thus any source policy matched
// against that identifier - describes.
func validateFullRemoteURL(originalRemote, fullURL string) (string, error) {
	if !isGitTransport(fullURL) {
		fullURL = "https://" + fullURL
	}
	fullHost, fullPath, err := gitRemoteHostPath(fullURL)
	if err != nil {
		return "", errors.New("failed to parse git.fullurl")
	}

	origHost, origPath, err := gitRemoteHostPath(originalRemote)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse Git source identifier")
	}

	if !strings.EqualFold(origHost, fullHost) || gitRepositoryPath(origPath) != gitRepositoryPath(fullPath) {
		return "", errors.New("git.fullurl does not match source identifier")
	}
	return fullURL, nil
}

// gitRepositoryPath canonicalizes a URL path for repository-identity
// comparison, matching the canonicalization NewGitIdentifier/ParseGitRef apply
// to identifiers (treating "/user/repo" and "/user/repo.git" as equivalent).
func gitRepositoryPath(p string) string {
	return strings.TrimSuffix(path.Join("/", p), ".git")
}

// gitRemoteHostPath extracts the host and path from a git remote URL, which
// may be an implicit SCP-style URL (user@host:path) that net/url.Parse cannot
// handle directly.
func gitRemoteHostPath(remote string) (host, repoPath string, err error) {
	if sshutil.IsImplicitSSHTransport(remote) {
		rest := remote
		if at := strings.Index(remote, "@"); at >= 0 {
			rest = remote[at+1:]
		}
		colon := strings.Index(rest, ":")
		if colon < 0 {
			return "", "", errors.Errorf("invalid scp-style git URL %q", remote)
		}
		return rest[:colon], rest[colon+1:], nil
	}
	u, err := url.Parse(remote)
	if err != nil {
		return "", "", err
	}
	return u.Host, u.Path, nil
}

func getRefAndSubdir(fragment string) (ref string, subdir string) {
	refAndDir := strings.SplitN(fragment, ":", 2)
	ref = ""
	if len(refAndDir[0]) != 0 {
		ref = refAndDir[0]
	}
	if len(refAndDir) > 1 && len(refAndDir[1]) != 0 {
		subdir = refAndDir[1]
	}
	subdir = path.Join("/", subdir)
	subdir = strings.TrimPrefix(subdir, "/")
	return
}
