package ip

import (
	"os"
	"regexp"

	"github.com/samber/oops"
)

func ReadIpFromFile(absPath string) (string, error) {
	oopsBuilder := oops.
		Code("read_ip_from_file_failed").
		With("absPath", absPath)
	// Define the regular expression pattern for matching an IP address
	ipPattern := "\\b(?:\\d{1,3}\\.){3}\\d{1,3}\\b"
	// Compile the regular expression
	ipRegex := regexp.MustCompile(ipPattern)

	// Read the contents of the file
	content, err := os.ReadFile(absPath)
	if err != nil {
		err = oopsBuilder.
			Wrapf(err, "error reading file")
		return "", err
	}

	// Find the first match in the content
	ip := ipRegex.FindString(string(content))
	if len(ip) == 0 {
		err := oopsBuilder.
			Errorf("no valid IP address found in file")
		return "", err
	}

	return ip, nil
}
