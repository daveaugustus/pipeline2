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

func ParseKeyDump(ctx context.Context, fileDest <-chan string) <-chan interface{} {
	keyDump := make(chan interface{})

	go func() {
		defer close(keyDump)
	}()
	return keyDump
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
		confUser <- User{Username: "Alpha"}
	}()
	return confUser
}

// TODO: FIX this
func OrgMembers(ctx context.Context, user <-chan User) <-chan interface{} {
	confUser := make(chan interface{})

	go func() {
		defer close(confUser)
		confUser <- Org{Name: "Add"}
	}()
	return confUser
}

func AdminUsers(ctx context.Context, user <-chan User) <-chan User {
	adminUser := make(chan User)

	go func() {
		defer close(adminUser)
		// Check if the user is in the DB
		adminUser <- User{Username: "Alpha"}
	}()
	return adminUser
}

func RunComplexPipeline(base int, lines []string) error {
	fmt.Printf("runComplexPipeline: base=%v, lines=%v\n", base, lines)

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	var errcList []<-chan error

	// Source pipeline stage.
	linec, errc, err := lineListSource(ctx, lines...)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	// Transformer pipeline stage 1.
	numberc, errc, err := lineParser(ctx, base, linec)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	// Transformer pipeline stage 2.
	numberc1, numberc2, errc, err := splitter(ctx, numberc)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	// Transformer pipeline stage 3.
	numberc3, errc, err := squarer(ctx, numberc1)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	// Sink pipeline stage 1.
	errc, err = sink(ctx, numberc3)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	// Sink pipeline stage 2.
	errc, err = sink(ctx, numberc2)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	fmt.Println("Pipeline started. Waiting for pipeline to complete.")

	return WaitForPipeline(errcList...)
}
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
}
