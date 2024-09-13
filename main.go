/*
Copyright © 2023 Syndis ehf. <syndis@syndis.is>
*/
package main

import (
	"context"
	"net/http"

	"github.com/syndis-software/aftra-cli/cmd"
)

func main() {
	cmd.Execute(context.Background(), &http.Client{})
}
