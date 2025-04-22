package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/glsubri/kex/internal/tui"
	"github.com/go-openapi/spec"
)

func loadOpenAPISpecFromKubectl() (*spec.Swagger, error) {
	cmd := exec.Command("kubectl", "get", "--raw", "/openapi/v2")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var swagger spec.Swagger
	if err := json.Unmarshal(output, &swagger); err != nil {
		return nil, err
	}

	return &swagger, nil
}

func main() {
	spec, err := loadOpenAPISpecFromKubectl()
	if err != nil {
		panic(err)
	}

	m := tui.NewModel(spec)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
