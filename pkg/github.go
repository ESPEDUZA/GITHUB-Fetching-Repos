package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"sort"
	"strings"
)

type Repository struct {
	Name        string `json:"name"`
	CloneURL    string `json:"clone_url"`
	Description string `json:"description"`
	UpdatedAt   string `json:"updated_at"`
}

func FetchRepositories(user string, token string) ([]Repository, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/users/%s/repos?sort=updated&per_page=100", user), nil)
	if err != nil {
		return nil, err
	}

	if token != "" {
		req.Header.Add("Authorization", "token "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var repos []Repository
	err = json.Unmarshal(body, &repos)
	if err != nil {
		return nil, err
	}

	// Trier les répos par date de dernière modification
	sort.Slice(repos, func(i, j int) bool {
		return repos[i].UpdatedAt > repos[j].UpdatedAt
	})

	return repos, nil
}

func ExecuteGitCommands(repoDir string) error {
	// Exécuter git fetch pour récupérer toutes les références de branches
	err := executeGitCommand(repoDir, "fetch", "--all")
	if err != nil {
		return fmt.Errorf("erreur lors de l'exécution de git fetch: %w", err)
	}

	// Obtenir le nom de la branche avec le dernier commit
	branch, err := getLastModifiedBranch(repoDir)
	if err != nil {
		return fmt.Errorf("erreur lors de l'obtention de la dernière branche modifiée: %w", err)
	}

	// Exécuter git pull sur la branche avec le dernier commit
	err = executeGitCommand(repoDir, "pull", "origin", branch)
	if err != nil {
		return fmt.Errorf("erreur lors de l'exécution de git pull sur la branche %s: %w", branch, err)
	}
	fmt.Printf("%s a bien été fetch  + pull\n", repoDir)

	return nil
}

func executeGitCommand(repoDir string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = repoDir
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func getLastModifiedBranch(repoDir string) (string, error) {
	cmd := exec.Command("git", "for-each-ref", "--sort=-committerdate", "--count=1", "--format=%(refname:short)", "refs/heads/")
	cmd.Dir = repoDir
	branchBytes, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(branchBytes)), nil
}
