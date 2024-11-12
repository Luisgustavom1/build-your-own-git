package commands

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func pktLine(command string, objID string, capabilities []string) string {
	capabilityList := strings.Join(capabilities, " ")
	line := fmt.Sprintf("%s %s %s", command, objID, capabilityList)

	totalLength := len(line) + 4
	hexLength := fmt.Sprintf("%04x", totalLength)

	return hexLength + line
}

func Clone(args []string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("usage: mygit clone <url> <some_dir>\n")
	}

	url := args[0]
	req, err := http.NewRequest(http.MethodGet, url+"/info/refs?service=git-upload-pack", nil)
	if err != nil {
		fmt.Printf("could not get packs: %s\n", err)
		os.Exit(1)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("error getting pack: %s\n", err)
		os.Exit(1)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("could not read response body: %s\n", err)
		os.Exit(1)
	}

	packs := strings.Split(string(resBody), "\n")
	refs := packs[1 : len(packs)-1] // remove first two and last one

	body := ""

	for i, ref := range refs {
		parts := strings.Split(ref, " ")[0] // size+hash
		hash := parts[len(parts)-40:]       // remove size

		capabilities := []string{}
		if i == 0 {
			capabilities = append(capabilities, "multi_ack_detailed",
				"side-band-64k",
				"agent=git/2.43.0")
		}

		body += pktLine("want", hash, capabilities) + "\n"
	}

	body += "0009done\n"
	body += "0000"

	req, err = http.NewRequest(http.MethodPost, url+"/git-upload-pack", strings.NewReader(body))
	if err != nil {
		fmt.Printf("could not perform git-upload-pack: %s\n", err)
		os.Exit(1)
	}

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("error performing git-upload-pack: %s\n", err)
		os.Exit(1)
	}

	resBody, err = io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("could not read response body: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(string(resBody))

	dir := args[1]
	os.Mkdir(dir, 0755)

	return "", nil
}
