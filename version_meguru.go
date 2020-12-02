package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

func main() {
	name := "AssemblyInfo.cs"
	b, err := ioutil.ReadFile(name)
	if err != nil {
		fmt.Println(os.Stderr, err)
		os.Exit(1)
	}
	lines := string(b)

	rep := regexp.MustCompile(`\d+.\d+.\d+.\d+`)
	arg := os.Args[1]
	if !rep.MatchString(arg) {
		fmt.Println("ex. 1.0.1000.0")
		os.Exit(1)
	}
	str := rep.ReplaceAllString(lines, arg)
	fmt.Println(str)

	if err := ioutil.WriteFile(name, []byte(str), 0666); err != nil {
		fmt.Println(os.Stderr, err)
		os.Exit(1)
	}
	convert_assembly_vlm(arg)
}

// アセンブリバージョンを V1.xLxx Mxx に変換する
// `.` で区切る
// 未完成
func convert_assembly_vlm(assembly string) string {
	fmt.Println(assembly)
	rep := regexp.MustCompile(`\.`)
	result := rep.Split(assembly, -1)
	fmt.Println(result)
	if len(result) != 4 {
		fmt.Println("error")
		os.Exit(1)
	}

	// Vの整形
	var major string
	major = "V" + result[0] + ".0" + result[1]
	// Lの整形
	var livision string
	if len(result[2]) != 4 {
		livision = "L0" + result[2][:1]
	} else {
		livision = "L" + result[2][:2]
	}
	// Mの整形
	var minor string
	if len(result[3]) != 4 {
		minor = " M00"
	} else {
		minor = " M0" + result[3][:1]
	}
	ver := major + livision + minor
	fmt.Println(ver)
	return ver
}
