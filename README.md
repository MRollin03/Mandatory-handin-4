# Project Setup and Running Instructions

## 1. Setup

To get started with this project, follow these steps:

1. **Clone the Repository**
   ```bash
   git clone https://github.com/MRollin03/Mandatory-handin-4
   cd Mandatory-handin-4
   ```

- import dependencies
  - `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`
  - `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`
- follow next steps

## 2: How to run the program''

To run the program, execute the following commands:

Navigate to the Main Directory

- cd from root
- `cd main`
- `go run tokenRing.go`
- Program should be up and running

Verify the program is running, It should be requesting a command. You can check the logs or access the relevant endpoints to ensure it's functioning correctly.
Once the program is requesting input:

- provide input using the format: `request <id>` where `<id>` is a int from 1 to 5
- this command asks the node in question to request access to the critical resource.
- The program creates 5 nodes on startup. node 1 is defined as the Primary node, hence the 1-5 constraint.
