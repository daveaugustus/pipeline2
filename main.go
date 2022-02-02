package main

import (
	"context"
	"fmt"
	"pipeline2/pipeline"
	"time"
)

type MigrationPipe func(<-chan pipeline.Result) <-chan pipeline.Result

// ParseOrg returns MigrationPipe
func UnzipSrc(res <-chan pipeline.Result) MigrationPipe {
	return func(in <-chan pipeline.Result) <-chan pipeline.Result {
		return unzip(context.Background(), res)
	}
}

func unzip(ctx context.Context, result <-chan pipeline.Result) <-chan pipeline.Result {
	fmt.Println("Unzipping...!")
	out := make(chan pipeline.Result)
	go func() {

		for res := range result {
			out <- res
		}

	}()
	fmt.Println("Unzipping done!")
	return out
}

// ParseOrg returns MigrationPipe
func ParseOrg(res <-chan pipeline.Result) MigrationPipe {
	return func(in <-chan pipeline.Result) <-chan pipeline.Result {
		return parseOrg(context.Background(), res)
	}
}

func parseOrg(ctx context.Context, result <-chan pipeline.Result) <-chan pipeline.Result {
	fmt.Println("Parsing orgs...!")

	out := make(chan pipeline.Result)
	defer close(out)

	go func() {

		for res := range result {
			out <- res
		}

	}()

	fmt.Println("Done parsing orgs!")
	return out
}

// ParseUser returns MigrationPipe
func ParseUser(res <-chan pipeline.Result) MigrationPipe {
	return func(in <-chan pipeline.Result) <-chan pipeline.Result {
		return parseUser(context.Background(), res)
	}
}

func parseUser(ctx context.Context, result <-chan pipeline.Result) <-chan pipeline.Result {
	fmt.Println("Parsing users...!")

	out := make(chan pipeline.Result)
	defer close(out)

	go func() {

		for res := range result {
			out <- res
		}

	}()

	fmt.Println("Done parsing users!")
	return out
}

// ConflictingUsers returns MigrationPipe
func ConflictingUsers(res <-chan pipeline.Result) MigrationPipe {
	return func(in <-chan pipeline.Result) <-chan pipeline.Result {
		return conflictingUsers(context.Background(), res)
	}
}

func conflictingUsers(ctx context.Context, result <-chan pipeline.Result) <-chan pipeline.Result {
	fmt.Println("Checking for conflicting users...!")

	out := make(chan pipeline.Result)
	defer close(out)

	go func() {

		for res := range result {
			out <- res
		}

	}()

	fmt.Println("Completed conflicting users check!")
	return out
}

// OrgMembers returns MigrationPipe
func OrgMembers(res <-chan pipeline.Result) MigrationPipe {
	return func(in <-chan pipeline.Result) <-chan pipeline.Result {
		return orgMembers(context.Background(), res)
	}
}

func orgMembers(ctx context.Context, result <-chan pipeline.Result) <-chan pipeline.Result {
	fmt.Println("Checking org, user associatian...!")

	out := make(chan pipeline.Result)
	defer close(out)

	go func() {

		for res := range result {
			out <- res
		}

	}()

	fmt.Println("Done checking associatian!")
	return out
}

// AdminUsers Return MigrationPipe
func AdminUsers(res <-chan pipeline.Result) MigrationPipe {
	return func(in <-chan pipeline.Result) <-chan pipeline.Result {
		return adminUsers(context.Background(), res)
	}
}

func adminUsers(ctx context.Context, result <-chan pipeline.Result) <-chan pipeline.Result {
	fmt.Println("Checking for admin users..!")

	out := make(chan pipeline.Result)
	defer close(out)

	go func() {

		for res := range result {
			out <- res
		}

	}()

	fmt.Println("Done admin user check!")
	return out
}

func migrationPipeline(source <-chan pipeline.Result, pipes ...MigrationPipe) {
	fmt.Println("Pipeline started...")
	msg := make(chan string)

	go func() {
		for _, pipe := range pipes {
			time.Sleep(time.Second)
			source = pipe(source)
		}

		for s := range source {
			fmt.Println(s)
		}
		msg <- "Done"
	}()

	fmt.Println("Pipeline Status: ", <-msg)
}

func RunPhaseOnePipeline(src string) {
	c := make(chan pipeline.Result, 1)

	c <- pipeline.Result{}
	close(c)

	migrationPipeline(c,
		UnzipSrc(c),
		ParseOrg(c),
		ParseUser(c),
		ConflictingUsers(c),
		OrgMembers(c),
		AdminUsers(c),
	)
}

func main() {
	RunPhaseOnePipeline("")
}
