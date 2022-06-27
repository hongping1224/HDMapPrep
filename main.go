package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hongping1224/lidario"
)

func generatePointCloudList(inputRoot string) {
	lasPaths := []string{inputRoot}
	lasPaths = findFile(inputRoot, ".las")
	pointCloudListPath := filepath.Join(inputRoot, "point_cloud_map_list.csv")
	f, err := os.Create(pointCloudListPath)
	if err != nil {
		log.Printf("Fail to create csv at %s :  %v", pointCloudListPath, err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()
	for _, lasPath := range lasPaths {
		fmt.Println(lasPath)
		las, err := lidario.NewLasFile(lasPath, "rh")
		if err != nil {
			log.Printf("%s read fail. Err : %v", las, err)
			continue
		}
		defer las.Close()
		xmin := las.Header.MinX
		ymin := las.Header.MinY
		zmin := las.Header.MinZ
		xmax := las.Header.MaxX
		ymax := las.Header.MaxY
		zmax := las.Header.MaxZ
		if _, err := w.WriteString(fmt.Sprintf("%s,%.5f,%.5f,%.5f,%.5f,%.5f,%.5f\n", lasPath, xmin, ymin, zmin, xmax, ymax, zmax)); err != nil {
			log.Printf("fail to write csv file %v", err)
			return
		}
	}
}

func generateVectorMapList(inputRoot string) {
	shpPaths := []string{inputRoot}
	shpPaths = findFile(inputRoot, ".shp")
	shpListPath := filepath.Join(inputRoot, "vector_map_list.csv")
	f, err := os.Create(shpListPath)
	if err != nil {
		log.Printf("Fail to create csv at %s :  %v", shpListPath, err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()
	for _, shpPath := range shpPaths {
		fmt.Println(shpPath)
		if _, err := w.WriteString(fmt.Sprintf("%s,0\n", shpPath)); err != nil {
			log.Printf("fail to write csv file %v", err)
			return
		}
	}

}

func main() {
	dir := flag.String("dir", "./hd_maps", "input Folder")
	flag.Parse()
	inputRoot, err := filepath.Abs(*dir)
	if err != nil {
		log.Fatal("Fail to read Abs input Folder path")
		return
	}
	fileinfo, err := os.Stat(inputRoot)
	if os.IsNotExist(err) {
		log.Fatal("path does not exist.")
	}
	if fileinfo.IsDir() == false {
		log.Fatal("Fail to read input Folder path")
		return
	}
	generatePointCloudList(inputRoot)
	generateVectorMapList(inputRoot)
}

func findFile(root string, match string) (file []string) {

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		path, err = filepath.Abs(path)
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if strings.HasSuffix(info.Name(), match) {
			file = append(file, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("Total shp file : ", len(file))
	return file
}
