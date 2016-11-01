package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sync"
	"time"
)

var remotes = []string{}

// Create wait group to make sure we can track when all is done
var wg sync.WaitGroup

// For checking if a string is in a slice
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func main() {
	fmt.Println("\n============GIT PARALLEL PUSH=============")
	var ignore = []string{}

	// Eventually use this to make sure we're in the root directory
	// git rev-parse --show-toplevel
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")

	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	out, outErr := cmd.Output()

	if outErr != nil {
		fmt.Printf(outErr.Error())
	}
	os.Chdir(string(out))

	// Set up the ignores array for git remotes to ignore
	fmt.Println("\nDeciding which remotes to ignore")
	fmt.Println("===============================================\n")
	gitIgnore, _ := os.Open(".gitp2ignore")
	defer gitIgnore.Close()
	scanner1 := bufio.NewScanner(gitIgnore)
	scanner1.Split(bufio.ScanLines)

	for scanner1.Scan() {
		fmt.Println("The " + scanner1.Text() + " remote will be ignored.")
		// Add the name of this remote to the list of remotes to ignore
		ignore = append(ignore, scanner1.Text())
	}

	// Set up the remotes array for git remotes to push to
	gitConfig, _ := os.Open(".git/config")
	defer gitConfig.Close()
	scanner2 := bufio.NewScanner(gitConfig)
	scanner2.Split(bufio.ScanLines)

	var remoteString = regexp.MustCompile(`\[(.*)remote ("|')(.*)('|")(.*)\]`)

	for scanner2.Scan() {
		if remoteString.MatchString(scanner2.Text()) {
			x := remoteString.FindAllStringSubmatch(scanner2.Text(), -1)
			remoteName := x[0][3]
			if stringInSlice(remoteName, ignore) == false {
				remotes = append(remotes, remoteName)
			}
		}
	}

	// Need to make this use go routines
	for _, remote := range remotes {
		// Increment the WaitGroup counter.
		wg.Add(1)
		go git_push(remote)
	}

	//Show output as running
	ticker := time.NewTicker(time.Millisecond * 500)
	go func() {
		for range ticker.C {
			fmt.Printf(".")
		}
	}()

	// Wait for all git pushes to complete
	wg.Wait()
	ticker.Stop()
	fmt.Println("\n\nAll done.")
}

func git_push(remote string) {
	// Decrement the counter when the goroutine completes.
	defer wg.Done()
	cmd := exec.Command("git", "push", remote, "master")

	out, outErr := cmd.CombinedOutput()

	fmt.Printf("\n\nOutput for pushing to " + remote)
	fmt.Printf("\n=================================\n")

	if outErr != nil {
		fmt.Printf(outErr.Error() + "\n")
	} else {
		fmt.Printf(string(out) + "\n")
	}
}
