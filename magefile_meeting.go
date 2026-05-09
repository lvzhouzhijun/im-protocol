//go:build mage
// +build mage

package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/magefile/mage/mg"
)

type Meeting mg.Namespace

var meetingPath = filepath.Join(".", "openmeeting")

var meetingModules = []string{
	"admin",
	"meeting",
	"user",
}

func (Meeting) GenGo() error {
	// 将日志输出重定向到标准输出
	log.SetOutput(os.Stdout)
	log.Println("Generating Go code from meeting proto files")

	protoc, err := getToolPath("protoc")
	if err != nil {
		return err
	}

	var meetingPath = filepath.Join(".", "openmeeting")
	for _, module := range meetingModules {
		meetingGoOutPath := filepath.Join(meetingPath, module, GO)
		if err := os.MkdirAll(meetingGoOutPath, 0755); err != nil {
			return err
		}
		args := []string{
			"--go_out=" + meetingGoOutPath,      // --go_out: 指定生成标准 Protobuf Go 代码的输出目录。
			"--go-grpc_out=" + meetingGoOutPath, // --go-grpc_out: 指定生成 gRPC Go 代码的输出目录。
			// --go_opt / --go-grpc_opt 设置生成代码的模块路径（Go Module Path），确保导入路径正确。这里使用了硬编码的 github.com/lvzhouzhijun/im-protocol/openmeeting/ 作为模块前缀。
			"--go_opt=module=github.com/lvzhouzhijun/im-protocol/openmeeting/" + strings.Join([]string{module}, "/"),
			"--go-grpc_opt=module=github.com/lvzhouzhijun/im-protocol/openmeeting/" + strings.Join([]string{module}, "/"),
			// 输入文件：指定 .proto 文件的路径（openmeeting/<module>/<module>.proto）
			filepath.Join(meetingPath, module, module) + ".proto",
		}
		// 创建命令
		cmd := exec.Command(protoc, args...)
		connectStd(cmd)
		// 执行命令
		if err := cmd.Run(); err != nil {
			log.Printf("Error generating Go code for meeting module %s: %v\n", module, err)
			continue
		}
	}

	if err := removeOmitemptyTags(); err != nil {
		log.Println("Remove Omitempty is Error", err)
		return err
	} else {
		log.Println("Remove Omitempty is Success")
	}
	return nil
}

func (Meeting) GenJava() error {
	log.SetOutput(os.Stdout)
	log.Println("Generating Java code from meeting proto files")

	protoc, err := getToolPath("protoc")
	if err != nil {
		return err
	}

	var meetingPath = filepath.Join(".", "openmeeting")
	for _, module := range meetingModules {
		meetingJavaOutPath := filepath.Join(meetingPath, module, JAVA)
		if err := os.MkdirAll(meetingJavaOutPath, 0755); err != nil {
			return err
		}
		args := []string{
			"--java_out=lite:" + meetingJavaOutPath,
			filepath.Join(meetingPath, module, module) + ".proto",
		}

		cmd := exec.Command(protoc, args...)
		connectStd(cmd)

		if err := cmd.Run(); err != nil {
			log.Printf("Error generating Java code for meeting module %s: %v\n", module, err)
			continue
		}
	}
	return nil
}

func (Meeting) GenKotlin() error {
	log.SetOutput(os.Stdout)
	log.Println("Generating Kotlin code from meeting proto files")

	protoc, err := getToolPath("protoc")
	if err != nil {
		return err
	}

	var meetingPath = filepath.Join(".", "openmeeting")
	for _, module := range meetingModules {
		meetingKotlinOutPath := filepath.Join(meetingPath, module, Kotlin)
		if err := os.MkdirAll(meetingKotlinOutPath, 0755); err != nil {
			return err
		}
		args := []string{
			"--kotlin_out=lite:" + meetingKotlinOutPath,
			filepath.Join(meetingPath, module, module) + ".proto",
		}

		cmd := exec.Command(protoc, args...)
		connectStd(cmd)

		if err := cmd.Run(); err != nil {
			log.Printf("Error generating Kotlin code for meeting module %s: %v\n", module, err)
			continue
		}
	}
	return nil
}

func (Meeting) GenCSharp() error {
	log.SetOutput(os.Stdout)
	log.Println("Generating C# code from proto files")

	protoc, err := getToolPath("protoc")
	if err != nil {
		return err
	}

	var meetingPath = filepath.Join(".", "openmeeting")
	for _, module := range meetingModules {
		meetingCsharpOutPath := filepath.Join(meetingPath, module, CSharp)
		if err := os.MkdirAll(meetingCsharpOutPath, 0755); err != nil {
			return err
		}
		args := []string{
			"--csharp_out=" + meetingCsharpOutPath,
			filepath.Join(meetingPath, module, module) + ".proto",
		}

		cmd := exec.Command(protoc, args...)
		connectStd(cmd)

		if err := cmd.Run(); err != nil {
			log.Printf("Error generating C# code for module %s: %v\n", module, err)
			continue
		}
	}
	return nil
}

func (Meeting) GenJavaScript() error {
	log.SetOutput(os.Stdout)
	log.Println("Generating JavaScript code from proto files")

	protoc, err := getToolPath("protoc")
	if err != nil {
		return err
	}

	var meetingPath = filepath.Join(".", "openmeeting")
	for _, module := range meetingModules {
		meetingJsOutPath := filepath.Join(meetingPath, module, JS)
		if err := os.MkdirAll(meetingJsOutPath, 0755); err != nil {
			return err
		}
		args := []string{
			"--js_out=import_style=commonjs,binary:" + meetingJsOutPath,
			filepath.Join(meetingPath, module, module) + ".proto",
		}

		cmd := exec.Command(protoc, args...)
		connectStd(cmd)

		if err := cmd.Run(); err != nil {
			log.Printf("Error generating JS code for module %s: %v\n", module, err)
			continue
		}
	}
	return nil
}

func (Meeting) GenTypeScript() error {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)
	log.Println("Generating TypeScript code from proto files")

	protoc, err := getToolPath("protoc")
	if err != nil {
		return err
	}

	tsProto := filepath.Join(".", "node_modules", ".bin", "protoc-gen-ts_proto")
	if runtime.GOOS == "windows" {
		tsProto = filepath.Join(".", "node_modules", ".bin", "protoc-gen-ts_proto.cmd")
	}

	if _, err := os.Stat(tsProto); err != nil {
		log.Println("tsProto Not Found. Error: ", err, " tsProto Path: ", tsProto)
		return err
	}

	for _, module := range meetingModules {
		meetingTsOutDir := filepath.Join(meetingPath, "pb", TS)
		if err := os.MkdirAll(meetingTsOutDir, 0755); err != nil {
			return err
		}
		args := []string{
			"--plugin=protoc-gen-ts_proto=" + tsProto,
			"--ts_proto_opt=messages=true,outputJsonMethods=false,outputPartialMethods=false,outputClientImpl=false,outputEncodeMethods=false,useOptionals=messages",
			"--ts_proto_out=" + meetingTsOutDir,
			filepath.Join(meetingPath, module, module) + ".proto",
		}

		cmd := exec.Command(protoc, args...)
		connectStd(cmd)

		if err := cmd.Run(); err != nil {
			log.Printf("Error generating TypeScript code for module %s: %v\n", module, err)
			continue
		}
	}
	return nil
}

func (Meeting) GenSwift() error {
	log.SetOutput(os.Stdout)
	log.Println("Generating Swift code from proto files")

	protoc, err := getToolPath("protoc")
	if err != nil {
		return err
	}

	for _, module := range meetingModules {
		swiftOutDir := filepath.Join(".", module, SWIFT)

		modulePath := filepath.Join(module, module+".proto")

		if err := os.MkdirAll(swiftOutDir, 0755); err != nil {
			return err
		}
		args := []string{
			"--swift_out=" + swiftOutDir,
			"--swift_opt=Visibility=" + "Public",
			modulePath,
		}

		cmd := exec.Command(protoc, args...)
		connectStd(cmd)

		if err := cmd.Run(); err != nil {
			log.Printf("Error generating Swift code for module %s: %v\n", module, err)
			continue
		}
		log.Printf("Successfully generated Swift code for module %s\n", module)
	}
	return nil
}
