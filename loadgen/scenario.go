package loadgen

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/temporalio/omes/loadgen/kitchensink"
	"go.temporal.io/sdk/client"
	"go.uber.org/zap"
)

type Scenario struct {
	Description string
	Executor    Executor
}

// Executor for a scenario.
type Executor interface {
	// Run the scenario
	Run(context.Context, ScenarioInfo) error
}

// ExecutorFunc is an [Executor] implementation for a function
type ExecutorFunc func(context.Context, ScenarioInfo) error

// Run implements [Executor.Run].
func (e ExecutorFunc) Run(ctx context.Context, info ScenarioInfo) error { return e(ctx, info) }

// HasDefaultConfiguration is an interface executors can implement to show their
// default configuration.
type HasDefaultConfiguration interface {
	GetDefaultConfiguration() RunConfiguration
}

var registeredScenarios = make(map[string]*Scenario)

// MustRegisterScenario registers a scenario in the global static registry.
// Panics if registration fails.
// The file name of the caller is be used as the scenario name.
func MustRegisterScenario(scenario Scenario) {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		panic("Could not infer caller when registering a nameless scenario")
	}
	scenarioName := strings.Replace(filepath.Base(file), ".go", "", 1)
	_, found := registeredScenarios[scenarioName]
	if found {
		panic(fmt.Errorf("duplicate scenario with name: %s", scenarioName))
	}
	registeredScenarios[scenarioName] = &scenario
}

// GetScenarios gets a copy of registered scenarios
func GetScenarios() map[string]*Scenario {
	ret := make(map[string]*Scenario, len(registeredScenarios))
	for k, v := range registeredScenarios {
		ret[k] = v
	}
	return ret
}

// GetScenario gets a scenario by name from the global static registry.
func GetScenario(name string) *Scenario {
	return registeredScenarios[name]
}

// ScenarioInfo contains information about the scenario under execution.
type ScenarioInfo struct {
	// Name of the scenario (inferred from the file name)
	ScenarioName string
	// Run ID of the current scenario run, used to generate a unique task queue
	// and workflow ID prefix. This is a single value for the whole scenario, and
	// not a Workflow RunID.
	RunID string
	// Metrics component for registering new metrics.
	MetricsHandler client.MetricsHandler
	// A zap logger.
	Logger *zap.SugaredLogger
	// A Temporal client.
	Client client.Client
	// Configuration info passed by user if any.
	Configuration RunConfiguration
	// ScenarioOptions are info passed from the command line. Do not mutate these.
	ScenarioOptions map[string]string
	// The namespace that was used when connecting the client.
	Namespace string
}

func (s *ScenarioInfo) ScenarioOptionInt(name string, defaultValue int) int {
	v := s.ScenarioOptions[name]
	if v == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		panic(err)
	}
	return i
}

func (s *ScenarioInfo) ScenarioOptionDuration(name string, defaultValue time.Duration) time.Duration {
	v := s.ScenarioOptions[name]
	if v == "" {
		return defaultValue
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		panic(err)
	}
	return d
}

func (s *ScenarioInfo) TaskQueue() string {
	return TaskQueueForRun(s.ScenarioName, s.RunID)
}

const DefaultIterations = 10
const DefaultMaxConcurrent = 10

type RunConfiguration struct {
	// Number of iterations to run of this scenario (mutually exclusive with Duration).
	Iterations int
	// Duration limit of this scenario (mutually exclusive with Iterations). If
	// neither iterations nor duration is set, default is DefaultIterations.
	Duration time.Duration
	// Maximum number of instances of the Execute method to run concurrently.
	// Default is DefaultMaxConcurrent.
	MaxConcurrent int
}

func (r *RunConfiguration) ApplyDefaults() {
	if r.Iterations == 0 && r.Duration == 0 {
		r.Iterations = DefaultIterations
	}
	if r.MaxConcurrent == 0 {
		r.MaxConcurrent = DefaultMaxConcurrent
	}
}

// Run represents an individual scenario run (many may be in a single instance (of possibly many) of a scenario).
type Run struct {
	// Do not mutate this, this is shared across the entire scenario
	*ScenarioInfo
	// Each run should have a unique iteration.
	Iteration int
	Logger    *zap.SugaredLogger
}

// NewRun creates a new run.
func (s *ScenarioInfo) NewRun(iteration int) *Run {
	return &Run{
		ScenarioInfo: s,
		Iteration:    iteration,
		Logger:       s.Logger.With("iteration", iteration),
	}
}

// TaskQueueForRun returns a default task queue name for the given scenario name and run ID.
func TaskQueueForRun(scenarioName, runID string) string {
	return fmt.Sprintf("%s:%s", scenarioName, runID)
}

// DefaultStartWorkflowOptions gets default start workflow info.
func (r *Run) DefaultStartWorkflowOptions() client.StartWorkflowOptions {
	return client.StartWorkflowOptions{
		TaskQueue:                                TaskQueueForRun(r.ScenarioName, r.RunID),
		ID:                                       fmt.Sprintf("w-%s-%d", r.RunID, r.Iteration),
		WorkflowExecutionErrorWhenAlreadyStarted: true,
	}
}

// DefaultKitchenSinkWorkflowOptions gets the default kitchen sink workflow info.
func (r *Run) DefaultKitchenSinkWorkflowOptions() KitchenSinkWorkflowOptions {
	return KitchenSinkWorkflowOptions{StartOptions: r.DefaultStartWorkflowOptions()}
}

type KitchenSinkWorkflowOptions struct {
	Params       kitchensink.WorkflowParams
	StartOptions client.StartWorkflowOptions
}

// ExecuteKitchenSinkWorkflow starts the generic "kitchen sink" workflow and waits for its
// completion ignoring its result.
func (r *Run) ExecuteKitchenSinkWorkflow(ctx context.Context, options *KitchenSinkWorkflowOptions) error {
	return r.ExecuteAnyWorkflow(ctx, options.StartOptions, "kitchenSink", nil, options.Params)
}

// ExecuteAnyWorkflow wraps calls to the client executing workflows to include some logging,
// returning an error if the execution fails.
func (r *Run) ExecuteAnyWorkflow(ctx context.Context, options client.StartWorkflowOptions, workflow interface{}, valuePtr interface{}, args ...interface{}) error {
	r.Logger.Debugf("Executing workflow %s with info: %v", workflow, options)
	execution, err := r.Client.ExecuteWorkflow(ctx, options, workflow, args...)
	if err != nil {
		return err
	}
	if err := execution.Get(ctx, valuePtr); err != nil {
		return fmt.Errorf("workflow execution failed (ID: %s, run ID: %s): %w", execution.GetID(), execution.GetRunID(), err)
	}
	return nil
}
