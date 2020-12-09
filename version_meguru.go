package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	fv         = "assembly: AssemblyFileVersion"
	av         = "assembly: AssemblyVersion"
	fvrc       = "FILEVERSION"
	pvrc       = "PRODUCTVERSION"
	vpvrc      = "VALUE \"FileVersion\""
	defaultVer = "1.0.1000.0"
)

func main() {
	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("please argument")
		os.Exit(-1)
	}

	fmt.Println("-----Search AssemblyInfo files----")
	rc := getAssemblyFiles(flag.Arg(0))

	showVersion(rc)

	fmt.Println("-----AssemblyInfo Update ?-----")
	// ユーザーの入力待ち
	fmt.Println("Please input update version. ex: 1.0.1000.0")
	fmt.Println("If you input nothing, update default version(1.0.1000.0")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	var changeVer string
	if changeVer = scanner.Text(); changeVer == "" {
		changeVer = defaultVer
	}

	updateVersion(rc, changeVer)
}

func dirwalk(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	var paths []string
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, dirwalk(filepath.Join(dir, file.Name()))...)
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
	}
	return paths
}

func getAssemblyFiles(path string) []string {
	var resources_path []string
	for _, file := range dirwalk(path) {
		if _, name := filepath.Split(file); name == "AssemblyInfo.cs" || filepath.Ext(name) == ".rc" {
			resources_path = append(resources_path, file)
		}
	}
	return resources_path
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func showVersion(files []string) {
	for _, n := range files {
		lines, err := readLines(n)
		if err != nil {
			fmt.Println(n, " is not read.")
		}
		for _, t := range lines {
			// FileVersionの出力
			if strings.Contains(t, fv) {
				fmt.Println(n, " : ", t)
			}

			if strings.Contains(t, av) {
				match, _ := regexp.MatchString("(\"*.*.*.0\")", t)
				if match {
					fmt.Println(n, " : ", t)
				}
			}
			if strings.Contains(t, fvrc) {
				fmt.Println(n, " : ", t)
			}
			if strings.Contains(t, pvrc) {
				fmt.Println(n, " : ", t)
			}
			if strings.Contains(t, vpvrc) {
				fmt.Println(n, " : ", t)
			}
		}
	}
}

func updateVersion(files []string, changeVer string) {
	for _, name := range files {
		b, err := ioutil.ReadFile(name)
		if err != nil {
			fmt.Println(os.Stderr, err)
			os.Exit(-1)
		}
		lines := string(b)

		rep := regexp.MustCompile(`(\"\d+.\d+.\d+.\d+\")`)
		arg := `"` + changeVer + `"`
		if !rep.MatchString(arg) {
			fmt.Println("ex. 1.0.1000.0")
			os.Exit(-1)
		}
		str := rep.ReplaceAllString(lines, arg)

		// .rc用に分岐
		if filepath.Ext(name) == ".rc" {
			repF := regexp.MustCompile(`FILEVERSION \d+,\d+,\d+,\d+`)
			aF := `FILEVERSION ` + convertRcFormat(changeVer)
			str = repF.ReplaceAllString(str, aF)

			repP := regexp.MustCompile(`PRODUCTVERSION \d+,\d+,\d+,\d+`)
			aP := `PRODUCTVERSION ` + convertRcFormat(changeVer)
			str = repP.ReplaceAllString(str, aP)

			repVLM := regexp.MustCompile(`V\d.\d\dL\d\d M\d\d`)
			str = repVLM.ReplaceAllString(str, convertAssemblyVLM(changeVer))
		}

		if err := ioutil.WriteFile(name, []byte(str), 0666); err != nil {
			fmt.Println(os.Stderr, err)
			os.Exit(-1)
		}
		fmt.Println(name, " : ", convertAssemblyVLM(changeVer))
	}
}

func convertAssemblyVLM(assembly string) string {
	rep := regexp.MustCompile(`\.`)
	result := rep.Split(assembly, -1)
	if len(result) != 4 {
		fmt.Println("error")
		os.Exit(-1)
	}

	// Vの整形
	var major string
	major = "V" + result[0] + ".0" + result[1]
	// L
	var livision string
	if len(result[2]) != 4 {
		livision = "L0" + result[2][:1]
	} else {
		livision = "L" + result[2][:2]
	}
	// M
	var minor string
	if len(result[3]) != 4 {
		if len(result[3]) != 3 {
			// case 0
			minor = " M00"
		} else {
			// 100 → M01
			minor = " M0" + result[3][:1]
		}
	} else {
		// 1000 → M10
		// 1100 → M11
		minor = " M" + result[3][:2]
	}
	ver := major + livision + minor
	return ver
}

// FILEVERSION 1,0,1000,0 に変換
func convertRcFormat(assembly string) string {
	rep := regexp.MustCompile(`\.`)
	result := rep.Split(assembly, -1)
	if len(result) != 4 {
		fmt.Println("error")
		os.Exit(-1)
	}

	return result[0] + "," + result[1] + "," + result[2] + "," + result[3]
}
