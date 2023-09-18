# Canvas LMS File Downloader

>Download files from currently-enrolled [Canvas](https://www.instructure.com/canvas) courses. 

Developed for students from [National University of Singapore](https://nus.edu.sg), extendable to other institutions (see [Non-NUS Canvas Users](#non-nus-canvas-users)).

<!-- omit in toc -->
## Table of Contents

- [Installing](#installing)
  - [Caveats for MacOS users](#caveats-for-macos-users)
- [Usage](#usage)
  - [Flags](#flags)
  - [Non-NUS Canvas Users](#non-nus-canvas-users)
- [Demo](#demo)
- [Contributing](#contributing)

## Installing

There are 2 ways to use this application:

1. Download the latest release from the [releases page](https://github.com/yusufaine/nus-canvas-cli/releases), or
2. Clone this repository and run `go run main.go` to install dependencies. (Requires [Go](https://golang.org/dl/) v1.20+).

### Caveats for MacOS users

MacOS users may encounter an error when running the binary from the releases page for the first time due to MacOS's security settings. To fix this, you will need to:

1. Head over to `System Settings > Privacy & Security`, and scroll down to the `Security` section, and click `Allow Anyway`, and
2. run the application binary again and allow it to run when prompted.

![Allow application to run](https://gist.githubusercontent.com/yusufaine/23cea8a7a4f0fe3714f81d19944cbda7/raw/dc64c05a08d5331355d75102ee71f56d1f1119ce/03_mac_caveat.png)

## Usage

1. Generate a Canvas access token [here](https://canvas.nus.edu.sg/profile/settings)
![Generate Canvas access token](https://gist.githubusercontent.com/yusufaine/23cea8a7a4f0fe3714f81d19944cbda7/raw/6b94cf370e05f1db4cf75215bdea845561603d78/01_generate_token.png)
2. Run the application and either:
   - (Recommended) Create a `.token` file in the same directory as the application and paste the token in there, or
   - Pass in the token as a command-line argument (specifying `--store` is recommended to save the token for future use).

### Flags

Run the application with the `--help` flag to see the list of available flags.

### Non-NUS Canvas Users

If you are not from NUS, you will need to specify the Canvas URL using the `--host` flag. For example, if you are from Yale University, you will need to run the application as follows, assuming you have a `.token` file in the same directory as the application:

```bash
./canvas-cli --host canvas.yale.edu
```

## Demo

> Using the recommended method of storing the token in a `.token` file. By default, the application will only download files that are <= 10MB in size. This can be changed by specifying the `--size=SIZE_IN_MB` flag (e.g. `--size=50` to download files <= 50MB in size)

![Demo](https://gist.github.com/yusufaine/23cea8a7a4f0fe3714f81d19944cbda7/raw/0ce4e75ac7024bda57b6bdb15c6b9cec215be76d/02_demo.gif)

## Contributing

This project is open to contributions. Feel free to open an issue or submit a pull request detailing your the issues encountered, or changes made respectively.
