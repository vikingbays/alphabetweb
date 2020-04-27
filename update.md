2020-04-27

1、优化可执行文件编译模式
1）、优化-build模式，可以更好的支持交叉平台编译，例如：android/arm64
    原来的方式：
      go run /xxx/alphabetweb/src/alphabet/cmd/abserver.go  -build  "project"  "project_out"  "mac"
    优化后的方式：
      go run /xxx/alphabetweb/src/alphabet/cmd/abserver.go  -build  "project"  "project_out"  "darwin/amd64"

    第三个参数的如参数信息，可以参考：go tool dist list

2）、如果涉及到android的编译，可以指定  CC、CXX、CGO_ENABLED、GO111MODULE参数，在执行命令前定义(用AB_前缀)，例如：
      AB_CC=/xxxx/android-ndk-r19c/toolchains/llvm/prebuilt/darwin-x86_64/bin/aarch64-linux-android21-clang        \
      AB_CXX=/xxxx/android-ndk-r19c/toolchains/llvm/prebuilt/darwin-x86_64/bin/aarch64-linux-android21-clang++     \
      AB_CGO_ENABLED=1                                                                                             \
      AB_GO111MODULE=off                                                                                           \
      go run /xxxx/alphabetweb/src/alphabet/cmd/abserver.go  -build "..." "..." "android/arm64"
    其中：CC、CXX指向外部安装的编译环境，这儿是NDK的android编译环境。
