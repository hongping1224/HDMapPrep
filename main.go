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
	shp "github.com/jonas-p/go-shp"
)

func generatePointCloudList(inputRoot string, useRelative bool) {
	lasPaths := []string{inputRoot}
	lasPaths = findFile(inputRoot, ".las", useRelative)
	pointCloudListPath := filepath.Join(inputRoot, "point_cloud_map_list.csv")
	f, err := os.Create(pointCloudListPath)
	if err != nil {
		log.Printf("Fail to create csv at %s :  %v", pointCloudListPath, err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()
	for _, lasPath := range lasPaths {
		las, err := lidario.NewLasFile(lasPath, "rh")
		if useRelative {
			lasPath = lasPath[strings.Index(lasPath, "point_cloud_map"):]
		}
		fmt.Println(lasPath)
		if err != nil {
			log.Printf("%s read fail. Err : %v", las, err)
			continue
		}
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
		las.Close()
	}
}

func generateVectorMapList(inputRoot string, useRelative bool) {
	shpPaths := []string{inputRoot}
	shpPaths = findFile(inputRoot, ".shp", useRelative)
	shpListPath := filepath.Join(inputRoot, "vector_map_list.csv")
	f, err := os.Create(shpListPath)
	if err != nil {
		log.Printf("Fail to create csv at %s :  %v", shpListPath, err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()
	for _, shpPath := range shpPaths {
		shape, err := shp.Open(shpPath)
		defer shape.Close()
		if useRelative {
			shpPath = shpPath[strings.Index(shpPath, "vector_map"):]
		}
		if err != nil {
			log.Printf("Fail to open %s  with Error: %s", shpPath, err)
			continue
		}
		xmin := shape.BBox().MinX
		ymin := shape.BBox().MinY
		xmax := shape.BBox().MaxX
		ymax := shape.BBox().MaxY
		fmt.Println(shpPath)

		if _, err := w.WriteString(fmt.Sprintf("%s,%.5f,%.5f,%.5f,%.5f\n", shpPath, xmin, ymin, xmax, ymax)); err != nil {
			log.Printf("fail to write csv file %v", err)
			return
		}
	}

}

func main() {
	dir := flag.String("dir", "./hd_maps", "input Folder")
	relativePath := flag.Bool("rel", false, "use relative path")
	flag.Parse()
	inputRoot := *dir
	fmt.Println(*dir)
	fmt.Println(*relativePath)
	if !*relativePath {
		inputRoot, _ = filepath.Abs(*dir)
	}
	fileinfo, err := os.Stat(inputRoot)
	if os.IsNotExist(err) {
		log.Fatal("path does not exist.")
	}
	if fileinfo.IsDir() == false {
		log.Fatal("Fail to read input Folder path")
		return
	}
	generatePointCloudList(inputRoot, *relativePath)
	generateVectorMapList(inputRoot, *relativePath)
}

func findFile(root string, match string, useRelative bool) (file []string) {

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if useRelative == false {
			path, err = filepath.Abs(path)
			if err != nil {
				fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
				return err
			}
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
