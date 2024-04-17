package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"text/tabwriter"
	"time"
)

type TestEvent struct {
	Time    time.Time // encodes as an RFC3339-format string
	Action  string
	Package string
	Test    string
	Elapsed float64 // seconds
	Output  string
}

func main() {
	eventsWithElapsedTimes := []TestEvent{}
	scanner := bufio.NewScanner(os.Stdin)
	passed := true
	for scanner.Scan() {
		var event TestEvent
		l := scanner.Text()
		if err := json.Unmarshal([]byte(l), &event); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		fmt.Printf("%s", event.Output)
		if event.Elapsed > 0 {
			eventsWithElapsedTimes = append(eventsWithElapsedTimes, event)
		}
		if event.Action == "fail" {
			passed = false
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	sort.Slice(eventsWithElapsedTimes, func(i, j int) bool {
		return eventsWithElapsedTimes[i].Elapsed < eventsWithElapsedTimes[j].Elapsed
	})

	fmt.Printf("\n--- Test Duration Summary ---\n")

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Package\tTest\tAction\tElapsed (seconds)\n")
	for _, event := range eventsWithElapsedTimes {
		if event.Test != "" {
			fmt.Fprintf(w, "%s\t%s\t%s\t%v\n", event.Package, event.Test, event.Action, event.Elapsed)
		}
	}
	w.Flush()
	fmt.Printf("%s", buf.String())

	if !passed {
		fmt.Printf("test(s) failed\n")
		os.Exit(1)
	}
	fmt.Printf("test(s) passed\n")
}
