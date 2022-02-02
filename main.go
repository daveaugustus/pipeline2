package main

import (
	"context"
	"fmt"
	"pipeline2/pipeline"
	"time"
)

type MigrationPipe func(context.Context, <-chan pipeline.Result) <-chan pipeline.Result

// ParseOrg returns MigrationPipe
func UnzipSrc(ctx context.Context, res <-chan pipeline.Result) MigrationPipe {
	return func(ctx context.Context, in <-chan pipeline.Result) <-chan pipeline.Result {
		return unzip(ctx, res)
	}
}

func unzip(ctx context.Context, result <-chan pipeline.Result) <-chan pipeline.Result {
	fmt.Println("Starting to unzip...")
	out := make(chan pipeline.Result)
	go func() {
		fmt.Println("Processing to unzip...")
		for res := range result {
			select {
			case out <- res:
			case <-ctx.Done():
				return
			}
		}

	}()
	fmt.Println("Done unzipping!")
	return out
}

// ParseOrg returns MigrationPipe
func ParseOrg(ctx context.Context, res <-chan pipeline.Result) MigrationPipe {
	return func(ctx context.Context, in <-chan pipeline.Result) <-chan pipeline.Result {
		return parseOrg(ctx, res)
	}
}

func parseOrg(ctx context.Context, result <-chan pipeline.Result) <-chan pipeline.Result {
	fmt.Println("Starting to parse orgs...")

	out := make(chan pipeline.Result)
	defer close(out)

	go func() {
		fmt.Println("Processing to parse orgs...")
		for res := range result {
			select {
			case out <- res:
			case <-ctx.Done():
				return
			}
		}

	}()

	fmt.Println("Done parsing orgs!")
	return out
}

// ParseUser returns MigrationPipe
func ParseUser(ctx context.Context, res <-chan pipeline.Result) MigrationPipe {
	return func(ctx context.Context, in <-chan pipeline.Result) <-chan pipeline.Result {
		return parseUser(ctx, res)
	}
}

func parseUser(ctx context.Context, result <-chan pipeline.Result) <-chan pipeline.Result {
	fmt.Println("Starting to parsing users...")

	out := make(chan pipeline.Result)
	defer close(out)

	go func() {
		fmt.Println("Processing to parsing users...")
		for res := range result {
			select {
			case out <- res:
			case <-ctx.Done():
				return
			}
		}
	}()

	fmt.Println("Done parsing users!")
	return out
}

// ConflictingUsers returns MigrationPipe
func ConflictingUsers(ctx context.Context, res <-chan pipeline.Result) MigrationPipe {
	return func(ctx context.Context, in <-chan pipeline.Result) <-chan pipeline.Result {
		return conflictingUsers(ctx, res)
	}
}

func conflictingUsers(ctx context.Context, result <-chan pipeline.Result) <-chan pipeline.Result {
	fmt.Println("Starting to check conflicting users...")

	out := make(chan pipeline.Result)
	defer close(out)

	go func() {
		fmt.Println("Processing to check conflicting users...")

		for res := range result {
			select {
			case out <- res:
			case <-ctx.Done():
				return
			}
		}

	}()

	fmt.Println("Done check conflicting users!")
	return out
}

// OrgMembers returns MigrationPipe
func OrgMembers(ctx context.Context, res <-chan pipeline.Result) MigrationPipe {
	return func(ctx context.Context, in <-chan pipeline.Result) <-chan pipeline.Result {
		return orgMembers(ctx, res)
	}
}

func orgMembers(ctx context.Context, result <-chan pipeline.Result) <-chan pipeline.Result {
	fmt.Println("Starting to check org users association...")

	out := make(chan pipeline.Result)
	defer close(out)

	go func() {
		fmt.Println("Processing to check org users association...")
		for res := range result {
			select {
			case out <- res:
			case <-ctx.Done():
				return
			}
		}

	}()

	fmt.Println("Done checking associatian!")
	return out
}

// AdminUsers Return MigrationPipe
func AdminUsers(ctx context.Context, res <-chan pipeline.Result) MigrationPipe {
	return func(ctx context.Context, in <-chan pipeline.Result) <-chan pipeline.Result {
		return adminUsers(ctx, res)
	}
}

func adminUsers(ctx context.Context, result <-chan pipeline.Result) <-chan pipeline.Result {
	fmt.Println("Starting to check admin users...")

	out := make(chan pipeline.Result)
	defer close(out)

	go func() {
		fmt.Println("Processing to to check admin users...")
		for res := range result {
			select {
			case out <- res:
			case <-ctx.Done():
				return
			}
		}

	}()

	fmt.Println("Done admin user check!")
	return out
}

func migrationPipeline(ctx context.Context, source <-chan pipeline.Result, pipes ...MigrationPipe) {
	fmt.Println("Pipeline started...")
	msg := make(chan string)

	go func() {
		for _, pipe := range pipes {
			time.Sleep(time.Second)
			fmt.Println()
			source = pipe(ctx, source)
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
	ctx := context.Background()
	c <- pipeline.Result{}
	close(c)

	migrationPipeline(ctx, c,
		UnzipSrc(ctx, c),
		ParseOrg(ctx, c),
		ParseUser(ctx, c),
		ConflictingUsers(ctx, c),
		OrgMembers(ctx, c),
		AdminUsers(ctx, c),
	)
}

func main() {
	RunPhaseOnePipeline("")
}
