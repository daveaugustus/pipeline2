package main

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Unzip(ctx context.Context, src string) <-chan string {
	filePath := make(chan string)
	go func() {
		r, err := zip.OpenReader(src)
		if err != nil {
			fmt.Println("cannot open reader")
			ctx.Done()
		}
		defer r.Close()

		for _, f := range r.File {
			fpath := filepath.Join("", f.Name)

			// Checking for any invalid file paths
			if !strings.HasPrefix(fpath, filepath.Clean("backup")+string(os.PathSeparator)) {
				fmt.Println("invalid path")
				ctx.Done()
			}

			// filenames = append(filenames, fpath)

			if f.FileInfo().IsDir() {
				os.MkdirAll(fpath, os.ModePerm)
				continue
			}

			// Creating the files in the target directory
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				fmt.Println("cannot create dir")
				ctx.Done()
			}

			// The created file will be stored in
			// outFile with permissions to write &/or truncate
			outFile, err := os.OpenFile(fpath,
				os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
				f.Mode())
			if err != nil {
				fmt.Println("cannot create a file")
				ctx.Done()
			}

			rc, err := f.Open()
			if err != nil {
				fmt.Println("cannot open file")
				ctx.Done()
			}

			_, err = io.Copy(outFile, rc)

			outFile.Close()
			rc.Close()

			if err != nil {
				fmt.Println("cannot copy a file")
				ctx.Done()
			}
			filePath <- fpath
		}
		close(filePath)
	}()
	return filePath
}

func ParseOrg(ctx context.Context, fileDest <-chan string) <-chan Org {
	orgChan := make(chan Org)

	go func() {
		defer close(orgChan)
		orgChan <- Org{Name: "Org Alpha"}
		orgChan <- Org{Name: "Org Beta"}
	}()
	return orgChan
}

func ParseKeyDump(ctx context.Context, fileDest <-chan string) <-chan KeyDump {
	keyDump := make(chan KeyDump)

	go func() {
		defer close(keyDump)
		keyDump <- KeyDump{
			Admin: true,
		}
	}()
	return keyDump
}

func ParseUser(ctx context.Context, fileDest <-chan string) <-chan User {
	userChan := make(chan User)

	go func() {
		defer close(userChan)
		userChan <- User{Username: "Alpha"}
		userChan <- User{Username: "Beta"}
	}()
	return userChan
}

func ConflictingUsers(ctx context.Context, user <-chan User) <-chan User {
	confUser := make(chan User)

	go func() {
		defer close(confUser)
		// Check if the user is in the DB
		confUser <- User{Username: "Gamma"}
	}()
	return confUser
}

// TODO: FIX this
func OrgMembers(ctx context.Context, user <-chan User) <-chan map[User][]Org {
	userOrg := make(chan map[User][]Org)
	go func() {
		defer close(userOrg)
		userOrgMap := map[User][]Org{}
		userOrg <- userOrgMap
	}()
	return userOrg
}

func AdminUsers(ctx context.Context, user <-chan User) <-chan User {
	adminUser := make(chan User)

	go func() {
		defer close(adminUser)
		// Check if the user is in the DB
		adminUser <- User{Username: "Delta"}
	}()
	return adminUser
}

func RunComplexPipeline() error {

	// var errcList []<-chan error
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	// Unzipping
	ch := Unzip(ctx, "/home/dave/eureka/data-pipeline-golang/backup.zip")

	orgs := ParseOrg(ctx, ch)

	keyDump := ParseKeyDump(ctx, ch)

	users := ParseUser(ctx, ch)

	existingUsers := ConflictingUsers(ctx, users)

	orgsUser := OrgMembers(ctx, users)

	adminUsers := AdminUsers(ctx, users)

	// orgs
	for org := range orgs {
		fmt.Println("Organization: ", org)
	}
	fmt.Println()

	// Keydump
	for kd := range keyDump {
		fmt.Println("Keydump: ", kd)
	}
	fmt.Println()

	// Existing user
	for eu := range existingUsers {
		fmt.Println("Existing user: ", eu)
	}
	fmt.Println()

	// User's Org
	for ou := range orgsUser {
		fmt.Println("User's Org: ", ou)
	}

	fmt.Println()
	// Admin User
	for au := range adminUsers {
		fmt.Println("Admin User: ", au)
	}

	fmt.Println("Pipeline started. Waiting for pipeline to complete.")

	return nil
}

//
//
//
//
//
//
//
//
//
//
// /
//
//
// /
//
//
//
//
//
// /
//
func main() {
	ch := Unzip(context.Background(), "/home/dave/eureka/data-pipeline-golang/backup.zip")
	count := 1
	for i := range ch {
		fmt.Println(count, i)
		count++
	}

	// Phase 2

	// org := ParseOrg(Unzip("/home/dave/eureka/data-pipeline-golang/backup.zip"))
	// kd := ParseKeyDump(Unzip("/home/dave/eureka/data-pipeline-golang/backup.zip"))
	// // count := 1
	// // for i := range org {
	// // 	// fmt.Println(count, i)
	// // 	count++
	// // }
	// SaveAndUpdate(org, kd)
	RunComplexPipeline()
}
