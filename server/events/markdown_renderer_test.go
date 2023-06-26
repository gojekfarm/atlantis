// Copyright 2017 HootSuite Media Inc.
//
// Licensed under the Apache License, Version 2.0 (the License);
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an AS IS BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// Modified hereafter by contributors to runatlantis/atlantis.

package events_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/gojekfarm/atlantis/server/events"
	"github.com/gojekfarm/atlantis/server/events/command"
	"github.com/gojekfarm/atlantis/server/events/models"
	. "github.com/gojekfarm/atlantis/testing"
)

func TestRenderErr(t *testing.T) {
	err := errors.New("err")
	cases := []struct {
		Description string
		Command     command.Name
		Error       error
		Expected    string
	}{
		{
			"apply error",
			command.Apply,
			err,
			"**Apply Error**\n```\nerr\n```\n",
		},
		{
			"plan error",
			command.Plan,
			err,
			"**Plan Error**\n```\nerr\n```\n",
		},
		{
			"policy check error",
			command.PolicyCheck,
			err,
			"**Policy Check Error**\n```\nerr\n```" +
				"\n* :heavy_check_mark: To **approve** failing policies an authorized approver can comment:\n" +
				"    * `atlantis approve_policies`\n" +
				"* :repeat: Or, address the policy failure by modifying the codebase and re-planning.\n\n",
		},
	}

	r := events.MarkdownRenderer{}
	for _, c := range cases {
		res := command.Result{
			Error: c.Error,
		}
		for _, verbose := range []bool{true, false} {
			t.Run(fmt.Sprintf("%s_%t", c.Description, verbose), func(t *testing.T) {
				s := r.Render(res, c.Command, "log", verbose, models.Github)
				if !verbose {
					Equals(t, c.Expected, s)
				} else {
					Equals(t, c.Expected+"<details><summary>Log</summary>\n  <p>\n\n```\nlog```\n</p></details>\n", s)
				}
			})
		}
	}
}

func TestRenderFailure(t *testing.T) {
	cases := []struct {
		Description string
		Command     command.Name
		Failure     string
		Expected    string
	}{
		{
			"apply failure",
			command.Apply,
			"failure",
			"**Apply Failed**: failure\n",
		},
		{
			"plan failure",
			command.Plan,
			"failure",
			"**Plan Failed**: failure\n",
		},
		{
			"policy check failure",
			command.PolicyCheck,
			"failure",
			"**Policy Check Failed**: failure\n",
		},
	}

	r := events.MarkdownRenderer{}
	for _, c := range cases {
		res := command.Result{
			Failure: c.Failure,
		}
		for _, verbose := range []bool{true, false} {
			t.Run(fmt.Sprintf("%s_%t", c.Description, verbose), func(t *testing.T) {
				s := r.Render(res, c.Command, "log", verbose, models.Github)
				if !verbose {
					Equals(t, c.Expected, s)
				} else {
					Equals(t, c.Expected+"<details><summary>Log</summary>\n  <p>\n\n```\nlog```\n</p></details>\n", s)
				}
			})
		}
	}
}

func TestRenderErrAndFailure(t *testing.T) {
	r := events.MarkdownRenderer{}
	res := command.Result{
		Error:   errors.New("error"),
		Failure: "failure",
	}
	s := r.Render(res, command.Plan, "", false, models.Github)
	Equals(t, "**Plan Error**\n```\nerror\n```\n", s)
}

func TestRenderProjectResults(t *testing.T) {
	cases := []struct {
		Description    string
		Command        command.Name
		ProjectResults []command.ProjectResult
		VCSHost        models.VCSHostType
		Expected       string
	}{
		{
			"no projects",
			command.Plan,
			[]command.ProjectResult{},
			models.Github,
			"Ran Plan for 0 projects:\n\n\n\n",
		},
		{
			"single successful plan",
			command.Plan,
			[]command.ProjectResult{
				{
					PlanSuccess: &models.PlanSuccess{
						TerraformOutput: "terraform-output",
						LockURL:         "lock-url",
						RePlanCmd:       "atlantis plan -d path -w workspace",
						ApplyCmd:        "atlantis apply -d path -w workspace",
					},
					Workspace:  "workspace",
					RepoRelDir: "path",
				},
			},
			models.Github,
			`Ran Plan for dir: $path$ workspace: $workspace$

$$$diff
terraform-output
$$$

* :arrow_forward: To **apply** this plan, comment:
    * $atlantis apply -d path -w workspace$
* :put_litter_in_its_place: To **delete** this plan click [here](lock-url)
* :repeat: To **plan** this project again, comment:
    * $atlantis plan -d path -w workspace$

---
* :fast_forward: To **apply** all unapplied plans from this pull request, comment:
    * $atlantis apply$
* :put_litter_in_its_place: To delete all plans and locks for the PR, comment:
    * $atlantis unlock$
`,
		},
		{
			"single successful plan with master ahead",
			command.Plan,
			[]command.ProjectResult{
				{
					PlanSuccess: &models.PlanSuccess{
						TerraformOutput: "terraform-output",
						LockURL:         "lock-url",
						RePlanCmd:       "atlantis plan -d path -w workspace",
						ApplyCmd:        "atlantis apply -d path -w workspace",
						HasDiverged:     true,
					},
					Workspace:  "workspace",
					RepoRelDir: "path",
				},
			},
			models.Github,
			`Ran Plan for dir: $path$ workspace: $workspace$

$$$diff
terraform-output
$$$

* :arrow_forward: To **apply** this plan, comment:
    * $atlantis apply -d path -w workspace$
* :put_litter_in_its_place: To **delete** this plan click [here](lock-url)
* :repeat: To **plan** this project again, comment:
    * $atlantis plan -d path -w workspace$

:warning: The branch we're merging into is ahead, it is recommended to pull new commits first.

---
* :fast_forward: To **apply** all unapplied plans from this pull request, comment:
    * $atlantis apply$
* :put_litter_in_its_place: To delete all plans and locks for the PR, comment:
    * $atlantis unlock$
`,
		},
		{
			"single successful plan with project name",
			command.Plan,
			[]command.ProjectResult{
				{
					PlanSuccess: &models.PlanSuccess{
						TerraformOutput: "terraform-output",
						LockURL:         "lock-url",
						RePlanCmd:       "atlantis plan -d path -w workspace",
						ApplyCmd:        "atlantis apply -d path -w workspace",
					},
					Workspace:   "workspace",
					RepoRelDir:  "path",
					ProjectName: "projectname",
				},
			},
			models.Github,
			`Ran Plan for project: $projectname$ dir: $path$ workspace: $workspace$

$$$diff
terraform-output
$$$

* :arrow_forward: To **apply** this plan, comment:
    * $atlantis apply -d path -w workspace$
* :put_litter_in_its_place: To **delete** this plan click [here](lock-url)
* :repeat: To **plan** this project again, comment:
    * $atlantis plan -d path -w workspace$

---
* :fast_forward: To **apply** all unapplied plans from this pull request, comment:
    * $atlantis apply$
* :put_litter_in_its_place: To delete all plans and locks for the PR, comment:
    * $atlantis unlock$
`,
		},
		{
			"single successful policy check with project name",
			command.PolicyCheck,
			[]command.ProjectResult{
				{
					PolicyCheckSuccess: &models.PolicyCheckSuccess{
						PolicyCheckOutput: "2 tests, 1 passed, 0 warnings, 0 failure, 0 exceptions",
						LockURL:           "lock-url",
						RePlanCmd:         "atlantis plan -d path -w workspace",
						ApplyCmd:          "atlantis apply -d path -w workspace",
					},
					Workspace:   "workspace",
					RepoRelDir:  "path",
					ProjectName: "projectname",
				},
			},
			models.Github,
			`Ran Policy Check for project: $projectname$ dir: $path$ workspace: $workspace$

$$$diff
2 tests, 1 passed, 0 warnings, 0 failure, 0 exceptions
$$$

* :arrow_forward: To **apply** this plan, comment:
    * $atlantis apply -d path -w workspace$
* :put_litter_in_its_place: To **delete** this plan click [here](lock-url)
* :repeat: To re-run policies **plan** this project again by commenting:
    * $atlantis plan -d path -w workspace$

---
* :fast_forward: To **apply** all unapplied plans from this pull request, comment:
    * $atlantis apply$
* :put_litter_in_its_place: To delete all plans and locks for the PR, comment:
    * $atlantis unlock$
`,
		},
		{
			"single successful apply",
			command.Apply,
			[]command.ProjectResult{
				{
					ApplySuccess: "success",
					Workspace:    "workspace",
					RepoRelDir:   "path",
				},
			},
			models.Github,
			`Ran Apply for dir: $path$ workspace: $workspace$

$$$diff
success
$$$

`,
		},
		{
			"single successful apply with project name",
			command.Apply,
			[]command.ProjectResult{
				{
					ApplySuccess: "success",
					Workspace:    "workspace",
					RepoRelDir:   "path",
					ProjectName:  "projectname",
				},
			},
			models.Github,
			`Ran Apply for project: $projectname$ dir: $path$ workspace: $workspace$

$$$diff
success
$$$

`,
		},
		{
			"multiple successful plans",
			command.Plan,
			[]command.ProjectResult{
				{
					Workspace:  "workspace",
					RepoRelDir: "path",
					PlanSuccess: &models.PlanSuccess{
						TerraformOutput: "terraform-output",
						LockURL:         "lock-url",
						ApplyCmd:        "atlantis apply -d path -w workspace",
						RePlanCmd:       "atlantis plan -d path -w workspace",
					},
				},
				{
					Workspace:   "workspace",
					RepoRelDir:  "path2",
					ProjectName: "projectname",
					PlanSuccess: &models.PlanSuccess{
						TerraformOutput: "terraform-output2",
						LockURL:         "lock-url2",
						ApplyCmd:        "atlantis apply -d path2 -w workspace",
						RePlanCmd:       "atlantis plan -d path2 -w workspace",
					},
				},
			},
			models.Github,
			`Ran Plan for 2 projects:

1. dir: $path$ workspace: $workspace$
1. project: $projectname$ dir: $path2$ workspace: $workspace$

### 1. dir: $path$ workspace: $workspace$
$$$diff
terraform-output
$$$

* :arrow_forward: To **apply** this plan, comment:
    * $atlantis apply -d path -w workspace$
* :put_litter_in_its_place: To **delete** this plan click [here](lock-url)
* :repeat: To **plan** this project again, comment:
    * $atlantis plan -d path -w workspace$

---
### 2. project: $projectname$ dir: $path2$ workspace: $workspace$
$$$diff
terraform-output2
$$$

* :arrow_forward: To **apply** this plan, comment:
    * $atlantis apply -d path2 -w workspace$
* :put_litter_in_its_place: To **delete** this plan click [here](lock-url2)
* :repeat: To **plan** this project again, comment:
    * $atlantis plan -d path2 -w workspace$

---
* :fast_forward: To **apply** all unapplied plans from this pull request, comment:
    * $atlantis apply$
* :put_litter_in_its_place: To delete all plans and locks for the PR, comment:
    * $atlantis unlock$
`,
		},
		{
			"multiple successful policy checks",
			command.PolicyCheck,
			[]command.ProjectResult{
				{
					Workspace:  "workspace",
					RepoRelDir: "path",
					PolicyCheckSuccess: &models.PolicyCheckSuccess{
						PolicyCheckOutput: "4 tests, 4 passed, 0 warnings, 0 failures, 0 exceptions",
						LockURL:           "lock-url",
						ApplyCmd:          "atlantis apply -d path -w workspace",
						RePlanCmd:         "atlantis plan -d path -w workspace",
					},
				},
				{
					Workspace:   "workspace",
					RepoRelDir:  "path2",
					ProjectName: "projectname",
					PolicyCheckSuccess: &models.PolicyCheckSuccess{
						PolicyCheckOutput: "4 tests, 4 passed, 0 warnings, 0 failures, 0 exceptions",
						LockURL:           "lock-url2",
						ApplyCmd:          "atlantis apply -d path2 -w workspace",
						RePlanCmd:         "atlantis plan -d path2 -w workspace",
					},
				},
			},
			models.Github,
			`Ran Policy Check for 2 projects:

1. dir: $path$ workspace: $workspace$
1. project: $projectname$ dir: $path2$ workspace: $workspace$

### 1. dir: $path$ workspace: $workspace$
$$$diff
4 tests, 4 passed, 0 warnings, 0 failures, 0 exceptions
$$$

* :arrow_forward: To **apply** this plan, comment:
    * $atlantis apply -d path -w workspace$
* :put_litter_in_its_place: To **delete** this plan click [here](lock-url)
* :repeat: To re-run policies **plan** this project again by commenting:
    * $atlantis plan -d path -w workspace$

---
### 2. project: $projectname$ dir: $path2$ workspace: $workspace$
$$$diff
4 tests, 4 passed, 0 warnings, 0 failures, 0 exceptions
$$$

* :arrow_forward: To **apply** this plan, comment:
    * $atlantis apply -d path2 -w workspace$
* :put_litter_in_its_place: To **delete** this plan click [here](lock-url2)
* :repeat: To re-run policies **plan** this project again by commenting:
    * $atlantis plan -d path2 -w workspace$

---
* :fast_forward: To **apply** all unapplied plans from this pull request, comment:
    * $atlantis apply$
* :put_litter_in_its_place: To delete all plans and locks for the PR, comment:
    * $atlantis unlock$
`,
		},
		{
			"multiple successful applies",
			command.Apply,
			[]command.ProjectResult{
				{
					RepoRelDir:   "path",
					Workspace:    "workspace",
					ProjectName:  "projectname",
					ApplySuccess: "success",
				},
				{
					RepoRelDir:   "path2",
					Workspace:    "workspace",
					ApplySuccess: "success2",
				},
			},
			models.Github,
			`Ran Apply for 2 projects:

1. project: $projectname$ dir: $path$ workspace: $workspace$
1. dir: $path2$ workspace: $workspace$

### 1. project: $projectname$ dir: $path$ workspace: $workspace$
$$$diff
success
$$$

---
### 2. dir: $path2$ workspace: $workspace$
$$$diff
success2
$$$

---

`,
		},
		{
			"single errored plan",
			command.Plan,
			[]command.ProjectResult{
				{
					Error:      errors.New("error"),
					RepoRelDir: "path",
					Workspace:  "workspace",
				},
			},
			models.Github,
			`Ran Plan for dir: $path$ workspace: $workspace$

**Plan Error**
$$$
error
$$$

`,
		},
		{
			"single failed plan",
			command.Plan,
			[]command.ProjectResult{
				{
					RepoRelDir: "path",
					Workspace:  "workspace",
					Failure:    "failure",
				},
			},
			models.Github,
			`Ran Plan for dir: $path$ workspace: $workspace$

**Plan Failed**: failure

`,
		},
		{
			"successful, failed, and errored plan",
			command.Plan,
			[]command.ProjectResult{
				{
					Workspace:  "workspace",
					RepoRelDir: "path",
					PlanSuccess: &models.PlanSuccess{
						TerraformOutput: "terraform-output",
						LockURL:         "lock-url",
						ApplyCmd:        "atlantis apply -d path -w workspace",
						RePlanCmd:       "atlantis plan -d path -w workspace",
					},
				},
				{
					Workspace:  "workspace",
					RepoRelDir: "path2",
					Failure:    "failure",
				},
				{
					Workspace:   "workspace",
					RepoRelDir:  "path3",
					ProjectName: "projectname",
					Error:       errors.New("error"),
				},
			},
			models.Github,
			`Ran Plan for 3 projects:

1. dir: $path$ workspace: $workspace$
1. dir: $path2$ workspace: $workspace$
1. project: $projectname$ dir: $path3$ workspace: $workspace$

### 1. dir: $path$ workspace: $workspace$
$$$diff
terraform-output
$$$

* :arrow_forward: To **apply** this plan, comment:
    * $atlantis apply -d path -w workspace$
* :put_litter_in_its_place: To **delete** this plan click [here](lock-url)
* :repeat: To **plan** this project again, comment:
    * $atlantis plan -d path -w workspace$

---
### 2. dir: $path2$ workspace: $workspace$
**Plan Failed**: failure

---
### 3. project: $projectname$ dir: $path3$ workspace: $workspace$
**Plan Error**
$$$
error
$$$

---
* :fast_forward: To **apply** all unapplied plans from this pull request, comment:
    * $atlantis apply$
* :put_litter_in_its_place: To delete all plans and locks for the PR, comment:
    * $atlantis unlock$
`,
		},
		{
			"successful, failed, and errored policy check",
			command.PolicyCheck,
			[]command.ProjectResult{
				{
					Workspace:  "workspace",
					RepoRelDir: "path",
					PolicyCheckSuccess: &models.PolicyCheckSuccess{
						PolicyCheckOutput: "4 tests, 4 passed, 0 warnings, 0 failures, 0 exceptions",
						LockURL:           "lock-url",
						ApplyCmd:          "atlantis apply -d path -w workspace",
						RePlanCmd:         "atlantis plan -d path -w workspace",
					},
				},
				{
					Workspace:  "workspace",
					RepoRelDir: "path2",
					Failure:    "failure",
				},
				{
					Workspace:   "workspace",
					RepoRelDir:  "path3",
					ProjectName: "projectname",
					Error:       errors.New("error"),
				},
			},
			models.Github,
			`Ran Policy Check for 3 projects:

1. dir: $path$ workspace: $workspace$
1. dir: $path2$ workspace: $workspace$
1. project: $projectname$ dir: $path3$ workspace: $workspace$

### 1. dir: $path$ workspace: $workspace$
$$$diff
4 tests, 4 passed, 0 warnings, 0 failures, 0 exceptions
$$$

* :arrow_forward: To **apply** this plan, comment:
    * $atlantis apply -d path -w workspace$
* :put_litter_in_its_place: To **delete** this plan click [here](lock-url)
* :repeat: To re-run policies **plan** this project again by commenting:
    * $atlantis plan -d path -w workspace$

---
### 2. dir: $path2$ workspace: $workspace$
**Policy Check Failed**: failure

---
### 3. project: $projectname$ dir: $path3$ workspace: $workspace$
**Policy Check Error**
$$$
error
$$$
* :heavy_check_mark: To **approve** failing policies an authorized approver can comment:
    * $atlantis approve_policies$
* :repeat: Or, address the policy failure by modifying the codebase and re-planning.


---
* :fast_forward: To **apply** all unapplied plans from this pull request, comment:
    * $atlantis apply$
* :put_litter_in_its_place: To delete all plans and locks for the PR, comment:
    * $atlantis unlock$
`,
		},
		{
			"successful, failed, and errored apply",
			command.Apply,
			[]command.ProjectResult{
				{
					Workspace:    "workspace",
					RepoRelDir:   "path",
					ApplySuccess: "success",
				},
				{
					Workspace:  "workspace",
					RepoRelDir: "path2",
					Failure:    "failure",
				},
				{
					Workspace:  "workspace",
					RepoRelDir: "path3",
					Error:      errors.New("error"),
				},
			},
			models.Github,
			`Ran Apply for 3 projects:

1. dir: $path$ workspace: $workspace$
1. dir: $path2$ workspace: $workspace$
1. dir: $path3$ workspace: $workspace$

### 1. dir: $path$ workspace: $workspace$
$$$diff
success
$$$

---
### 2. dir: $path2$ workspace: $workspace$
**Apply Failed**: failure

---
### 3. dir: $path3$ workspace: $workspace$
**Apply Error**
$$$
error
$$$

---

`,
		},
		{
			"successful, failed, and errored apply",
			command.Apply,
			[]command.ProjectResult{
				{
					Workspace:    "workspace",
					RepoRelDir:   "path",
					ApplySuccess: "success",
				},
				{
					Workspace:  "workspace",
					RepoRelDir: "path2",
					Failure:    "failure",
				},
				{
					Workspace:  "workspace",
					RepoRelDir: "path3",
					Error:      errors.New("error"),
				},
			},
			models.Github,
			`Ran Apply for 3 projects:

1. dir: $path$ workspace: $workspace$
1. dir: $path2$ workspace: $workspace$
1. dir: $path3$ workspace: $workspace$

### 1. dir: $path$ workspace: $workspace$
$$$diff
success
$$$

---
### 2. dir: $path2$ workspace: $workspace$
**Apply Failed**: failure

---
### 3. dir: $path3$ workspace: $workspace$
**Apply Error**
$$$
error
$$$

---

`,
		},
	}

	r := events.MarkdownRenderer{}
	for _, c := range cases {
		t.Run(c.Description, func(t *testing.T) {
			res := command.Result{
				ProjectResults: c.ProjectResults,
			}
			for _, verbose := range []bool{true, false} {
				t.Run(c.Description, func(t *testing.T) {
					s := r.Render(res, c.Command, "log", verbose, c.VCSHost)
					expWithBackticks := strings.Replace(c.Expected, "$", "`", -1)
					if !verbose {
						Equals(t, expWithBackticks, s)
					} else {
						Equals(t, expWithBackticks+"<details><summary>Log</summary>\n  <p>\n\n```\nlog```\n</p></details>\n", s)
					}
				})
			}
		})
	}
}

// Test that if disable apply all is set then the apply all footer is not added
func TestRenderProjectResultsDisableApplyAll(t *testing.T) {
	cases := []struct {
		Description    string
		Command        command.Name
		ProjectResults []command.ProjectResult
		VCSHost        models.VCSHostType
		Expected       string
	}{
		{
			"single successful plan with disable apply all set",
			command.Plan,
			[]command.ProjectResult{
				{
					PlanSuccess: &models.PlanSuccess{
						TerraformOutput: "terraform-output",
						LockURL:         "lock-url",
						RePlanCmd:       "atlantis plan -d path -w workspace",
						ApplyCmd:        "atlantis apply -d path -w workspace",
					},
					Workspace:  "workspace",
					RepoRelDir: "path",
				},
			},
			models.Github,
			`Ran Plan for dir: $path$ workspace: $workspace$

$$$diff
terraform-output
$$$

* :arrow_forward: To **apply** this plan, comment:
    * $atlantis apply -d path -w workspace$
* :put_litter_in_its_place: To **delete** this plan click [here](lock-url)
* :repeat: To **plan** this project again, comment:
    * $atlantis plan -d path -w workspace$


`,
		},
		{
			"single successful plan with project name with disable apply all set",
			command.Plan,
			[]command.ProjectResult{
				{
					PlanSuccess: &models.PlanSuccess{
						TerraformOutput: "terraform-output",
						LockURL:         "lock-url",
						RePlanCmd:       "atlantis plan -d path -w workspace",
						ApplyCmd:        "atlantis apply -d path -w workspace",
					},
					Workspace:   "workspace",
					RepoRelDir:  "path",
					ProjectName: "projectname",
				},
			},
			models.Github,
			`Ran Plan for project: $projectname$ dir: $path$ workspace: $workspace$

$$$diff
terraform-output
$$$

* :arrow_forward: To **apply** this plan, comment:
    * $atlantis apply -d path -w workspace$
* :put_litter_in_its_place: To **delete** this plan click [here](lock-url)
* :repeat: To **plan** this project again, comment:
    * $atlantis plan -d path -w workspace$


`,
		},
		{
			"multiple successful plans, disable apply all set",
			command.Plan,
			[]command.ProjectResult{
				{
					Workspace:  "workspace",
					RepoRelDir: "path",
					PlanSuccess: &models.PlanSuccess{
						TerraformOutput: "terraform-output",
						LockURL:         "lock-url",
						ApplyCmd:        "atlantis apply -d path -w workspace",
						RePlanCmd:       "atlantis plan -d path -w workspace",
					},
				},
				{
					Workspace:   "workspace",
					RepoRelDir:  "path2",
					ProjectName: "projectname",
					PlanSuccess: &models.PlanSuccess{
						TerraformOutput: "terraform-output2",
						LockURL:         "lock-url2",
						ApplyCmd:        "atlantis apply -d path2 -w workspace",
						RePlanCmd:       "atlantis plan -d path2 -w workspace",
					},
				},
			},
			models.Github,
			`Ran Plan for 2 projects:

1. dir: $path$ workspace: $workspace$
1. project: $projectname$ dir: $path2$ workspace: $workspace$

### 1. dir: $path$ workspace: $workspace$
$$$diff
terraform-output
$$$

* :arrow_forward: To **apply** this plan, comment:
    * $atlantis apply -d path -w workspace$
* :put_litter_in_its_place: To **delete** this plan click [here](lock-url)
* :repeat: To **plan** this project again, comment:
    * $atlantis plan -d path -w workspace$

### 2. project: $projectname$ dir: $path2$ workspace: $workspace$
$$$diff
terraform-output2
$$$

* :arrow_forward: To **apply** this plan, comment:
    * $atlantis apply -d path2 -w workspace$
* :put_litter_in_its_place: To **delete** this plan click [here](lock-url2)
* :repeat: To **plan** this project again, comment:
    * $atlantis plan -d path2 -w workspace$


`,
		},
	}
	r := events.MarkdownRenderer{
		DisableApplyAll: true,
	}
	for _, c := range cases {
		t.Run(c.Description, func(t *testing.T) {
			res := command.Result{
				ProjectResults: c.ProjectResults,
			}
			for _, verbose := range []bool{true, false} {
				t.Run(c.Description, func(t *testing.T) {
					s := r.Render(res, c.Command, "log", verbose, c.VCSHost)
					expWithBackticks := strings.Replace(c.Expected, "$", "`", -1)
					if !verbose {
						Equals(t, expWithBackticks, s)
					} else {
						Equals(t, expWithBackticks+"<details><summary>Log</summary>\n  <p>\n\n```\nlog```\n</p></details>\n", s)
					}
				})
			}
		})
	}
}

// Test that if disable apply is set then the apply  footer is not added
func TestRenderProjectResultsDisableApply(t *testing.T) {
	cases := []struct {
		Description    string
		Command        command.Name
		ProjectResults []command.ProjectResult
		VCSHost        models.VCSHostType
		Expected       string
	}{
		{
			"single successful plan with disable apply set",
			command.Plan,
			[]command.ProjectResult{
				{
					PlanSuccess: &models.PlanSuccess{
						TerraformOutput: "terraform-output",
						LockURL:         "lock-url",
						RePlanCmd:       "atlantis plan -d path -w workspace",
						ApplyCmd:        "atlantis apply -d path -w workspace",
					},
					Workspace:  "workspace",
					RepoRelDir: "path",
				},
			},
			models.Github,
			`Ran Plan for dir: $path$ workspace: $workspace$

$$$diff
terraform-output
$$$

* :put_litter_in_its_place: To **delete** this plan click [here](lock-url)
* :repeat: To **plan** this project again, comment:
    * $atlantis plan -d path -w workspace$


`,
		},
		{
			"single successful plan with project name with disable apply set",
			command.Plan,
			[]command.ProjectResult{
				{
					PlanSuccess: &models.PlanSuccess{
						TerraformOutput: "terraform-output",
						LockURL:         "lock-url",
						RePlanCmd:       "atlantis plan -d path -w workspace",
						ApplyCmd:        "atlantis apply -d path -w workspace",
					},
					Workspace:   "workspace",
					RepoRelDir:  "path",
					ProjectName: "projectname",
				},
			},
			models.Github,
			`Ran Plan for project: $projectname$ dir: $path$ workspace: $workspace$

$$$diff
terraform-output
$$$

* :put_litter_in_its_place: To **delete** this plan click [here](lock-url)
* :repeat: To **plan** this project again, comment:
    * $atlantis plan -d path -w workspace$


`,
		},
		{
			"multiple successful plans, disable apply set",
			command.Plan,
			[]command.ProjectResult{
				{
					Workspace:  "workspace",
					RepoRelDir: "path",
					PlanSuccess: &models.PlanSuccess{
						TerraformOutput: "terraform-output",
						LockURL:         "lock-url",
						ApplyCmd:        "atlantis apply -d path -w workspace",
						RePlanCmd:       "atlantis plan -d path -w workspace",
					},
				},
				{
					Workspace:   "workspace",
					RepoRelDir:  "path2",
					ProjectName: "projectname",
					PlanSuccess: &models.PlanSuccess{
						TerraformOutput: "terraform-output2",
						LockURL:         "lock-url2",
						ApplyCmd:        "atlantis apply -d path2 -w workspace",
						RePlanCmd:       "atlantis plan -d path2 -w workspace",
					},
				},
			},
			models.Github,
			`Ran Plan for 2 projects:

1. dir: $path$ workspace: $workspace$
1. project: $projectname$ dir: $path2$ workspace: $workspace$

### 1. dir: $path$ workspace: $workspace$
$$$diff
terraform-output
$$$

* :put_litter_in_its_place: To **delete** this plan click [here](lock-url)
* :repeat: To **plan** this project again, comment:
    * $atlantis plan -d path -w workspace$

### 2. project: $projectname$ dir: $path2$ workspace: $workspace$
$$$diff
terraform-output2
$$$

* :put_litter_in_its_place: To **delete** this plan click [here](lock-url2)
* :repeat: To **plan** this project again, comment:
    * $atlantis plan -d path2 -w workspace$


`,
		},
	}
	r := events.MarkdownRenderer{
		DisableApplyAll: true,
		DisableApply:    true,
	}
	for _, c := range cases {
		t.Run(c.Description, func(t *testing.T) {
			res := command.Result{
				ProjectResults: c.ProjectResults,
			}
			for _, verbose := range []bool{true, false} {
				t.Run(c.Description, func(t *testing.T) {
					s := r.Render(res, c.Command, "log", verbose, c.VCSHost)
					expWithBackticks := strings.Replace(c.Expected, "$", "`", -1)
					if !verbose {
						Equals(t, expWithBackticks, s)
					} else {
						Equals(t, expWithBackticks+"<details><summary>Log</summary>\n  <p>\n\n```\nlog```\n</p></details>\n", s)
					}
				})
			}
		})
	}
}

// Test that if folding is disabled that it's not used.
func TestRenderProjectResults_DisableFolding(t *testing.T) {
	mr := events.MarkdownRenderer{
		DisableMarkdownFolding: true,
	}

	rendered := mr.Render(command.Result{
		ProjectResults: []command.ProjectResult{
			{
				RepoRelDir: ".",
				Workspace:  "default",
				Error:      errors.New(strings.Repeat("line\n", 13)),
			},
		},
	}, command.Plan, "log", false, models.Github)
	Equals(t, false, strings.Contains(rendered, "<details>"))
}

// Test that if the output is longer than 12 lines, it gets wrapped on the right
// VCS hosts during an error.
func TestRenderProjectResults_WrappedErr(t *testing.T) {
	cases := []struct {
		VCSHost                 models.VCSHostType
		GitlabCommonMarkSupport bool
		Output                  string
		ShouldWrap              bool
	}{
		{
			VCSHost:    models.Github,
			Output:     strings.Repeat("line\n", 1),
			ShouldWrap: false,
		},
		{
			VCSHost:    models.Github,
			Output:     strings.Repeat("line\n", 13),
			ShouldWrap: true,
		},
		{
			VCSHost:                 models.Gitlab,
			GitlabCommonMarkSupport: false,
			Output:                  strings.Repeat("line\n", 1),
			ShouldWrap:              false,
		},
		{
			VCSHost:                 models.Gitlab,
			GitlabCommonMarkSupport: false,
			Output:                  strings.Repeat("line\n", 13),
			ShouldWrap:              false,
		},
		{
			VCSHost:                 models.Gitlab,
			GitlabCommonMarkSupport: true,
			Output:                  strings.Repeat("line\n", 1),
			ShouldWrap:              false,
		},
		{
			VCSHost:                 models.Gitlab,
			GitlabCommonMarkSupport: true,
			Output:                  strings.Repeat("line\n", 13),
			ShouldWrap:              true,
		},
		{
			VCSHost:    models.BitbucketCloud,
			Output:     strings.Repeat("line\n", 1),
			ShouldWrap: false,
		},
		{
			VCSHost:    models.BitbucketCloud,
			Output:     strings.Repeat("line\n", 13),
			ShouldWrap: false,
		},
		{
			VCSHost:    models.BitbucketServer,
			Output:     strings.Repeat("line\n", 1),
			ShouldWrap: false,
		},
		{
			VCSHost:    models.BitbucketServer,
			Output:     strings.Repeat("line\n", 13),
			ShouldWrap: false,
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%s_%v", c.VCSHost.String(), c.ShouldWrap),
			func(t *testing.T) {
				mr := events.MarkdownRenderer{
					GitlabSupportsCommonMark: c.GitlabCommonMarkSupport,
				}

				rendered := mr.Render(command.Result{
					ProjectResults: []command.ProjectResult{
						{
							RepoRelDir: ".",
							Workspace:  "default",
							Error:      errors.New(c.Output),
						},
					},
				}, command.Plan, "log", false, c.VCSHost)
				var exp string
				if c.ShouldWrap {
					exp = `Ran Plan for dir: $.$ workspace: $default$

**Plan Error**
<details><summary>Show Output</summary>

$$$
` + c.Output + `
$$$
</details>

`
				} else {
					exp = `Ran Plan for dir: $.$ workspace: $default$

**Plan Error**
$$$
` + c.Output + `
$$$

`
				}

				expWithBackticks := strings.Replace(exp, "$", "`", -1)
				Equals(t, expWithBackticks, rendered)
			})
	}
}

// Test that if the output is longer than 12 lines, it gets wrapped on the right
// VCS hosts for a single project.
func TestRenderProjectResults_WrapSingleProject(t *testing.T) {
	cases := []struct {
		VCSHost                 models.VCSHostType
		GitlabCommonMarkSupport bool
		Output                  string
		ShouldWrap              bool
	}{
		{
			VCSHost:    models.Github,
			Output:     strings.Repeat("line\n", 1),
			ShouldWrap: false,
		},
		{
			VCSHost:    models.Github,
			Output:     strings.Repeat("line\n", 13) + "No changes. Infrastructure is up-to-date.",
			ShouldWrap: true,
		},
		{
			VCSHost:                 models.Gitlab,
			GitlabCommonMarkSupport: false,
			Output:                  strings.Repeat("line\n", 1),
			ShouldWrap:              false,
		},
		{
			VCSHost:                 models.Gitlab,
			GitlabCommonMarkSupport: false,
			Output:                  strings.Repeat("line\n", 13),
			ShouldWrap:              false,
		},
		{
			VCSHost:                 models.Gitlab,
			GitlabCommonMarkSupport: true,
			Output:                  strings.Repeat("line\n", 1),
			ShouldWrap:              false,
		},
		{
			VCSHost:                 models.Gitlab,
			GitlabCommonMarkSupport: true,
			Output:                  strings.Repeat("line\n", 13) + "No changes. Infrastructure is up-to-date.",
			ShouldWrap:              true,
		},
		{
			VCSHost:    models.BitbucketCloud,
			Output:     strings.Repeat("line\n", 1),
			ShouldWrap: false,
		},
		{
			VCSHost:    models.BitbucketCloud,
			Output:     strings.Repeat("line\n", 13),
			ShouldWrap: false,
		},
		{
			VCSHost:    models.BitbucketServer,
			Output:     strings.Repeat("line\n", 1),
			ShouldWrap: false,
		},
		{
			VCSHost:    models.BitbucketServer,
			Output:     strings.Repeat("line\n", 13),
			ShouldWrap: false,
		},
	}

	for _, c := range cases {
		for _, cmd := range []command.Name{command.Plan, command.Apply} {
			t.Run(fmt.Sprintf("%s_%s_%v", c.VCSHost.String(), cmd.String(), c.ShouldWrap),
				func(t *testing.T) {
					mr := events.MarkdownRenderer{
						GitlabSupportsCommonMark: c.GitlabCommonMarkSupport,
					}
					var pr command.ProjectResult
					switch cmd {
					case command.Plan:
						pr = command.ProjectResult{
							RepoRelDir: ".",
							Workspace:  "default",
							PlanSuccess: &models.PlanSuccess{
								TerraformOutput: c.Output,
								LockURL:         "lock-url",
								RePlanCmd:       "replancmd",
								ApplyCmd:        "applycmd",
							},
						}
					case command.Apply:
						pr = command.ProjectResult{
							RepoRelDir:   ".",
							Workspace:    "default",
							ApplySuccess: c.Output,
						}
					}
					rendered := mr.Render(command.Result{
						ProjectResults: []command.ProjectResult{pr},
					}, cmd, "log", false, c.VCSHost)

					// Check result.
					var exp string
					switch cmd {
					case command.Plan:
						if c.ShouldWrap {
							exp = `Ran Plan for dir: $.$ workspace: $default$

<details><summary>Show Output</summary>

$$$diff
` + c.Output + `
$$$

* :arrow_forward: To **apply** this plan, comment:
    * $applycmd$
* :put_litter_in_its_place: To **delete** this plan click [here](lock-url)
* :repeat: To **plan** this project again, comment:
    * $replancmd$
</details>
No changes. Infrastructure is up-to-date.

---
* :fast_forward: To **apply** all unapplied plans from this pull request, comment:
    * $atlantis apply$
* :put_litter_in_its_place: To delete all plans and locks for the PR, comment:
    * $atlantis unlock$
`
						} else {
							exp = `Ran Plan for dir: $.$ workspace: $default$

$$$diff
` + c.Output + `
$$$

* :arrow_forward: To **apply** this plan, comment:
    * $applycmd$
* :put_litter_in_its_place: To **delete** this plan click [here](lock-url)
* :repeat: To **plan** this project again, comment:
    * $replancmd$

---
* :fast_forward: To **apply** all unapplied plans from this pull request, comment:
    * $atlantis apply$
* :put_litter_in_its_place: To delete all plans and locks for the PR, comment:
    * $atlantis unlock$
`
						}
					case command.Apply:
						if c.ShouldWrap {
							exp = `Ran Apply for dir: $.$ workspace: $default$

<details><summary>Show Output</summary>

$$$diff
` + c.Output + `
$$$
</details>

`
						} else {
							exp = `Ran Apply for dir: $.$ workspace: $default$

$$$diff
` + c.Output + `
$$$

`
						}
					}

					expWithBackticks := strings.Replace(exp, "$", "`", -1)
					Equals(t, expWithBackticks, rendered)
				})
		}
	}
}

func TestRenderProjectResults_MultiProjectApplyWrapped(t *testing.T) {
	mr := events.MarkdownRenderer{}
	tfOut := strings.Repeat("line\n", 13)
	rendered := mr.Render(command.Result{
		ProjectResults: []command.ProjectResult{
			{
				RepoRelDir:   ".",
				Workspace:    "staging",
				ApplySuccess: tfOut,
			},
			{
				RepoRelDir:   ".",
				Workspace:    "production",
				ApplySuccess: tfOut,
			},
		},
	}, command.Apply, "log", false, models.Github)
	exp := `Ran Apply for 2 projects:

1. dir: $.$ workspace: $staging$
1. dir: $.$ workspace: $production$

### 1. dir: $.$ workspace: $staging$
<details><summary>Show Output</summary>

$$$diff
` + tfOut + `
$$$
</details>

---
### 2. dir: $.$ workspace: $production$
<details><summary>Show Output</summary>

$$$diff
` + tfOut + `
$$$
</details>

---

`
	expWithBackticks := strings.Replace(exp, "$", "`", -1)
	Equals(t, expWithBackticks, rendered)
}

func TestRenderProjectResults_MultiProjectPlanWrapped(t *testing.T) {
	mr := events.MarkdownRenderer{}
	tfOut := strings.Repeat("line\n", 13) + "Plan: 1 to add, 0 to change, 0 to destroy."
	rendered := mr.Render(command.Result{
		ProjectResults: []command.ProjectResult{
			{
				RepoRelDir: ".",
				Workspace:  "staging",
				PlanSuccess: &models.PlanSuccess{
					TerraformOutput: tfOut,
					LockURL:         "staging-lock-url",
					ApplyCmd:        "staging-apply-cmd",
					RePlanCmd:       "staging-replan-cmd",
				},
			},
			{
				RepoRelDir: ".",
				Workspace:  "production",
				PlanSuccess: &models.PlanSuccess{
					TerraformOutput: tfOut,
					LockURL:         "production-lock-url",
					ApplyCmd:        "production-apply-cmd",
					RePlanCmd:       "production-replan-cmd",
				},
			},
		},
	}, command.Plan, "log", false, models.Github)
	exp := `Ran Plan for 2 projects:

1. dir: $.$ workspace: $staging$
1. dir: $.$ workspace: $production$

### 1. dir: $.$ workspace: $staging$
<details><summary>Show Output</summary>

$$$diff
` + tfOut + `
$$$

* :arrow_forward: To **apply** this plan, comment:
    * $staging-apply-cmd$
* :put_litter_in_its_place: To **delete** this plan click [here](staging-lock-url)
* :repeat: To **plan** this project again, comment:
    * $staging-replan-cmd$
</details>
Plan: 1 to add, 0 to change, 0 to destroy.

---
### 2. dir: $.$ workspace: $production$
<details><summary>Show Output</summary>

$$$diff
` + tfOut + `
$$$

* :arrow_forward: To **apply** this plan, comment:
    * $production-apply-cmd$
* :put_litter_in_its_place: To **delete** this plan click [here](production-lock-url)
* :repeat: To **plan** this project again, comment:
    * $production-replan-cmd$
</details>
Plan: 1 to add, 0 to change, 0 to destroy.

---
* :fast_forward: To **apply** all unapplied plans from this pull request, comment:
    * $atlantis apply$
* :put_litter_in_its_place: To delete all plans and locks for the PR, comment:
    * $atlantis unlock$
`
	expWithBackticks := strings.Replace(exp, "$", "`", -1)
	Equals(t, expWithBackticks, rendered)
}

// Test rendering when there was an error in one of the plans and we deleted
// all the plans as a result.
func TestRenderProjectResults_PlansDeleted(t *testing.T) {
	cases := map[string]struct {
		cr  command.Result
		exp string
	}{
		"one failure": {
			cr: command.Result{
				ProjectResults: []command.ProjectResult{
					{
						RepoRelDir: ".",
						Workspace:  "staging",
						Failure:    "failure",
					},
				},
				PlansDeleted: true,
			},
			exp: `Ran Plan for dir: $.$ workspace: $staging$

**Plan Failed**: failure

`,
		},
		"two failures": {
			cr: command.Result{
				ProjectResults: []command.ProjectResult{
					{
						RepoRelDir: ".",
						Workspace:  "staging",
						Failure:    "failure",
					},
					{
						RepoRelDir: ".",
						Workspace:  "production",
						Failure:    "failure",
					},
				},
				PlansDeleted: true,
			},
			exp: `Ran Plan for 2 projects:

1. dir: $.$ workspace: $staging$
1. dir: $.$ workspace: $production$

### 1. dir: $.$ workspace: $staging$
**Plan Failed**: failure

---
### 2. dir: $.$ workspace: $production$
**Plan Failed**: failure

---

`,
		},
		"one failure, one success": {
			cr: command.Result{
				ProjectResults: []command.ProjectResult{
					{
						RepoRelDir: ".",
						Workspace:  "staging",
						Failure:    "failure",
					},
					{
						RepoRelDir: ".",
						Workspace:  "production",
						PlanSuccess: &models.PlanSuccess{
							TerraformOutput: "tf out",
							LockURL:         "lock-url",
							RePlanCmd:       "re-plan cmd",
							ApplyCmd:        "apply cmd",
						},
					},
				},
				PlansDeleted: true,
			},
			exp: `Ran Plan for 2 projects:

1. dir: $.$ workspace: $staging$
1. dir: $.$ workspace: $production$

### 1. dir: $.$ workspace: $staging$
**Plan Failed**: failure

---
### 2. dir: $.$ workspace: $production$
$$$diff
tf out
$$$

This plan was not saved because one or more projects failed and automerge requires all plans pass.

---

`,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			mr := events.MarkdownRenderer{}
			rendered := mr.Render(c.cr, command.Plan, "log", false, models.Github)
			expWithBackticks := strings.Replace(c.exp, "$", "`", -1)
			Equals(t, expWithBackticks, rendered)
		})
	}
}

// test that id repo locking is disabled the link to unlock the project is not rendered
func TestRenderProjectResultsWithRepoLockingDisabled(t *testing.T) {
	cases := []struct {
		Description    string
		Command        command.Name
		ProjectResults []command.ProjectResult
		VCSHost        models.VCSHostType
		Expected       string
	}{
		{
			"no projects",
			command.Plan,
			[]command.ProjectResult{},
			models.Github,
			"Ran Plan for 0 projects:\n\n\n\n",
		},
		{
			"single successful plan",
			command.Plan,
			[]command.ProjectResult{
				{
					PlanSuccess: &models.PlanSuccess{
						TerraformOutput: "terraform-output",
						LockURL:         "lock-url",
						RePlanCmd:       "atlantis plan -d path -w workspace",
						ApplyCmd:        "atlantis apply -d path -w workspace",
					},
					Workspace:  "workspace",
					RepoRelDir: "path",
				},
			},
			models.Github,
			`Ran Plan for dir: $path$ workspace: $workspace$

$$$diff
terraform-output
$$$

* :arrow_forward: To **apply** this plan, comment:
    * $atlantis apply -d path -w workspace$
* :repeat: To **plan** this project again, comment:
    * $atlantis plan -d path -w workspace$

---
* :fast_forward: To **apply** all unapplied plans from this pull request, comment:
    * $atlantis apply$
* :put_litter_in_its_place: To delete all plans and locks for the PR, comment:
    * $atlantis unlock$
`,
		},
		{
			"single successful plan with master ahead",
			command.Plan,
			[]command.ProjectResult{
				{
					PlanSuccess: &models.PlanSuccess{
						TerraformOutput: "terraform-output",
						LockURL:         "lock-url",
						RePlanCmd:       "atlantis plan -d path -w workspace",
						ApplyCmd:        "atlantis apply -d path -w workspace",
						HasDiverged:     true,
					},
					Workspace:  "workspace",
					RepoRelDir: "path",
				},
			},
			models.Github,
			`Ran Plan for dir: $path$ workspace: $workspace$

$$$diff
terraform-output
$$$

* :arrow_forward: To **apply** this plan, comment:
    * $atlantis apply -d path -w workspace$
* :repeat: To **plan** this project again, comment:
    * $atlantis plan -d path -w workspace$

:warning: The branch we're merging into is ahead, it is recommended to pull new commits first.

---
* :fast_forward: To **apply** all unapplied plans from this pull request, comment:
    * $atlantis apply$
* :put_litter_in_its_place: To delete all plans and locks for the PR, comment:
    * $atlantis unlock$
`,
		},
		{
			"single successful plan with project name",
			command.Plan,
			[]command.ProjectResult{
				{
					PlanSuccess: &models.PlanSuccess{
						TerraformOutput: "terraform-output",
						LockURL:         "lock-url",
						RePlanCmd:       "atlantis plan -d path -w workspace",
						ApplyCmd:        "atlantis apply -d path -w workspace",
					},
					Workspace:   "workspace",
					RepoRelDir:  "path",
					ProjectName: "projectname",
				},
			},
			models.Github,
			`Ran Plan for project: $projectname$ dir: $path$ workspace: $workspace$

$$$diff
terraform-output
$$$

* :arrow_forward: To **apply** this plan, comment:
    * $atlantis apply -d path -w workspace$
* :repeat: To **plan** this project again, comment:
    * $atlantis plan -d path -w workspace$

---
* :fast_forward: To **apply** all unapplied plans from this pull request, comment:
    * $atlantis apply$
* :put_litter_in_its_place: To delete all plans and locks for the PR, comment:
    * $atlantis unlock$
`,
		},
		{
			"single successful apply",
			command.Apply,
			[]command.ProjectResult{
				{
					ApplySuccess: "success",
					Workspace:    "workspace",
					RepoRelDir:   "path",
				},
			},
			models.Github,
			`Ran Apply for dir: $path$ workspace: $workspace$

$$$diff
success
$$$

`,
		},
		{
			"single successful apply with project name",
			command.Apply,
			[]command.ProjectResult{
				{
					ApplySuccess: "success",
					Workspace:    "workspace",
					RepoRelDir:   "path",
					ProjectName:  "projectname",
				},
			},
			models.Github,
			`Ran Apply for project: $projectname$ dir: $path$ workspace: $workspace$

$$$diff
success
$$$

`,
		},
		{
			"multiple successful plans",
			command.Plan,
			[]command.ProjectResult{
				{
					Workspace:  "workspace",
					RepoRelDir: "path",
					PlanSuccess: &models.PlanSuccess{
						TerraformOutput: "terraform-output",
						LockURL:         "lock-url",
						ApplyCmd:        "atlantis apply -d path -w workspace",
						RePlanCmd:       "atlantis plan -d path -w workspace",
					},
				},
				{
					Workspace:   "workspace",
					RepoRelDir:  "path2",
					ProjectName: "projectname",
					PlanSuccess: &models.PlanSuccess{
						TerraformOutput: "terraform-output2",
						LockURL:         "lock-url2",
						ApplyCmd:        "atlantis apply -d path2 -w workspace",
						RePlanCmd:       "atlantis plan -d path2 -w workspace",
					},
				},
			},
			models.Github,
			`Ran Plan for 2 projects:

1. dir: $path$ workspace: $workspace$
1. project: $projectname$ dir: $path2$ workspace: $workspace$

### 1. dir: $path$ workspace: $workspace$
$$$diff
terraform-output
$$$

* :arrow_forward: To **apply** this plan, comment:
    * $atlantis apply -d path -w workspace$
* :repeat: To **plan** this project again, comment:
    * $atlantis plan -d path -w workspace$

---
### 2. project: $projectname$ dir: $path2$ workspace: $workspace$
$$$diff
terraform-output2
$$$

* :arrow_forward: To **apply** this plan, comment:
    * $atlantis apply -d path2 -w workspace$
* :repeat: To **plan** this project again, comment:
    * $atlantis plan -d path2 -w workspace$

---
* :fast_forward: To **apply** all unapplied plans from this pull request, comment:
    * $atlantis apply$
* :put_litter_in_its_place: To delete all plans and locks for the PR, comment:
    * $atlantis unlock$
`,
		},
		{
			"multiple successful applies",
			command.Apply,
			[]command.ProjectResult{
				{
					RepoRelDir:   "path",
					Workspace:    "workspace",
					ProjectName:  "projectname",
					ApplySuccess: "success",
				},
				{
					RepoRelDir:   "path2",
					Workspace:    "workspace",
					ApplySuccess: "success2",
				},
			},
			models.Github,
			`Ran Apply for 2 projects:

1. project: $projectname$ dir: $path$ workspace: $workspace$
1. dir: $path2$ workspace: $workspace$

### 1. project: $projectname$ dir: $path$ workspace: $workspace$
$$$diff
success
$$$

---
### 2. dir: $path2$ workspace: $workspace$
$$$diff
success2
$$$

---

`,
		},
		{
			"single errored plan",
			command.Plan,
			[]command.ProjectResult{
				{
					Error:      errors.New("error"),
					RepoRelDir: "path",
					Workspace:  "workspace",
				},
			},
			models.Github,
			`Ran Plan for dir: $path$ workspace: $workspace$

**Plan Error**
$$$
error
$$$

`,
		},
		{
			"single failed plan",
			command.Plan,
			[]command.ProjectResult{
				{
					RepoRelDir: "path",
					Workspace:  "workspace",
					Failure:    "failure",
				},
			},
			models.Github,
			`Ran Plan for dir: $path$ workspace: $workspace$

**Plan Failed**: failure

`,
		},
		{
			"successful, failed, and errored plan",
			command.Plan,
			[]command.ProjectResult{
				{
					Workspace:  "workspace",
					RepoRelDir: "path",
					PlanSuccess: &models.PlanSuccess{
						TerraformOutput: "terraform-output",
						LockURL:         "lock-url",
						ApplyCmd:        "atlantis apply -d path -w workspace",
						RePlanCmd:       "atlantis plan -d path -w workspace",
					},
				},
				{
					Workspace:  "workspace",
					RepoRelDir: "path2",
					Failure:    "failure",
				},
				{
					Workspace:   "workspace",
					RepoRelDir:  "path3",
					ProjectName: "projectname",
					Error:       errors.New("error"),
				},
			},
			models.Github,
			`Ran Plan for 3 projects:

1. dir: $path$ workspace: $workspace$
1. dir: $path2$ workspace: $workspace$
1. project: $projectname$ dir: $path3$ workspace: $workspace$

### 1. dir: $path$ workspace: $workspace$
$$$diff
terraform-output
$$$

* :arrow_forward: To **apply** this plan, comment:
    * $atlantis apply -d path -w workspace$
* :repeat: To **plan** this project again, comment:
    * $atlantis plan -d path -w workspace$

---
### 2. dir: $path2$ workspace: $workspace$
**Plan Failed**: failure

---
### 3. project: $projectname$ dir: $path3$ workspace: $workspace$
**Plan Error**
$$$
error
$$$

---
* :fast_forward: To **apply** all unapplied plans from this pull request, comment:
    * $atlantis apply$
* :put_litter_in_its_place: To delete all plans and locks for the PR, comment:
    * $atlantis unlock$
`,
		},
		{
			"successful, failed, and errored apply",
			command.Apply,
			[]command.ProjectResult{
				{
					Workspace:    "workspace",
					RepoRelDir:   "path",
					ApplySuccess: "success",
				},
				{
					Workspace:  "workspace",
					RepoRelDir: "path2",
					Failure:    "failure",
				},
				{
					Workspace:  "workspace",
					RepoRelDir: "path3",
					Error:      errors.New("error"),
				},
			},
			models.Github,
			`Ran Apply for 3 projects:

1. dir: $path$ workspace: $workspace$
1. dir: $path2$ workspace: $workspace$
1. dir: $path3$ workspace: $workspace$

### 1. dir: $path$ workspace: $workspace$
$$$diff
success
$$$

---
### 2. dir: $path2$ workspace: $workspace$
**Apply Failed**: failure

---
### 3. dir: $path3$ workspace: $workspace$
**Apply Error**
$$$
error
$$$

---

`,
		},
		{
			"successful, failed, and errored apply",
			command.Apply,
			[]command.ProjectResult{
				{
					Workspace:    "workspace",
					RepoRelDir:   "path",
					ApplySuccess: "success",
				},
				{
					Workspace:  "workspace",
					RepoRelDir: "path2",
					Failure:    "failure",
				},
				{
					Workspace:  "workspace",
					RepoRelDir: "path3",
					Error:      errors.New("error"),
				},
			},
			models.Github,
			`Ran Apply for 3 projects:

1. dir: $path$ workspace: $workspace$
1. dir: $path2$ workspace: $workspace$
1. dir: $path3$ workspace: $workspace$

### 1. dir: $path$ workspace: $workspace$
$$$diff
success
$$$

---
### 2. dir: $path2$ workspace: $workspace$
**Apply Failed**: failure

---
### 3. dir: $path3$ workspace: $workspace$
**Apply Error**
$$$
error
$$$

---

`,
		},
	}

	r := events.MarkdownRenderer{}
	r.DisableRepoLocking = true
	for _, c := range cases {
		t.Run(c.Description, func(t *testing.T) {
			res := command.Result{
				ProjectResults: c.ProjectResults,
			}
			for _, verbose := range []bool{true, false} {
				t.Run(c.Description, func(t *testing.T) {
					s := r.Render(res, c.Command, "log", verbose, c.VCSHost)
					expWithBackticks := strings.Replace(c.Expected, "$", "`", -1)
					if !verbose {
						Equals(t, expWithBackticks, s)
					} else {
						Equals(t, expWithBackticks+"<details><summary>Log</summary>\n  <p>\n\n```\nlog```\n</p></details>\n", s)
					}
				})
			}
		})
	}
}

func TestRenderProjectResultsWithEnableDiffMarkdownFormat(t *testing.T) {
	tfOutput := `An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
~ update in-place
-/+ destroy and then create replacement

Terraform will perform the following actions:

  # module.redacted.aws_instance.redacted must be replaced
-/+ resource "aws_instance" "redacted" {
      ~ ami                          = "ami-redacted" -> "ami-redacted" # forces replacement
      ~ arn                          = "arn:aws:ec2:us-east-1:redacted:instance/i-redacted" -> (known after apply)
      ~ associate_public_ip_address  = false -> (known after apply)
        availability_zone            = "us-east-1b"
      ~ cpu_core_count               = 4 -> (known after apply)
      ~ cpu_threads_per_core         = 2 -> (known after apply)
      - disable_api_termination      = false -> null
      - ebs_optimized                = false -> null
        get_password_data            = false
      - hibernation                  = false -> null
      + host_id                      = (known after apply)
        iam_instance_profile         = "remote_redacted_profile"
      ~ id                           = "i-redacted" -> (known after apply)
      ~ instance_state               = "running" -> (known after apply)
        instance_type                = "c5.2xlarge"
      ~ ipv6_address_count           = 0 -> (known after apply)
      ~ ipv6_addresses               = [] -> (known after apply)
        key_name                     = "RedactedRedactedRedacted"
      - monitoring                   = false -> null
      + outpost_arn                  = (known after apply)
      + password_data                = (known after apply)
      + placement_group              = (known after apply)
      ~ primary_network_interface_id = "eni-redacted" -> (known after apply)
      ~ private_dns                  = "ip-redacted.ec2.internal" -> (known after apply)
      ~ private_ip                   = "redacted" -> (known after apply)
      + public_dns                   = (known after apply)
      + public_ip                    = (known after apply)
      ~ secondary_private_ips        = [] -> (known after apply)
      ~ security_groups              = [] -> (known after apply)
        source_dest_check            = true
        subnet_id                    = "subnet-redacted"
        tags                         = {
            "Name" = "redacted-redacted"
        }
      ~ tenancy                      = "default" -> (known after apply)
        user_data                    = "redacted"
      ~ volume_tags                  = {} -> (known after apply)
        vpc_security_group_ids       = [
            "sg-redactedsecuritygroup",
        ]

      + ebs_block_device {
          + delete_on_termination = (known after apply)
          + device_name           = (known after apply)
          + encrypted             = (known after apply)
          + iops                  = (known after apply)
          + kms_key_id            = (known after apply)
          + snapshot_id           = (known after apply)
          + volume_id             = (known after apply)
          + volume_size           = (known after apply)
          + volume_type           = (known after apply)
        }

      + ephemeral_block_device {
          + device_name  = (known after apply)
          + no_device    = (known after apply)
          + virtual_name = (known after apply)
        }

      ~ metadata_options {
          ~ http_endpoint               = "enabled" -> (known after apply)
          ~ http_put_response_hop_limit = 1 -> (known after apply)
          ~ http_tokens                 = "optional" -> (known after apply)
        }

      + network_interface {
          + delete_on_termination = (known after apply)
          + device_index          = (known after apply)
          + network_interface_id  = (known after apply)
        }

      ~ root_block_device {
          ~ delete_on_termination = true -> (known after apply)
          ~ device_name           = "/dev/sda1" -> (known after apply)
          ~ encrypted             = false -> (known after apply)
          ~ iops                  = 600 -> (known after apply)
          + kms_key_id            = (known after apply)
          ~ volume_id             = "vol-redacted" -> (known after apply)
          ~ volume_size           = 200 -> (known after apply)
          ~ volume_type           = "gp2" -> (known after apply)
        }
    }

  # module.redacted.aws_route53_record.redacted_record will be updated in-place
~ resource "aws_route53_record" "redacted_record" {
        fqdn    = "redacted.redacted.redacted.io"
        id      = "redacted_redacted.redacted.redacted.io_A"
        name    = "redacted.redacted.redacted.io"
      ~ records = [
            "foo",
          - "redacted",
        ] -> (known after apply)
        ttl     = 300
        type    = "A"
        zone_id = "redacted"
    }

  # module.redacted.aws_route53_record.redacted_record_2 will be created
+ resource "aws_route53_record" "redacted_record" {
      + fqdn    = "redacted.redacted.redacted.io"
      + id      = "redacted_redacted.redacted.redacted.io_A"
      + name    = "redacted.redacted.redacted.io"
      + records = [
            "foo",
        ]
      + ttl     = 300
      + type    = "A"
      + zone_id = "redacted"
    }

# helm_release.external_dns[0] will be updated in-place
~ resource "helm_release" "external_dns" {
      id                         = "external-dns"
      name                       = "external-dns"
    ~ values                     = [
        - <<-EOT
            image:
              tag: "0.12.0"
              pullSecrets:
              - XXXXX

            domainFilters: ["xxxxx","xxxxx"]
            base64:
              +dGhpcyBpcyBzb21lIHN0cmluZyBvciBzb21ldGhpbmcKCg==
          EOT,
        + <<-EOT
            image:
              tag: "0.12.0"
              pullSecrets:
              - XXXXX

            domainFilters: ["xxxxx","xxxxx"]
            base64:
              +dGhpcyBpcyBzb21lIHN0cmluZyBvciBzb21ldGhpbmcKCg==
          EOT,
      ]
    }

Plan: 1 to add, 1 to change, 1 to destroy.
`
	cases := []struct {
		Description    string
		Command        command.Name
		ProjectResults []command.ProjectResult
		VCSHost        models.VCSHostType
		Expected       string
	}{
		{
			"single successful plan with diff markdown formatted",
			command.Plan,
			[]command.ProjectResult{
				{
					PlanSuccess: &models.PlanSuccess{
						TerraformOutput: tfOutput,
						LockURL:         "lock-url",
						RePlanCmd:       "atlantis plan -d path -w workspace",
						ApplyCmd:        "atlantis apply -d path -w workspace",
					},
					Workspace:  "workspace",
					RepoRelDir: "path",
				},
			},
			models.Github,
			`Ran Plan for dir: $path$ workspace: $workspace$

<details><summary>Show Output</summary>

$$$diff
An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
! update in-place
-/+ destroy and then create replacement

Terraform will perform the following actions:

  # module.redacted.aws_instance.redacted must be replaced
-/+ resource "aws_instance" "redacted" {
!       ami                          = "ami-redacted" -> "ami-redacted" # forces replacement
!       arn                          = "arn:aws:ec2:us-east-1:redacted:instance/i-redacted" -> (known after apply)
!       associate_public_ip_address  = false -> (known after apply)
        availability_zone            = "us-east-1b"
!       cpu_core_count               = 4 -> (known after apply)
!       cpu_threads_per_core         = 2 -> (known after apply)
-       disable_api_termination      = false -> null
-       ebs_optimized                = false -> null
        get_password_data            = false
-       hibernation                  = false -> null
+       host_id                      = (known after apply)
        iam_instance_profile         = "remote_redacted_profile"
!       id                           = "i-redacted" -> (known after apply)
!       instance_state               = "running" -> (known after apply)
        instance_type                = "c5.2xlarge"
!       ipv6_address_count           = 0 -> (known after apply)
!       ipv6_addresses               = [] -> (known after apply)
        key_name                     = "RedactedRedactedRedacted"
-       monitoring                   = false -> null
+       outpost_arn                  = (known after apply)
+       password_data                = (known after apply)
+       placement_group              = (known after apply)
!       primary_network_interface_id = "eni-redacted" -> (known after apply)
!       private_dns                  = "ip-redacted.ec2.internal" -> (known after apply)
!       private_ip                   = "redacted" -> (known after apply)
+       public_dns                   = (known after apply)
+       public_ip                    = (known after apply)
!       secondary_private_ips        = [] -> (known after apply)
!       security_groups              = [] -> (known after apply)
        source_dest_check            = true
        subnet_id                    = "subnet-redacted"
        tags                         = {
            "Name" = "redacted-redacted"
        }
!       tenancy                      = "default" -> (known after apply)
        user_data                    = "redacted"
!       volume_tags                  = {} -> (known after apply)
        vpc_security_group_ids       = [
            "sg-redactedsecuritygroup",
        ]

+       ebs_block_device {
+           delete_on_termination = (known after apply)
+           device_name           = (known after apply)
+           encrypted             = (known after apply)
+           iops                  = (known after apply)
+           kms_key_id            = (known after apply)
+           snapshot_id           = (known after apply)
+           volume_id             = (known after apply)
+           volume_size           = (known after apply)
+           volume_type           = (known after apply)
        }

+       ephemeral_block_device {
+           device_name  = (known after apply)
+           no_device    = (known after apply)
+           virtual_name = (known after apply)
        }

!       metadata_options {
!           http_endpoint               = "enabled" -> (known after apply)
!           http_put_response_hop_limit = 1 -> (known after apply)
!           http_tokens                 = "optional" -> (known after apply)
        }

+       network_interface {
+           delete_on_termination = (known after apply)
+           device_index          = (known after apply)
+           network_interface_id  = (known after apply)
        }

!       root_block_device {
!           delete_on_termination = true -> (known after apply)
!           device_name           = "/dev/sda1" -> (known after apply)
!           encrypted             = false -> (known after apply)
!           iops                  = 600 -> (known after apply)
+           kms_key_id            = (known after apply)
!           volume_id             = "vol-redacted" -> (known after apply)
!           volume_size           = 200 -> (known after apply)
!           volume_type           = "gp2" -> (known after apply)
        }
    }

  # module.redacted.aws_route53_record.redacted_record will be updated in-place
! resource "aws_route53_record" "redacted_record" {
        fqdn    = "redacted.redacted.redacted.io"
        id      = "redacted_redacted.redacted.redacted.io_A"
        name    = "redacted.redacted.redacted.io"
!       records = [
            "foo",
-           "redacted",
        ] -> (known after apply)
        ttl     = 300
        type    = "A"
        zone_id = "redacted"
    }

  # module.redacted.aws_route53_record.redacted_record_2 will be created
+ resource "aws_route53_record" "redacted_record" {
+       fqdn    = "redacted.redacted.redacted.io"
+       id      = "redacted_redacted.redacted.redacted.io_A"
+       name    = "redacted.redacted.redacted.io"
+       records = [
            "foo",
        ]
+       ttl     = 300
+       type    = "A"
+       zone_id = "redacted"
    }

# helm_release.external_dns[0] will be updated in-place
! resource "helm_release" "external_dns" {
      id                         = "external-dns"
      name                       = "external-dns"
!     values                     = [
-         <<-EOT
            image:
              tag: "0.12.0"
              pullSecrets:
              - XXXXX

            domainFilters: ["xxxxx","xxxxx"]
            base64:
              +dGhpcyBpcyBzb21lIHN0cmluZyBvciBzb21ldGhpbmcKCg==
          EOT,
+         <<-EOT
            image:
              tag: "0.12.0"
              pullSecrets:
              - XXXXX

            domainFilters: ["xxxxx","xxxxx"]
            base64:
              +dGhpcyBpcyBzb21lIHN0cmluZyBvciBzb21ldGhpbmcKCg==
          EOT,
      ]
    }

Plan: 1 to add, 1 to change, 1 to destroy.

$$$

* :put_litter_in_its_place: To **delete** this plan click [here](lock-url)
* :repeat: To **plan** this project again, comment:
    * $atlantis plan -d path -w workspace$
</details>
Plan: 1 to add, 1 to change, 1 to destroy.


`,
		},
	}
	r := events.MarkdownRenderer{
		DisableApplyAll:          true,
		DisableApply:             true,
		EnableDiffMarkdownFormat: true,
	}
	for _, c := range cases {
		t.Run(c.Description, func(t *testing.T) {
			res := command.Result{
				ProjectResults: c.ProjectResults,
			}
			for _, verbose := range []bool{true, false} {
				t.Run(c.Description, func(t *testing.T) {
					s := r.Render(res, c.Command, "log", verbose, c.VCSHost)
					expWithBackticks := strings.Replace(c.Expected, "$", "`", -1)
					if !verbose {
						Equals(t, expWithBackticks, s)
					} else {
						Equals(t, expWithBackticks+"<details><summary>Log</summary>\n  <p>\n\n```\nlog```\n</p></details>\n", s)
					}
				})
			}
		})
	}
}
