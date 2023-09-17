# NUS Canvas File Downloader

> Download files from currently-enrolled Canvas courses.

## Installing

There are 2 ways to use this application:

1. Download the latest release from the [releases page](https://github.com/yusufaine/nus-canvas-cli/releases), or
2. Clone this repository and run `go run main.go` to install dependencies. (Requires [Go](https://golang.org/dl/) v1.20+).

## Usage

1. Generate a Canvas access token [here](https://canvas.nus.edu.sg/profile/settings)
![Generate Canvas access token](https://gist.githubusercontent.com/yusufaine/23cea8a7a4f0fe3714f81d19944cbda7/raw/6b94cf370e05f1db4cf75215bdea845561603d78/01_generate_token.png)
2. Run the application and either:
   * (Recommended) Create a `.token` file in the same directory as the application and paste the token in there, or
   * Pass in the token as a command-line argument (specifying `--store` is recommended to save the token for future use).

## Demo

> Using the recommended method of storing the token in a `.token` file.

![Demo](https://gist.githubusercontent.com/yusufaine/23cea8a7a4f0fe3714f81d19944cbda7/raw/d1acda94510f6a6de9d67c62b61e0e5bf76c6c2e/02_demo.gif)
