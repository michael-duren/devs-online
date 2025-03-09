package views

import "fmt"

func Layout(header, body, footer string) string {
	return fmt.Sprintf("%s\n%s\n%s", header, body, footer)
}
