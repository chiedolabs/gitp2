package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
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
	/////////////////////////////////////////////////
	// PREPARE SOME OUTPUT COLORS
	/////////////////////////////////////////////////
	boldGreen := color.New(color.FgGreen, color.Bold)
	boldWhite := color.New(color.FgWhite, color.Bold)
	//boldBlue := color.New(color.FgBlue, color.Bold)
	//boldRed := color.New(color.FgRed, color.Bold)

	start := time.Now()
	boldWhite.Println("\n============GIT PARALLEL PUSH=============")
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
	gitIgnore, _ := os.Open(".gitp2ignore")
	defer gitIgnore.Close()
	scanner1 := bufio.NewScanner(gitIgnore)
	scanner1.Split(bufio.ScanLines)

	for scanner1.Scan() {
		fmt.Println("The " + scanner1.Text() + " remote will be ignored.")
		// Add the name of this remote to the list of remotes to ignore
		ignore = append(ignore, scanner1.Text())
	}

	// Display if any remotes will be ignored
	if len(ignore) == 0 {
		fmt.Println("\nNo remotes will be ignored")
	}

	// Allow the overidding of the git config file if desired
	configFile := ".git/config"
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	// Set up the remotes array for git remotes to push to
	gitConfig, _ := os.Open(configFile)
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

	// Wait for all git pushes to complete
	wg.Wait()
	ticker.Stop()
	elapsed := time.Since(start)
	boldGreen.Printf("\n\nYour scripts took %f seconds\n\n", elapsed.Seconds())
}

func git_push(remote string) {
	// Decrement the counter when the goroutine completes.
	defer wg.Done()
	cmd := exec.Command("git", "push", remote, "master")
	// Prefix each line with the remote name and output

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	// start the command after having set up the pipe
	cmd.Start()

	// read command's stdout line by line
	outIn := bufio.NewScanner(stdout)
	errIn := bufio.NewScanner(stderr)

	for errIn.Scan() {
		fmt.Printf("\n" + remote + ": " + errIn.Text())
	}

	for outIn.Scan() {
		fmt.Printf("\n" + remote + ": " + outIn.Text())
	}

	//delete me

	cmd.Wait()

}
