package node

import (
	"io"
	"strconv"

	"github.com/cloudfoundry/packit"
	"github.com/cloudfoundry/packit/postal"
	"github.com/cloudfoundry/packit/scribe"
)

type LogEmitter struct {
	// Logger is embedded and therefore delegates all of its functions to the
	// LogEmitter.
	scribe.Logger
}

func NewLogEmitter(output io.Writer) LogEmitter {
	return LogEmitter{
		Logger: scribe.NewLogger(output),
	}
}

func (e LogEmitter) SelectedDependency(entry packit.BuildpackPlanEntry, dependency postal.Dependency) {
	source, ok := entry.Metadata["version-source"].(string)
	if !ok {
		source = "<unknown>"
	}

	e.Subprocess("Selected Node Engine version (using %s): %s", source, dependency.Version)
	e.Break()
}

func (e LogEmitter) Environment(env packit.Environment, optimizeMemory bool) {
	e.Process("Configuring environment")
	e.Subprocess("%s", scribe.NewFormattedMapFromEnvironment(env))
	e.Break()
	e.Subprocess("Writing profile.d/0_memory_available.sh")
	e.Action("Calculates available memory based on container limits at launch time.")
	e.Action("Made available in the MEMORY_AVAILABLE environment variable.")
	if optimizeMemory {
		e.Break()
		e.Subprocess("Writing profile.d/1_optimize_memory.sh")
		e.Action("Assigns the NODE_OPTIONS environment variable with flag setting to optimize memory.")
		e.Action("Limits the total size of all objects on the heap to 75%% of the MEMORY_AVAILABLE.")
	}
	e.Break()
}

func (e LogEmitter) Candidates(entries []packit.BuildpackPlanEntry) {
	e.Subprocess("Candidate version sources (in priority order):")

	var (
		sources [][2]string
		maxLen  int
	)

	for _, entry := range entries {
		versionSource, ok := entry.Metadata["version-source"].(string)
		if !ok {
			versionSource = "<unknown>"
		}

		if len(versionSource) > maxLen {
			maxLen = len(versionSource)
		}

		if entry.Version == "" {
			entry.Version = "*"
		}

		sources = append(sources, [2]string{versionSource, entry.Version})
	}

	for _, source := range sources {
		e.Action(("%-" + strconv.Itoa(maxLen) + "s -> %q"), source[0], source[1])
	}

	e.Break()
}
