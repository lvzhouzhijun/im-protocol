//go:build mage
// +build mage

package main

import (
	"fmt"
	"go/build"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var Default = InstallDepend

var Aliases = map[string]any{
	"go":     GenGo,
	"java":   GenJava,
	"kotlin": GenKotlin,
	"csharp": GenCSharp,
	"js":     GenJavaScript,
	"ts":     GenTypeScript,
	"swift":  GenSwift,

	"dep": InstallDepend,

	"m:go":     Meeting.GenGo,
	"m:java":   Meeting.GenJava,
	"m:kotlin": Meeting.GenKotlin,
	"m:csharp": Meeting.GenCSharp,
	"m:js":     Meeting.GenJavaScript,
	"m:ts":     Meeting.GenTypeScript,
	"m:swift":  Meeting.GenSwift,
}

// 语言目标
// 为每种目标语言定义输出目录
const (
	GO     = "go"
	JAVA   = "java"
	CSharp = "csharp"
	Kotlin = "kotlin"
	JS     = "js"
	TS     = "ts"
	RS     = "rust"
	SWIFT  = "swift"
)

var protoModules = []string{
	"auth",
	"conversation",
	"errinfo",
	"group",
	"jssdk",
	"msg",
	"msggateway",
	"push",
	"relation",
	"rtc",
	"sdkws",
	"third",
	"user",
	"wrapperspb",
}

func InstallDepend() error {
	log.SetOutput(os.Stdout)
	log.Println("installing protobuf dependencies in Go.")

	cmds := [][]string{
		{"install", "google.golang.org/protobuf/cmd/protoc-gen-go@latest"},
		{"install", "google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"},
	}

	for _, cmdArgs := range cmds {
		cmd := exec.Command("go", cmdArgs...)
		connectStd(cmd)
		if err := cmd.Run(); err != nil {
			log.Printf("command %v error: %v", cmdArgs, err)
			return err
		}
	}
	return nil
}

func GenDocs() error {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)
	log.Println("Generating documentation from proto files")

	protoc, err := getToolPath("protoc")
	if err != nil {
		return err
	}

	docsOutDir := filepath.Join(".", "docs")
	for _, module := range protoModules {
		if err := os.MkdirAll(filepath.Join(docsOutDir, module), 0755); err != nil {
			return err
		}

		args := []string{
			"--doc_out=" + filepath.Join(docsOutDir),
			"--doc_opt=markdown," + strings.Join([]string{module, "md"}, "."),
			filepath.Join(module, module) + ".proto",
		}

		cmd := exec.Command(protoc, args...)
		connectStd(cmd)
		if err := cmd.Run(); err != nil {
			log.Printf("Error generating documentation for module %s: %v\n", module, err)
			continue
		}
	}
	return nil
}

func AllProtobuf() error {
	if err := GenGo(); err != nil {
		return err
	}
	if err := GenJava(); err != nil {
		return err
	}
	if err := GenCSharp(); err != nil {
		return err
	}
	if err := GenJavaScript(); err != nil {
		return err
	}
	if err := GenTypeScript(); err != nil {
		return err
	}
	return nil
}

func GenGo() error {
	// 将日志输出重定向到标准输出
	log.SetOutput(os.Stdout)
	log.Println("Generating Go code from proto files")

	protoc, err := getToolPath("protoc")
	if err != nil {
		return err
	}

	for _, module := range protoModules {
		args := []string{
			"--go_out=" + filepath.Join(".", module),      // --go_out: 指定生成标准 Protobuf Go 代码的输出目录。
			"--go-grpc_out=" + filepath.Join(".", module), // --go-grpc_out: 指定生成 gRPC Go 代码的输出目录。
			// --go_opt / --go-grpc_opt 设置生成代码的模块路径（Go Module Path），确保导入路径正确。这里使用了硬编码的 github.com/lvzhouzhijun/im-protocol/ 作为模块前缀。
			"--go_opt=module=github.com/lvzhouzhijun/im-protocol/" + strings.Join([]string{module}, "/"),
			"--go-grpc_opt=module=github.com/lvzhouzhijun/im-protocol/" + strings.Join([]string{module}, "/"),
			// 输入文件：指定 .proto 文件的路径（openmeeting/<module>/<module>.proto）
			filepath.Join(module, module) + ".proto",
		}
		// 创建命令
		cmd := exec.Command(protoc, args...)
		connectStd(cmd)
		// 执行命令
		if err := cmd.Run(); err != nil {
			log.Printf("Error generating Go code for module %s: %v\n", module, err)
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

func GenJava() error {
	log.SetOutput(os.Stdout)
	log.Println("Generating Java code from proto files")

	protoc, err := getToolPath("protoc")
	if err != nil {
		return err
	}

	for _, module := range protoModules {
		javaOutDir := filepath.Join(".", module, JAVA)
		if err := os.MkdirAll(javaOutDir, 0755); err != nil {
			return err
		}
		args := []string{
			"--java_out=lite:" + javaOutDir,
			filepath.Join(module, module) + ".proto",
		}
		log.Println(javaOutDir)

		cmd := exec.Command(protoc, args...)
		connectStd(cmd)

		if err := cmd.Run(); err != nil {
			log.Printf("Error generating Java code for module %s: %v\n", module, err)
			continue
		}
	}
	return nil
}

func GenKotlin() error {
	log.SetOutput(os.Stdout)
	log.Println("Generating Kotlin code from proto files")

	protoc, err := getToolPath("protoc")
	if err != nil {
		return err
	}

	for _, module := range protoModules {
		kotlinOutDir := filepath.Join(".", module, Kotlin)
		if err := os.MkdirAll(kotlinOutDir, 0755); err != nil {
			return err
		}
		args := []string{
			"--kotlin_out=" + kotlinOutDir,
			filepath.Join(module, module) + ".proto",
		}

		cmd := exec.Command(protoc, args...)
		connectStd(cmd)

		if err := cmd.Run(); err != nil {
			log.Printf("Error generating Kotlin code for module %s: %v\n", module, err)
			continue
		}
	}
	return nil
}

func GenCSharp() error {
	log.SetOutput(os.Stdout)
	log.Println("Generating C# code from proto files")

	protoc, err := getToolPath("protoc")
	if err != nil {
		return err
	}

	for _, module := range protoModules {
		csharpOutDir := filepath.Join(".", module, CSharp)
		if err := os.MkdirAll(csharpOutDir, 0755); err != nil {
			return err
		}
		args := []string{
			"--csharp_out=" + csharpOutDir,
			filepath.Join(module, module) + ".proto",
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

func GenJavaScript() error {
	log.SetOutput(os.Stdout)
	log.Println("Generating JavaScript code from proto files")

	protoc, err := getToolPath("protoc")
	if err != nil {
		return err
	}

	jsDir := filepath.Join(".", "pb", "js")
	args := []string{
		"--js_out=import_style=commonjs,binary:" + jsDir,
	}

	if err := os.MkdirAll(jsDir, 0755); err != nil {
		return err
	}

	for _, module := range protoModules {
		jsOutDir := filepath.Join(".", module, JS)
		if err := os.MkdirAll(jsOutDir, 0755); err != nil {
			return err
		}
		args = append(args, filepath.Join(module, module)+".proto")
	}

	cmd := exec.Command(protoc, args...)
	connectStd(cmd)

	if err := cmd.Run(); err != nil {
		log.Printf("Error generating JS code %v\n", err)
	}

	return nil
}

func GenTypeScript() error {
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

	for _, module := range protoModules {
		tsOutDir := filepath.Join("pb", TS)
		if err := os.MkdirAll(tsOutDir, 0755); err != nil {
			return err
		}
		args := []string{
			"--plugin=protoc-gen-ts_proto=" + tsProto,
			"--ts_proto_opt=messages=true,outputJsonMethods=false,outputPartialMethods=false,outputClientImpl=false,outputEncodeMethods=false,useOptionals=messages",
			"--ts_proto_out=" + tsOutDir,
			filepath.Join(module, module) + ".proto",
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

func GenSwift() error {
	log.SetOutput(os.Stdout)
	log.Println("Generating Swift code from proto files")

	protoc, err := getToolPath("protoc")
	if err != nil {
		return err
	}

	for _, module := range protoModules {
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

func GenHarmonyTS() error {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)

	log.Println("Generating Harmony TypeScript code from proto files")

	outJSFile := "proto.js"
	args := []string{
		"-t", "static-module",
		"-w", "es6",
		"-o", outJSFile,
	}

	for _, module := range protoModules {
		protoFile := filepath.Join(module, module) + ".proto"
		args = append(args, protoFile)
	}

	jscmd := exec.Command("pbjs", args...)
	jscmd.Env = os.Environ()
	connectStd(jscmd)

	log.Println("Running harmony js command", jscmd.String())
	if err := jscmd.Run(); err != nil {
		log.Printf("Error generating Harmony JS code: %v\n", err)
	}

	outTSDefFile := "proto.d.ts"
	tscmd := exec.Command("pbts", outJSFile, "-o", outTSDefFile)
	tscmd.Env = os.Environ()
	connectStd(tscmd)

	log.Println("Running harmony ts command", tscmd.String())
	if err := tscmd.Run(); err != nil {
		log.Printf("Error generating Harmony TS code: %v\n", err)
	}

	replaceStr := func(filePath, oldStr, newStr string) {
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Panic("failed to read file: %w", err)
		}

		originalContent := string(content)
		modifiedContent := strings.Replace(originalContent, oldStr, newStr, 1)
		if originalContent == modifiedContent {
			return
		}
		err = os.WriteFile(filePath, []byte(modifiedContent), 0644)
		if err != nil {
			log.Panic("failed to write file: %w", err)
		}

	}
	replaceStr(outJSFile, "import * as $protobuf from \"protobufjs/minimal\";", "import { index } from \"@ohos/protobufjs\"; \nconst $protobuf = index; \n import Long from 'long';\n$protobuf.util.Long=Long \n$protobuf.configure()")
	replaceStr(outTSDefFile, "import * as $protobuf from \"protobufjs\";\nimport Long = require(\"long\");", "import * as $protobuf from \"@ohos/protobufjs\"\nimport Long from 'long';")

	return nil
}

// connectStd 将外部命令的标准输出（Stdout）和标准错误（Stderr）直接“透传”到当前 Go 程序的终端。
// 简单来说，它让被调用的外部命令（如 protoc、pbjs 等）像在当前终端直接运行一样，实时打印日志，而不是把输出憋在内存里。
func connectStd(cmd *exec.Cmd) {
	// cmd.Stdout：这是外部命令产生正常输出（比如打印的日志、生成的数据）的地方。
	// os.Stdout：这是当前 Go 程序的标准输出，通常对应你运行程序的终端窗口。
	// 这行代码将外部命令的“嘴巴”接到了当前终端的“耳朵”上。外部命令打印什么，你的屏幕上就立刻显示什么。
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
}

// getWorkDirToolPath 在当前工作目录下的 tools 子目录中，递归查找指定名称的工具。
func getWorkDirToolPath(name string) string {
	toolPath := ""
	// os.Getwd() 是一个标准库函数，用于获取程序运行时的“当前工作目录”（Working Directory）的绝对路径。
	workDir, err := os.Getwd()
	if err != nil {
		log.Println("Error", err)
		return toolPath
	}
	// 这里它将 workDir（当前工作目录）和字符串 "tools" 拼接，形成目标目录的完整路径。
	toolsPath := filepath.Join(workDir, "tools")
	// filepath.Walk 函数会递归地遍历 toolsPath 目录下的所有文件和子目录。
	filepath.Walk(toolsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		/*
			info.Name()：获取当前遍历到的文件或目录的名称（例如 gofmt.exe 或 protoc）。
			filepath.Ext(info.Name())：获取文件名的扩展名（例如 .exe 或 ""）。
			strings.TrimSuffix(...)：从文件名中移除扩展名。
			如果文件名是 gofmt.exe，处理后变成 gofmt。
			如果文件名是 protoc（无扩展名），处理后仍然是 protoc。
			这行代码的目的是比较不带扩展名的文件名是否与要查找的 name 相等。
		*/
		if strings.TrimSuffix(info.Name(), filepath.Ext(info.Name())) == name {
			toolPath = path
		}
		return nil
	})
	return toolPath
}

// getToolPath 按照特定优先级顺序查找指定工具（可执行文件）的完整路径。
// 它依次在三个位置进行查找：
//
//	项目级：优先在项目工作目录内查找，方便使用项目特定的工具版本。
//	系统级：其次在系统全局路径 (PATH) 中查找，这是最常见的查找方式。
//	Go 级：最后在 Go 语言特有的 GOPATH/bin 目录下查找，这是 go install 或 go get 安装工具的默认位置。
func getToolPath(name string) (string, error) {
	// 在工作目录中查找工具路径。
	toolPath := getWorkDirToolPath(name)
	if toolPath != "" {
		return toolPath, nil
	}

	/*
		在系统 PATH 环境变量中查找。
		如果第一步未找到，则使用标准库 exec.LookPath(name) 进行查找。
		这个函数会在 PATH 环境变量列出的所有目录中搜索名为 name 的可执行文件
	*/
	if p, err := exec.LookPath(name); err != nil {
		return p, nil
	}
	/*
		如果前两步都失败了，代码开始准备在 Go 的 GOPATH 目录下查找。
		os.Getenv("GOPATH") 获取当前系统设置的 GOPATH 环境变量的值。
	*/
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		// 如果系统环境变量中没有设置 GOPATH，则使用 Go 标准库 build 包中的默认 GOPATH 值作为后备。
		gopath = build.Default.GOPATH
	}
	p := filepath.Join(gopath, "bin", name)
	// os.Stat(p) 用于获取文件 p 的信息。如果文件不存在，它会返回一个错误。
	if _, err := os.Stat(p); err != nil {
		return "", err
	}
	return p, nil
}

// removeOmitemptyTags 批量移除protobuf生成的Go文件中的 omitempty 标签。
func removeOmitemptyTags() error {
	/*
		正则表达式 ,\s*omitempty 匹配以下内容：
			,：匹配字面量逗号
			\s*：匹配零个或多个空白字符（空格、制表符等）
			omitempty：匹配字面量字符串 "omitempty"
			这个正则表达式用于查找类似 ,omitempty 或 , omitempty 的模式
	*/
	re := regexp.MustCompile(`,\s*omitempty`)
	/*
		filepath.Walk()：递归遍历指定目录下的所有文件和子目录
			"."：表示从当前目录开始遍历
			第二个参数是一个回调函数，对每个访问的文件/目录都会调用这个函数
			回调函数参数：
				path string：当前文件/目录的完整路径
				info os.FileInfo：包含文件信息的对象
				err error：可能发生的错误
	*/
	return filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		// 如果在访问路径过程中发生错误，直接返回错误
		if err != nil {
			fmt.Println("access path error:", err)
			return err
		}

		/*
			文件类型判断：
				!info.IsDir()：检查是否不是目录（即为普通文件）
				strings.HasSuffix(path, ".pb.go")：检查文件名是否以 ".pb.go" 结尾
				只处理非目录且扩展名为 .pb.go 的文件（protobuf生成的Go文件）
		*/
		if !info.IsDir() && strings.HasSuffix(path, ".pb.go") {
			// 读取整个文件的内容到内存中
			input, err := os.ReadFile(path)
			if err != nil {
				fmt.Println("ReadFile error. Path: %s, Error %v", path, err)
				return err
			}
			/*
				执行替换操作：
				string(input)：将字节数组转换为字符串
				re.ReplaceAllString()：使用之前定义的正则表达式，在字符串中查找所有匹配项，并将其替换为空字符串（即删除）
				结果存储在 output 变量中
			*/
			output := re.ReplaceAllString(string(input), "")
			// 检查是否发生了替换：比较原始内容和处理后的内容，如果不相同，说明有匹配项被替换了。
			if string(input) != output {
				/*
					写回文件：
						[]byte(output)：将修改后的字符串转换回字节数组
						info.Mode()：使用原始文件的权限模式
						将修改后的内容写回到原文件
				*/
				err = os.WriteFile(path, []byte(output), info.Mode())
				if err != nil {
					fmt.Println("Error writing file: %s, error: %v\n", path, err)
					return err
				}
			}
		}
		return nil
	})
}
