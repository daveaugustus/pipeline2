package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"pipeline2/pipeline"
	"sync"
)

type PipelineData struct {
	Result pipeline.Result
	Done   chan<- error
	Ctx    context.Context
}

type PhaseOnePipleine struct {
	in chan<- PipelineData
}

type PhaseOnePipelineProcessor func(<-chan PipelineData) <-chan PipelineData

// ParseOrg returns PhaseOnePipelineProcessor
func UnzipSrc() PhaseOnePipelineProcessor {
	return func(result <-chan PipelineData) <-chan PipelineData {
		return unzip(result)
	}
}

func unzip(result <-chan PipelineData) <-chan PipelineData {
	fmt.Println("Starting unzip pipeline")
	out := make(chan PipelineData, 100)
	go func() {

		for res := range result {
			res.Result.Meta.UnzipFolder = "backup"
			log.Printf("UnzipSrc: %+v", res)
			select {
			case out <- res:
			case <-res.Ctx.Done():
				res.Done <- nil
			}
		}
		fmt.Println("Closing unzip")
		close(out)
	}()
	return out
}

// ParseOrg returns PhaseOnePipelineProcessor
func ParseOrg() PhaseOnePipelineProcessor {
	return func(result <-chan PipelineData) <-chan PipelineData {
		return parseOrg(result)
	}
}

func parseOrg(result <-chan PipelineData) <-chan PipelineData {
	fmt.Println("Starting parse orgs pipeline")

	out := make(chan PipelineData, 100)

	go func() {
		fmt.Println("\nProcessing to parse orgs...")
		orgs := []pipeline.Org{
			{
				Name:      "ABC",
				FullName:  "ABC",
				ActionOps: 1,
			},
			{
				Name:      "PQR",
				FullName:  "PQR",
				ActionOps: 1,
			},
		}

		for res := range result {
			res.Result.ParsedResult.Orgs = orgs
			log.Printf("\nParseOrg: %+v", res)
			select {
			case out <- res:
			case <-res.Ctx.Done():
				res.Done <- nil
			}
			fmt.Println("after write")
		}
		fmt.Println("CLosing orgs pipeline")
		close(out)
	}()
	return out
}

// ParseUser returns PhaseOnePipelineProcessor
func ParseUser() PhaseOnePipelineProcessor {
	return func(result <-chan PipelineData) <-chan PipelineData {
		return parseUser(result)
	}
}

func parseUser(result <-chan PipelineData) <-chan PipelineData {
	fmt.Println("Starting parse user pipeline")

	out := make(chan PipelineData, 100)

	go func() {
		fmt.Println("Processing to parsing users...")
		users := []pipeline.User{
			{
				Username:    "davetweetlive",
				Email:       "dave@mail.com",
				DisplayName: "Dave",
			},
			{
				Username:    "somerandomnun",
				Email:       "someting@mail.com",
				DisplayName: "SOmething",
				IsAdmin:     true,
			},
		}

		for res := range result {
			res.Result.ParsedResult.Users = users
			log.Printf("ParseUser: %+v", res)
			select {
			case out <- res:
			case <-res.Ctx.Done():
				res.Done <- nil
			}
		}
		fmt.Println("Closing parseUser")
		close(out)
	}()

	return out
}

// ConflictingUsers returns PhaseOnePipelineProcessor
func ConflictingUsers() PhaseOnePipelineProcessor {
	return func(result <-chan PipelineData) <-chan PipelineData {
		return conflictingUsers(result)
	}
}

func conflictingUsers(result <-chan PipelineData) <-chan PipelineData {
	fmt.Println("Starting conflicting user check pipeline")

	out := make(chan PipelineData, 100)

	go func() {
		fmt.Println("Processing to check conflicting users...")

		for res := range result {
			res.Result.ParsedResult.Users[0].IsConflicting = true
			log.Printf("Post conflictingUsers process: %+v", res)
			select {
			case out <- res:
			case <-res.Ctx.Done():
				res.Done <- nil
			}
		}
		fmt.Println("Closing conflictUser")
		close(out)

	}()

	return out
}

// OrgMembers returns PhaseOnePipelineProcessor
func OrgMembers() PhaseOnePipelineProcessor {
	return func(result <-chan PipelineData) <-chan PipelineData {
		return orgMembers(result)
	}
}

func orgMembers(result <-chan PipelineData) <-chan PipelineData {
	fmt.Println("Starting org user check pipeline")

	out := make(chan PipelineData, 100)

	go func() {
		fmt.Println("Processing to check org users association...")
		orgUser := []pipeline.OrgsUsersAssociations{
			{
				OrgName: pipeline.Org{
					Name:      "ABC",
					FullName:  "ABC",
					ActionOps: 1,
				},
				Users: []pipeline.UserAssociation{
					{
						Username: "davetweetlive",
						IsAdmin:  false,
					},
				},
			},
		}

		for res := range result {
			res.Result.ParsedResult.OrgsUsers = orgUser
			log.Printf("Post orgMembers process: %+v", res)
			select {
			case out <- res:
			case <-res.Ctx.Done():
				res.Done <- nil
			}
		}
		fmt.Println("Closing orgmember")
		close(out)

	}()

	return out
}

// AdminUsers Return PhaseOnePipelineProcessor
func AdminUsers() PhaseOnePipelineProcessor {
	return func(result <-chan PipelineData) <-chan PipelineData {
		return adminUsers(result)
	}
}

func adminUsers(result <-chan PipelineData) <-chan PipelineData {
	fmt.Println("Starting org admin user check pipeline")

	out := make(chan PipelineData, 100)

	go func() {
		// fmt.Println("Processing to to check admin users...")
		for res := range result {
			select {
			case out <- res:
			case <-res.Ctx.Done():
				res.Done <- nil
			}
		}
		fmt.Println("Closing adminUsers")
		close(out)

	}()

	return out
}

// AdminUsers Return MigrationPipe
func ParseJSON() PhaseOnePipelineProcessor {
	return func(result <-chan PipelineData) <-chan PipelineData {
		return parseJSON(result)
	}
}
func parseJSON(result <-chan PipelineData) <-chan PipelineData {
	fmt.Println("JSON pipeline")

	out := make(chan PipelineData, 100)

	go func() {
		fmt.Println("Processing to JSON...")
		for res := range result {
			if err := res.TOJSON(); err != nil {
				return
			}
			// log.Printf("Binary Result: %+v", res)
			select {
			case out <- res:
			case <-res.Ctx.Done():
				res.Done <- nil
			}
		}
		fmt.Println("Closing JSON")
		// close(out)

	}()

	return out
}

func migrationPipeline(source <-chan PipelineData, pipes ...PhaseOnePipelineProcessor) {
	fmt.Println("Pipeline started...")

	go func() {
		for _, pipe := range pipes {
			source = pipe(source)
		}

		for s := range source {
			s.Done <- nil
		}
	}()
}

func SetupPhaseOnePipeline() PhaseOnePipleine {
	c := make(chan PipelineData, 100)
	migrationPipeline(c,
		UnzipSrc(),
		ParseOrg(),
		ParseUser(),
		ConflictingUsers(),
		OrgMembers(),
		AdminUsers(),
		// ParseJSON(),
	)
	return PhaseOnePipleine{in: c}
}

func (p *PhaseOnePipleine) Run(result pipeline.Result) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		done := make(chan error)
		select {
		case p.in <- PipelineData{Result: result, Done: done, Ctx: ctx}:
		}
		err := <-done
		if err != nil {
			fmt.Println("received error")
		}
		fmt.Println("received done")
		// status <- "done"
	}()
	wg.Wait()
}

func main() {
	fmt.Println("Pipeline running")
	msg := pipeline.Result{Meta: pipeline.Meta{ZipFile: "/home/dave/eureka/pipeline2/backup.zip"}}

	p1Pipeline := SetupPhaseOnePipeline()
	p1Pipeline.Run(msg)
}

func (p PipelineData) TOJSON() error {
	bte, err := json.Marshal(p.Result)
	if err != nil {
		fmt.Println("Marshal Error: ", err)
		return err
	}

	if err := ioutil.WriteFile("dem.json", bte, 0666); err != nil {
		return err
	}
	return nil
}
