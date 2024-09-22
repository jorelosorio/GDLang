# GDLang

GDLang is a fun, easy-to-learn, and type-safe programming language for modern development.

[![Test](https://github.com/jorelosorio/gdlang/actions/workflows/test.yml/badge.svg)](https://github.com/jorelosorio/gdlang/actions/workflows/test.yml) [![Release](https://github.com/jorelosorio/gdlang/actions/workflows/release.yml/badge.svg)](https://github.com/jorelosorio/gdlang/actions/workflows/release.yml)

## Features

- 🙂 **Super Syntax**: The syntax is designed to be exceptionally friendly and easy to read and write. Inspired by `Swift`, `TypeScript`, and `Go`, it combines modern, expressive syntax with a robust statically typed design.
- 🚀 **Fast and Efficient**: GDLang compiles code into `bytecode` that runs on an optimized `Virtual Machine (VM)`. The binaries are incredibly lightweight, occupying minimal space while ensuring high performance and fast execution across various environments.
- 🎉 **Single Binary, No Hassle**: The `VM` is a single binary with a size of just less than ~1.7MB and a basic `binary program` is ~163 bytes.
- 💪🏽 **Powered by [GoLang](https://golang.org)**: Built on the robust GoLang ecosystem, GDLang inherits Go's performance, reliability and networking capabilities.

## Why a new language?

The motivation behind creating this new programming language is to bridge the gap between complexity and accessibility. While many existing languages are powerful, they often come with steep learning curves and intricate syntax, which can reduce productivity and hinder adoption—especially for newcomers and those in fast-paced development environments.

The language prioritizes simplicity, readability, and ease of use, offering an all-in-one tool that can be extended to support a broad spectrum of applications but particularly well-suited for networking and I/O operations.

The syntax is designed to feel intuitive for developers with a background in web development, making it an ideal transitional language. This approach not only lowers the barrier to entry for newcomers but also provides a seamless adoption to start building and innovating quickly and efficiently.

## How it looks like?

Here's a simple example of a GDLang program, that demonstrates the basic syntax, with functions, variables, comments, and the `main` function:

```gdlang
typealias route = {
    name: string,
    handler: func() => string,
}

set index: route = {
    name: "/",
    handler: func() => string {
        return "GDLang!"
    },
}

// Every program needs a main function
pub func main() {
    set message = "Hello, " + index.handler() + "!"
    println(message)
}
```

## Getting Started 🧑🏽‍💻

Download the [Latest version](https://github.com/jorelosorio/gdlang/releases/latest) that matches your OS and architecture.

> For instance, if you are using a `Mac` with an `M1`/`M2` chip, you should download the `darwin-arm64` version.

Unzip the downloaded file and inside the folder, you will find the following binaries:

- `gdc` - The Compiler
- `gdvm` - The Virtual Machine
- `gdcvm` - The Compiler and Virtual Machine

Add `gdc` and `gdvm` to your `$PATH` to use them globally.

    export PATH=$PATH:/path/to/gdlang

### Lets compile a basic program

Create a new folder called `hello` and navigate to it.

> `NOTE:` this folder is going to be the `package` name.

Inside the `hello` folder, create a new file called `main.gd` and append the following code:

```gdlang
pub func main() {
    println("Hello, GDLang!")
}
```
The following command can be used to create the folder package and the file all together:

```bash
mkdir hello && echo 'pub func main() {
    println("Hello, GDLang!")
}' > hello/main.gd
```

Run the following command to compile your package:

```bash
gdc -pkg hello
```

> Note: The `-pkg` flag is used to specify the package name. In this case, the package name is `hello`.

After running the compiler, you will see a new file called `hello.gdbin` in the `hello` folder. This file contains the compiled bytecode of your program.

> NOTE: There is also a new file called `hello.gdmap` that maps the bytecode to the source code.

### Lets run the compiled program

Run the following command to execute the compiled bytecode:

```bash
gdvm -gdbin ./hello/hello.gdbin
```

You should see the following output:

```bash
Hello, GDLang!
```

🎉 Congratulations! You have successfully compiled and run your first GDLang program.


## 🚀 What's Coming Next?

- 🔄 **Threads and Channels Support**: Introduce support for threads and channels, similar to Go's concurrency model.
- 📡 **Networking and I/O**: Add support for networking and I/O operations to enable communication with external systems. `WIP`
- 📦 **Shared Libraries and SDKs**: Build a robust library ecosystem to be shared across other programs.
- 📜 **Standard Library**: Introduce a comprehensive standard library with built-in functions and utilities.

> **Note:** This roadmap is subject to change based on community feedback and contributions.

## Contributing

We’re excited to welcome contributions! Whether you’re fixing bugs, improving documentation, or adding new features, we’d love your help in making GDLang better. Before you start, please take a moment to read our [CONTRIBUTING](./CONTRIBUTING.md) guide. It outlines our processes and best practices to ensure smooth collaboration.

## License

This project is licensed under the GNU General Public License v3.0 (GPLv3). This means you are free to copy, modify, and distribute the code. However, any derivative works or modifications must also be licensed under the GPLv3. You can't use this code to create proprietary software; any software derived from or incorporating this code must also be open-sourced under the GPLv3.

See the [LICENSE](./LICENSE.txt) file for details.